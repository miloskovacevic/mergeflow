package domain

import (
	"context"
	"time"
)

type MergeRequestCollection []MergeRequest
type MergeRequest struct {
	ID        string
	Title     string
	Author    string
	Status    string
	CreatedAt time.Time
}

type MergeRequestFilter struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string
	ProjectId int
}

type GitProvider interface {
	GetMergeRequests(ctx context.Context, filter MergeRequestFilter) (MergeRequestCollection, error)
}
