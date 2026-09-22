package engine

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPickBestPrefersExactMajor(t *testing.T) {
	root := t.TempDir()
	path16 := writeFakeBin(t, root, "16", "pg_restore")
	path17 := writeFakeBin(t, root, "17", "pg_restore")

	got, err := pickBest("pg_restore", 16, []string{path17, path16})
	if err != nil {
		t.Fatal(err)
	}
	if got != path16 {
		t.Fatalf("expected exact major 16 %q, got %q", path16, got)
	}
}

func TestPickBestClosestNewerWhenNoExact(t *testing.T) {
	root := t.TempDir()
	path17 := writeFakeBin(t, root, "17", "pg_restore")
	path18 := writeFakeBin(t, root, "18", "pg_restore")

	got, err := pickBest("pg_restore", 16, []string{path18, path17})
	if err != nil {
		t.Fatal(err)
	}
	if got != path17 {
		t.Fatalf("expected closest newer 17 %q, got %q", path17, got)
	}
}

func TestPickBestRejectsOlderOnly(t *testing.T) {
	root := t.TempDir()
	path15 := writeFakeBin(t, root, "15", "pg_restore")

	_, err := pickBest("pg_restore", 16, []string{path15})
	if err == nil {
		t.Fatal("expected error when only older client is available")
	}
}

func TestMajorFromPath(t *testing.T) {
	cases := []struct {
		path string
		want int
	}{
		{"/usr/lib/postgresql/16/bin/pg_restore", 16},
		{"/usr/pgsql/17/bin/pg_dump", 17},
		{"/usr/bin/pg_restore", 0},
	}
	for _, tc := range cases {
		if got := majorFromPath(tc.path); got != tc.want {
			t.Fatalf("majorFromPath(%q)=%d want %d", tc.path, got, tc.want)
		}
	}
}

func TestFilterTransactionTimeout(t *testing.T) {
	in := strings.Join([]string{
		"SET statement_timeout = 0;",
		"SET lock_timeout = 0;",
		"SET idle_in_transaction_session_timeout = 0;",
		"SET transaction_timeout = 0;",
		"SET client_encoding = 'UTF8';",
		"CREATE TABLE t (id int);",
	}, "\n") + "\n"

	outBytes, err := io.ReadAll(filterTransactionTimeout(strings.NewReader(in)))
	if err != nil {
		t.Fatal(err)
	}
	out := string(outBytes)
	if strings.Contains(out, "transaction_timeout") {
		t.Fatalf("filter left transaction_timeout in output:\n%s", out)
	}
	for _, want := range []string{
		"SET statement_timeout = 0;",
		"SET lock_timeout = 0;",
		"SET idle_in_transaction_session_timeout = 0;",
		"SET client_encoding = 'UTF8';",
		"CREATE TABLE t (id int);",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing line %q in:\n%s", want, out)
		}
	}
}

func writeFakeBin(t *testing.T, root, major, name string) string {
	t.Helper()
	dir := filepath.Join(root, "postgresql", major, "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/true\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}
