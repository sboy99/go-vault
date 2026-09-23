package engine

import (
	"bufio"
	"bytes"
	"compress/gzip"
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
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/sboy99/go-vault/internal/domain"
)

// PostgresEngine runs official pg_dump / psql binaries.
type PostgresEngine struct {
	BinRoots []string
}

var _ domain.DumpEngine = (*PostgresEngine)(nil)

var versionMajorRE = regexp.MustCompile(`\(PostgreSQL\)\s+(\d+)`)

const (
	dumpHeaderMarker   = "PostgreSQL database dump"
	brandingMarker     = "go-vault backup"
	dumpBrandingBanner = "--\n-- go-vault backup\n--\n"
)

func NewPostgresEngine() *PostgresEngine {
	return &PostgresEngine{
		BinRoots: []string{
			"/usr/lib/postgresql",
			"/usr/pgsql",
			"/usr/bin",
		},
	}
}

// DumpTo streams a gzip-compressed plain-SQL dump to w.
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
		"--format=plain",
		"--clean",
		"--if-exists",
		"--no-owner",
		"--no-acl",
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
		gz := gzip.NewWriter(w)
		_, bannerErr := io.WriteString(gz, dumpBrandingBanner)
		var copyErr error
		if bannerErr == nil {
			_, copyErr = io.Copy(gz, stdout)
		}
		closeErr := gz.Close()
		waitErr := cmd.Wait()
		if bannerErr != nil {
			done <- bannerErr
			return
		}
		if copyErr != nil {
			done <- copyErr
			return
		}
		if closeErr != nil {
			done <- closeErr
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

// Verify checks that dumpPath is a valid gzip-compressed PostgreSQL plain SQL dump.
func (e *PostgresEngine) Verify(ctx context.Context, dumpPath string, serverMajor int) error {
	_ = ctx
	_ = serverMajor
	return verifyGzipSQLDump(dumpPath)
}

// Restore applies a gzip-compressed plain-SQL dump from dumpPath into the target database.
// jobs is ignored; plain SQL restores are single-threaded via psql.
func (e *PostgresEngine) Restore(ctx context.Context, p domain.ConnParams, dumpPath string, jobs int) error {
	_ = jobs
	if _, err := os.Stat(dumpPath); err != nil {
		return fmt.Errorf("dump file: %w", err)
	}
	major, _, err := e.serverMajor(ctx, p)
	if err != nil {
		return err
	}
	psql, err := e.resolveBinary("psql", major)
	if err != nil {
		return err
	}

	f, err := os.Open(dumpPath)
	if err != nil {
		return fmt.Errorf("open dump: %w", err)
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("gzip open: %w", err)
	}

	filtered, waitFilter := filterTransactionTimeout(gz)
	defer func() {
		_ = filtered.Close()
		waitFilter()
		_ = gz.Close()
	}()

	dsn := buildDSN(p)
	cmd := exec.CommandContext(ctx, psql,
		"--dbname="+dsn,
		"--set=ON_ERROR_STOP=1",
		"--no-password",
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+p.Password)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stderr := &boundedBuffer{max: 256 << 10}
	cmd.Stderr = stderr
	cmd.Stdout = io.Discard
	cmd.Stdin = filtered

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start psql: %w", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	return waitOrKill(ctx, cmd, done, "psql", stderr)
}

// ServerVersion returns the remote PostgreSQL version string and major number.
func (e *PostgresEngine) ServerVersion(ctx context.Context, p domain.ConnParams) (string, int, error) {
	major, version, err := e.serverMajor(ctx, p)
	return version, major, err
}

func verifyGzipSQLDump(dumpPath string) error {
	f, err := os.Open(dumpPath)
	if err != nil {
		return fmt.Errorf("open dump: %w", err)
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("gzip open: %w", err)
	}

	header := make([]byte, 512)
	n, err := io.ReadFull(gz, header)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		_ = gz.Close()
		return fmt.Errorf("read dump header: %w", err)
	}
	if !bytes.Contains(header[:n], []byte(brandingMarker)) {
		_ = gz.Close()
		return fmt.Errorf("missing go-vault branding header")
	}
	if !bytes.Contains(header[:n], []byte(dumpHeaderMarker)) {
		_ = gz.Close()
		return fmt.Errorf("not a PostgreSQL plain SQL dump")
	}

	if _, err := io.Copy(io.Discard, gz); err != nil {
		_ = gz.Close()
		return fmt.Errorf("gzip verify: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("gzip close: %w", err)
	}
	return nil
}

func (e *PostgresEngine) serverMajor(ctx context.Context, p domain.ConnParams) (int, string, error) {
	ssl := p.SSLMode
	if ssl == "" {
		ssl = "require"
	}
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.Username, p.Password, p.Name, ssl,
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
// Long lines (COPY rows) that exceed the read buffer are passed through unchanged.
// The returned wait function blocks until the filter goroutine has exited; call it
// after closing the PipeReader and before closing r.
func filterTransactionTimeout(r io.Reader) (*io.PipeReader, func()) {
	pr, pw := io.Pipe()
	var wg sync.WaitGroup
	wg.Go(func() {
		defer func() { _ = pw.Close() }()
		br := bufio.NewReaderSize(r, 64*1024)
		for {
			chunk, err := br.ReadSlice('\n')
			if len(chunk) > 0 {
				if err == bufio.ErrBufferFull {
					// Fragment of a long line — never a short SET statement. Pass through.
					if _, werr := pw.Write(chunk); werr != nil {
						_ = pw.CloseWithError(werr)
						return
					}
				} else {
					trimmed := strings.TrimSpace(string(chunk))
					skip := strings.HasPrefix(trimmed, "SET transaction_timeout =") ||
						strings.HasPrefix(trimmed, "set transaction_timeout =")
					if !skip {
						if _, werr := pw.Write(chunk); werr != nil {
							_ = pw.CloseWithError(werr)
							return
						}
					}
				}
			}
			if err == bufio.ErrBufferFull {
				continue
			}
			if err == io.EOF {
				return
			}
			if err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
	})
	return pr, wg.Wait
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
