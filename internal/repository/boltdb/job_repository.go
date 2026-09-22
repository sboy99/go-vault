package boltdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/sboy99/go-vault/internal/domain"
	boltclient "github.com/sboy99/go-vault/pkg/boltdb"
)

const bucketJob = "job_meta"

// JobRepository persists Job entities in BoltDB.
type JobRepository struct {
	client *boltclient.Client
}

var _ domain.JobRepository = (*JobRepository)(nil)

func NewJobRepository(client *boltclient.Client) (*JobRepository, error) {
	if client == nil {
		return nil, errors.New("boltdb client is required")
	}
	if err := client.CreateBucket(bucketJob); err != nil {
		return nil, fmt.Errorf("create job bucket: %w", err)
	}
	return &JobRepository{client: client}, nil
}

func (r *JobRepository) Save(ctx context.Context, j *domain.Job) error {
	_ = ctx
	data, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return r.client.Save(bucketJob, j.ID, data)
}

func (r *JobRepository) FindByID(ctx context.Context, id string) (*domain.Job, error) {
	_ = ctx
	data, err := r.client.Get(bucketJob, id)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, domain.ErrJobNotFound
	}
	var j domain.Job
	if err := json.Unmarshal(data, &j); err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *JobRepository) List(ctx context.Context, limit, offset int) ([]*domain.Job, error) {
	_ = ctx
	rows, err := r.client.ListAll(bucketJob)
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Job, 0, len(rows))
	for _, row := range rows {
		var j domain.Job
		if err := json.Unmarshal(row, &j); err != nil {
			return nil, err
		}
		out = append(out, &j)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	if offset > len(out) {
		return []*domain.Job{}, nil
	}
	out = out[offset:]
	if limit > 0 && limit < len(out) {
		out = out[:limit]
	}
	return out, nil
}
