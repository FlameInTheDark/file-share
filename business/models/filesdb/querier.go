// Code generated by sqlc. DO NOT EDIT.

package filesdb

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	Create(ctx context.Context, fileName string) (File, error)
	Find(ctx context.Context, fileID uuid.UUID) (File, error)
	Increase(ctx context.Context, id int32) error
	Stats(ctx context.Context, fileID uuid.UUID) (int64, error)
	Uploaded(ctx context.Context, fileID uuid.UUID) (File, error)
}

var _ Querier = (*Queries)(nil)
