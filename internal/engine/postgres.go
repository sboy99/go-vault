package engine

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	_ "github.com/lib/pq"
)

// ConnParams holds connection settings for pg_dump/pg_restore.
type ConnParams struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
	SSLMode  string
}

// DumpResult holds metadata about a successful dump stream setup.
type DumpResult struct {
	ServerVersion string
	ServerMajor   int
	BinaryPath    string
}

// PostgresEngine runs official pg_dump / pg_restore binaries.
type PostgresEngine struct {
	BinRoots []string
}

func NewPostgresEngine() *PostgresEngine {
	return &PostgresEngine{
		BinRoots: []string{
			"/usr/lib/postgresql",
			"/usr/pgsql",
			"/usr/bin",
		},
	}
}

// Dump streams a custom-format dump to w. Caller must consume the reader fully.
// Prefer DumpTo for typical use.
func (e *PostgresEngine) DumpTo(ctx context.Context, p ConnParams, w io.Writer) (*DumpResult, error) {
	major, version, err := e.serverMajor(ctx, p)
	if err != nil {
		return nil, err
	}
	bin, err := e.resolveBinary("pg_dump", major)
	if err != nil {
		return nil, err
	}

	dsn := buildDSN(p)
	cmd := exec.CommandContext(ctx, bin,
		"--dbname="+dsn,
		"--format=custom",
		"--compress=9",
		"--lock-wait-timeout=60000",
		"--no-password",
		"--verbose",
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+p.Password)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr := &boundedBuffer{max: 64 << 10}
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start pg_dump: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(w, stdout)
		waitErr := cmd.Wait()
		if copyErr != nil {
			done <- copyErr
			return
		}
		done <- waitErr
	}()

	select {
	case <-ctx.Done():
		killProcessGroup(cmd)
		<-done
		return nil, fmt.Errorf("pg_dump cancelled: %w", ctx.Err())
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("pg_dump failed: %w\nstderr: %s", err, stderr.String())
		}
	}

	return &DumpResult{
		ServerVersion: version,
		ServerMajor:   major,
		BinaryPath:    bin,
	}, nil
}

// Verify runs pg_restore --list on a dump file path.
func (e *PostgresEngine) Verify(ctx context.Context, dumpPath string, serverMajor int) error {
	bin, err := e.resolveBinary("pg_restore", serverMajor)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, bin, "--list", dumpPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stderr := &boundedBuffer{max: 64 << 10}
	cmd.Stderr = stderr
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("pg_restore --list failed: %w\nstderr: %s", err, stderr.String())
	}
	if len(bytes.TrimSpace(out)) == 0 {
		return fmt.Errorf("pg_restore --list produced empty listing")
	}
	return nil
}

// Restore applies a custom-format dump from dumpPath into the target database.
func (e *PostgresEngine) Restore(ctx context.Context, p ConnParams, dumpPath string, jobs int) error {
	major, _, err := e.serverMajor(ctx, p)
	if err != nil {
		return err
	}
	bin, err := e.resolveBinary("pg_restore", major)
	if err != nil {
		return err
	}
	if jobs < 1 {
		jobs = 1
	}
	dsn := buildDSN(p)
	cmd := exec.CommandContext(ctx, bin,
		"--dbname="+dsn,
		"--clean",
		"--if-exists",
		"--no-owner",
		"--no-acl",
		"--jobs="+strconv.Itoa(jobs),
		"--verbose",
		"--no-password",
		dumpPath,
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+p.Password)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stderr := &boundedBuffer{max: 256 << 10}
	cmd.Stderr = stderr
	cmd.Stdout = io.Discard

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start pg_restore: %w", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-ctx.Done():
		killProcessGroup(cmd)
		<-done
		return fmt.Errorf("pg_restore cancelled: %w", ctx.Err())
	case err := <-done:
		if err != nil {
			return fmt.Errorf("pg_restore failed: %w\nstderr: %s", err, stderr.String())
		}
	}
	return nil
}

// ServerVersion returns the remote PostgreSQL version string and major number.
func (e *PostgresEngine) ServerVersion(ctx context.Context, p ConnParams) (string, int, error) {
	major, version, err := e.serverMajor(ctx, p)
	return version, major, err
}

func (e *PostgresEngine) serverMajor(ctx context.Context, p ConnParams) (int, string, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.Username, p.Password, p.Name, p.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, "", fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return 0, "", fmt.Errorf("ping db: %w", err)
	}

	var versionNum int
	if err := db.QueryRowContext(ctx, "SHOW server_version_num").Scan(&versionNum); err != nil {
		return 0, "", fmt.Errorf("server_version_num: %w", err)
	}
	var version string
	if err := db.QueryRowContext(ctx, "SHOW server_version").Scan(&version); err != nil {
		return 0, "", fmt.Errorf("server_version: %w", err)
	}
	major := versionNum / 10000
	return major, version, nil
}

func (e *PostgresEngine) resolveBinary(name string, serverMajor int) (string, error) {
	// Prefer a client major >= server major.
	candidates := []string{}
	for _, root := range e.BinRoots {
		if root == "/usr/bin" {
			candidates = append(candidates, filepath.Join(root, name))
			continue
		}
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, ent := range entries {
			if !ent.IsDir() {
				continue
			}
			maj, err := strconv.Atoi(ent.Name())
			if err != nil {
				continue
			}
			if maj >= serverMajor {
				candidates = append(candidates, filepath.Join(root, ent.Name(), "bin", name))
			}
		}
	}
	// Also try PATH.
	if p, err := exec.LookPath(name); err == nil {
		candidates = append(candidates, p)
	}

	var lastErr error
	bestPath := ""
	bestMajor := -1
	for _, path := range candidates {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			lastErr = err
			continue
		}
		maj := binaryMajor(path)
		if maj >= serverMajor && maj > bestMajor {
			bestMajor = maj
			bestPath = path
		} else if bestPath == "" && maj == 0 {
			// Unknown major from PATH; accept as last resort after loop.
			if bestPath == "" {
				bestPath = path
			}
		}
	}
	if bestPath != "" {
		return bestPath, nil
	}
	if lastErr != nil {
		return "", fmt.Errorf("no compatible %s for server major %d: %w", name, serverMajor, lastErr)
	}
	return "", fmt.Errorf("no compatible %s found for server major %d (install postgresql-client-%d or newer)", name, serverMajor, serverMajor)
}

func binaryMajor(path string) int {
	// Path like /usr/lib/postgresql/16/bin/pg_dump
	parts := strings.Split(path, string(os.PathSeparator))
	for i, p := range parts {
		if p == "postgresql" || p == "pgsql" {
			if i+1 < len(parts) {
				if maj, err := strconv.Atoi(parts[i+1]); err == nil {
					return maj
				}
			}
		}
	}
	return 0
}

func buildDSN(p ConnParams) string {
	ssl := p.SSLMode
	if ssl == "" {
		ssl = "require"
	}
	return fmt.Sprintf(
		"postgresql://%s@%s:%d/%s?sslmode=%s",
		p.Username, p.Host, p.Port, p.Name, ssl,
	)
}

func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

type boundedBuffer struct {
	max int
	buf bytes.Buffer
}

func (b *boundedBuffer) Write(p []byte) (int, error) {
	remain := b.max - b.buf.Len()
	if remain <= 0 {
		return len(p), nil
	}
	if len(p) > remain {
		p = p[:remain]
	}
	return b.buf.Write(p)
}

func (b *boundedBuffer) String() string {
	return b.buf.String()
}
