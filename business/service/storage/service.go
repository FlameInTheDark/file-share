package storage

import (
	"context"
)

type Storage interface {
	GetUploadURL(ctx context.Context, id string) (string, error)
	GetDownloadURL(ctx context.Context, id, name string) (string, error)
}

// Service contains storage service interface and prepares data to create an upload and download URLs
type Service struct {
	storage Storage
}

func NewStorageService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// Upload create an upload URL
func (s *Service) Upload(ctx context.Context, id string) (string, error) {
	url, err := s.storage.GetUploadURL(ctx, id)
	return url, err
}

// Download create a download URL
func (s *Service) Download(ctx context.Context, id, name string) (string, error) {
	return s.storage.GetDownloadURL(ctx, id, name)
}
