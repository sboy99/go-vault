package job

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDecodeJobList(t *testing.T) {
	now := time.Now().UTC()
	j1 := Job{ID: "job-1", Type: TypeBackup, Status: StatusSucceeded, CreatedAt: now}
	j2 := Job{ID: "job-2", Type: TypeRestore, Status: StatusFailed, CreatedAt: now}
	b1, err := json.Marshal(j1)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := json.Marshal(j2)
	if err != nil {
		t.Fatal(err)
	}

	jobs, err := decodeJobList([][]byte{b1, b2})
	if err != nil {
		t.Fatalf("decodeJobList: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("want 2 jobs, got %d", len(jobs))
	}
	if jobs[0].ID != "job-1" || jobs[1].ID != "job-2" {
		t.Fatalf("unexpected jobs: %+v %+v", jobs[0], jobs[1])
	}
}

func TestDecodeJobListEmpty(t *testing.T) {
	jobs, err := decodeJobList(nil)
	if err != nil {
		t.Fatalf("decodeJobList: %v", err)
	}
	if len(jobs) != 0 {
		t.Fatalf("want empty, got %d", len(jobs))
	}
}

func TestDecodeJobListInvalidJSON(t *testing.T) {
	_, err := decodeJobList([][]byte{[]byte("not-json")})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
