package engine

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/sboy99/go-vault/internal/domain"
)

// PostgresEngine runs official pg_dump / pg_restore binaries.
type PostgresEngine struct {
	BinRoots []string
}

var _ domain.DumpEngine = (*PostgresEngine)(nil)

var versionMajorRE = regexp.MustCompile(`\(PostgreSQL\)\s+(\d+)`)

func NewPostgresEngine() *PostgresEngine {
	return &PostgresEngine{
		BinRoots: []string{
			"/usr/lib/postgresql",
			"/usr/pgsql",
			"/usr/bin",
		},
	}
}

// DumpTo streams a custom-format dump to w.
func (e *PostgresEngine) DumpTo(ctx context.Context, p domain.ConnParams, w io.Writer) (*domain.DumpResult, error) {
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

	if err := waitOrKill(ctx, cmd, done, "pg_dump", stderr); err != nil {
		return nil, err
	}

	return &domain.DumpResult{
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
	return listArchive(ctx, bin, dumpPath)
}

// Restore applies a custom-format dump from dumpPath into the target database.
func (e *PostgresEngine) Restore(ctx context.Context, p domain.ConnParams, dumpPath string, jobs int) error {
	if _, err := os.Stat(dumpPath); err != nil {
		return fmt.Errorf("dump file: %w", err)
	}
	major, _, err := e.serverMajor(ctx, p)
	if err != nil {
		return err
	}
	matching, err := e.resolveBinary("pg_restore", major)
	if err != nil {
		return err
	}
	if jobs < 1 {
		jobs = 1
	}

	if err := listArchive(ctx, matching, dumpPath); err == nil {
		return e.directRestore(ctx, p, dumpPath, matching, jobs)
	}

	newer, err := e.resolveReadableRestore(ctx, dumpPath, major)
	if err != nil {
		return fmt.Errorf("archive unreadable by pg_restore %d and no newer client can list it: %w", major, err)
	}
	return e.compatRestore(ctx, p, dumpPath, newer, major)
}

func (e *PostgresEngine) directRestore(ctx context.Context, p domain.ConnParams, dumpPath, bin string, jobs int) error {
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

	return waitOrKill(ctx, cmd, done, "pg_restore", stderr)
}

// compatRestore converts a newer-format dump to SQL with a newer pg_restore,
// strips SET transaction_timeout (unsupported on servers < 17), and loads via psql.
func (e *PostgresEngine) compatRestore(ctx context.Context, p domain.ConnParams, dumpPath, pgRestore string, serverMajor int) error {
	psql, err := e.resolveBinary("psql", serverMajor)
	if err != nil {
		return err
	}

	restoreCmd := exec.CommandContext(ctx, pgRestore,
		"--clean",
		"--if-exists",
		"--no-owner",
		"--no-acl",
		"--verbose",
		"--no-password",
		"-f", "-",
		dumpPath,
	)
	restoreCmd.Env = append(os.Environ(), "PGPASSWORD="+p.Password)
	restoreCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	restoreErr := &boundedBuffer{max: 256 << 10}
	restoreCmd.Stderr = restoreErr

	stdout, err := restoreCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("pg_restore stdout pipe: %w", err)
	}

	filtered := filterTransactionTimeout(stdout)
	defer func() { _ = filtered.Close() }()

	dsn := buildDSN(p)
	psqlCmd := exec.CommandContext(ctx, psql,
		"--dbname="+dsn,
		"--set=ON_ERROR_STOP=1",
		"--no-password",
	)
	psqlCmd.Env = append(os.Environ(), "PGPASSWORD="+p.Password)
	psqlCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	psqlErr := &boundedBuffer{max: 256 << 10}
	psqlCmd.Stderr = psqlErr
	psqlCmd.Stdout = io.Discard
	psqlCmd.Stdin = filtered

	if err := restoreCmd.Start(); err != nil {
		return fmt.Errorf("start pg_restore (compat): %w", err)
	}
	if err := psqlCmd.Start(); err != nil {
		killProcessGroup(restoreCmd)
		_ = restoreCmd.Wait()
		return fmt.Errorf("start psql (compat): %w", err)
	}

	doneRestore := make(chan error, 1)
	donePsql := make(chan error, 1)
	go func() { doneRestore <- restoreCmd.Wait() }()
	go func() { donePsql <- psqlCmd.Wait() }()

	var restoreWaitErr, psqlWaitErr error
	gotRestore, gotPsql := false, false
	for !gotRestore || !gotPsql {
		var restoreCh <-chan error
		var psqlCh <-chan error
		if !gotRestore {
			restoreCh = doneRestore
		}
		if !gotPsql {
			psqlCh = donePsql
		}
		select {
		case <-ctx.Done():
			_ = filtered.Close()
			killProcessGroup(restoreCmd)
			killProcessGroup(psqlCmd)
			if !gotRestore {
				<-doneRestore
			}
			if !gotPsql {
				<-donePsql
			}
			return fmt.Errorf("compat restore cancelled: %w", ctx.Err())
		case err := <-restoreCh:
			restoreWaitErr = err
			gotRestore = true
			if err != nil {
				_ = filtered.Close()
				killProcessGroup(psqlCmd)
			}
		case err := <-psqlCh:
			psqlWaitErr = err
			gotPsql = true
			_ = filtered.Close()
			if !gotRestore {
				killProcessGroup(restoreCmd)
			}
		}
	}

	if restoreWaitErr != nil {
		return fmt.Errorf("pg_restore (compat) failed: %w\nstderr: %s", restoreWaitErr, restoreErr.String())
	}
	if psqlWaitErr != nil {
		return fmt.Errorf("psql (compat) failed: %w\nstderr: %s", psqlWaitErr, psqlErr.String())
	}
	return nil
}

