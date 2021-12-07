package files

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/FlameInTheDark/file-share/business/models/filesdb"
)

type Service struct {
	files filesdb.Querier
}

func NewFilesService(db *sqlx.DB) *Service {
	return &Service{
		files: filesdb.New(db),
	}
}

func (s *Service) Find(ctx context.Context, id string) (string, string, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return "", "", err
	}
	file, err := s.files.Find(ctx, uid)
	if err != nil {
		return "", "", err
	}
	err = s.files.Increase(ctx, file.ID)
	if err != nil {
		return "", "", err
	}
	return file.FileID.String(), file.FileName, nil
}

func (s *Service) Create(ctx context.Context, name string) (string, error) {
	file, err := s.files.Create(ctx, name)
	if err != nil {
		return "", err
	}
	return file.FileID.String(), err
}

func (s *Service) Uploaded(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	_, err = s.files.Uploaded(ctx, uid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Statistics(ctx context.Context, id string) (int64, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return 0, err
	}
	return s.files.Stats(ctx, uid)
}