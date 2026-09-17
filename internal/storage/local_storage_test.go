package storage_test

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"testing"

	"github.com/sboy99/go-vault/internal/storage"
)

func TestLocalStorageSaveOpenDelete(t *testing.T) {
	dir := t.TempDir()
	store := storage.NewLocalStorage(filepath.Join(dir, "backups"))
	ctx := context.Background()

	info, err := store.Save(ctx, "test.dump", bytes.NewReader([]byte("hello-backup")))
	if err != nil {
		t.Fatal(err)
	}
	if info.Size != 12 {
		t.Fatalf("size=%d", info.Size)
	}
	if info.SHA256 == "" {
		t.Fatal("missing sha256")
	}

	rc, err := store.Open(ctx, "test.dump")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(rc)
	_ = rc.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello-backup" {
		t.Fatalf("got %q", data)
	}

	list, err := store.List(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("list=%v", list)
	}

	if err := store.Delete(ctx, "test.dump"); err != nil {
		t.Fatal(err)
	}
}