// ServerVersion returns the remote PostgreSQL version string and major number.
func (e *PostgresEngine) ServerVersion(ctx context.Context, p domain.ConnParams) (string, int, error) {
	major, version, err := e.serverMajor(ctx, p)
	return version, major, err
}

func (e *PostgresEngine) serverMajor(ctx context.Context, p domain.ConnParams) (int, string, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.Username, p.Password, p.Name, p.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, "", fmt.Errorf("open db: %w", err)
	}
	defer func() { _ = db.Close() }()

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
	candidates := e.collectCandidates(name, serverMajor)
	return pickBest(name, serverMajor, candidates)
}

// resolveReadableRestore finds a pg_restore newer than serverMajor that can list dumpPath.
func (e *PostgresEngine) resolveReadableRestore(ctx context.Context, dumpPath string, serverMajor int) (string, error) {
	candidates := e.collectCandidates("pg_restore", serverMajor+1)
	var lastErr error
	bestPath := ""
	bestMajor := -1
	for _, path := range candidates {
		info, err := os.Stat(path)
		if err != nil {
			lastErr = err
			continue
		}
		if info.IsDir() {
			continue
		}
		maj := binaryMajor(path)
		if maj <= serverMajor {
			continue
		}
		if err := listArchive(ctx, path, dumpPath); err != nil {
			lastErr = err
			continue
		}
		if bestMajor < 0 || maj < bestMajor {
			bestMajor = maj
			bestPath = path
		}
	}
	if bestPath != "" {
		return bestPath, nil
	}
	if lastErr != nil {
		return "", lastErr
	}
	return "", fmt.Errorf("no newer pg_restore found that can read the archive")
}

func (e *PostgresEngine) collectCandidates(name string, minMajor int) []string {
	var candidates []string
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
			if maj >= minMajor {
				candidates = append(candidates, filepath.Join(root, ent.Name(), "bin", name))
			}
		}
	}
	if p, err := exec.LookPath(name); err == nil {
		candidates = append(candidates, p)
	}
	return candidates
}

func pickBest(name string, serverMajor int, candidates []string) (string, error) {
	var lastErr error
	exactPath := ""
	closestPath := ""
	closestMajor := -1

	for _, path := range candidates {
		info, err := os.Stat(path)
		if err != nil {
			lastErr = err
			continue
		}
		if info.IsDir() {
			continue
		}
		maj := binaryMajor(path)
		if maj < serverMajor {
			continue
		}
		if maj == serverMajor {
			if exactPath == "" {
				exactPath = path
			}
			continue
		}
		// maj > serverMajor: keep the closest (lowest) newer major
		if closestMajor < 0 || maj < closestMajor {
			closestMajor = maj
			closestPath = path
		}
	}

	if exactPath != "" {
		return exactPath, nil
	}
	if closestPath != "" {
		return closestPath, nil
	}
	if lastErr != nil {
		return "", fmt.Errorf("no compatible %s for server major %d: %w", name, serverMajor, lastErr)
	}
	return "", fmt.Errorf("no compatible %s found for server major %d (install postgresql-client-%d or newer)", name, serverMajor, serverMajor)
}

func listArchive(ctx context.Context, bin, dumpPath string) error {
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

func waitOrKill(ctx context.Context, cmd *exec.Cmd, done <-chan error, name string, stderr *boundedBuffer) error {
	select {
	case <-ctx.Done():
		killProcessGroup(cmd)
		<-done
		return fmt.Errorf("%s cancelled: %w", name, ctx.Err())
	case err := <-done:
		if err != nil {
			return fmt.Errorf("%s failed: %w\nstderr: %s", name, err, stderr.String())
		}
		return nil
	}
}

func binaryMajor(path string) int {
	if maj := majorFromPath(path); maj > 0 {
		return maj
	}
	return majorFromVersion(path)
}

func majorFromPath(path string) int {
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

func majorFromVersion(path string) int {
	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return 0
	}
	m := versionMajorRE.FindSubmatch(out)
	if len(m) < 2 {
		return 0
	}
	maj, err := strconv.Atoi(string(m[1]))
	if err != nil {
		return 0
	}
	return maj
}

// filterTransactionTimeout drops SET transaction_timeout lines (PostgreSQL 17+ only).
func filterTransactionTimeout(r io.Reader) *io.PipeReader {
	pr, pw := io.Pipe()
	go func() {
		defer func() { _ = pw.Close() }()
		scanner := bufio.NewScanner(r)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 16*1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "SET transaction_timeout =") ||
				strings.HasPrefix(trimmed, "set transaction_timeout =") {
				continue
			}
			if _, err := io.WriteString(pw, line+"\n"); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
		if err := scanner.Err(); err != nil {
			_ = pw.CloseWithError(err)
		}
	}()
	return pr
}

func buildDSN(p domain.ConnParams) string {
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
