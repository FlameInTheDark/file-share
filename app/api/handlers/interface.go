package handlers

import "context"

type FileService interface {
	Find(ctx context.Context, id string) (string, string, error)
	Create(ctx context.Context, name string) (string, error)
	Uploaded(ctx context.Context, id string) error
	Statistics(ctx context.Context, id string) (int64, error)
}

type StorageService interface {
	Upload(ctx context.Context, id string) (string, error)
	Download(ctx context.Context, id, name string) (string, error)
}
