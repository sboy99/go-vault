//go:build integration

package engine_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/sboy99/go-vault/internal/engine"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestDumpDropRestore(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "app",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}
	t.Cleanup(func() { _ = pg.Terminate(ctx) })

	host, err := pg.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := pg.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	params := engine.ConnParams{
		Host:     host,
		Port:     port.Int(),
		Name:     "app",
		Username: "postgres",
		Password: "postgres",
		SSLMode:  "disable",
	}

	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=postgres dbname=app sslmode=disable", host, port.Int())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	deadline := time.Now().Add(30 * time.Second)
	for {
		if err := db.Ping(); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("postgres not ready")
		}
		time.Sleep(500 * time.Millisecond)
	}

	if _, err := db.Exec(`CREATE TABLE items (id serial primary key, name text not null);
		INSERT INTO items (name) VALUES ('alpha'), ('beta');`); err != nil {
		t.Fatal(err)
	}

	eng := engine.NewPostgresEngine()
	dir := t.TempDir()
	dumpPath := filepath.Join(dir, "test.dump")
	f, err := os.Create(dumpPath)
	if err != nil {
		t.Fatal(err)
	}

	dumpCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	res, err := eng.DumpTo(dumpCtx, params, f)
	_ = f.Close()
	if err != nil {
		t.Fatalf("dump: %v", err)
	}
	if err := eng.Verify(dumpCtx, dumpPath, res.ServerMajor); err != nil {
		t.Fatalf("verify: %v", err)
	}

	if _, err := db.Exec(`DROP TABLE items;`); err != nil {
		t.Fatal(err)
	}

	if err := eng.Restore(dumpCtx, params, dumpPath, 2); err != nil {
		t.Fatalf("restore: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT count(*) FROM items`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected 2 rows after restore, got %d", count)
	}
}
