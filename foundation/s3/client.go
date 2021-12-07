package s3

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/file-share/business/service/storage"
)

var _ storage.Storage = (*MinioClient)(nil)

const (
	bucketName          = "files"
	presignedExpiration = time.Hour
	ctxTimeout          = time.Second * 10
)

// MinioClient S3 storage client. Have to satisfy the storage.Storage interface
type MinioClient struct {
	client *minio.Client

	logger *zap.Logger

	close chan struct{}
}

func NewMinioClient(endpoint, accessTokenID, secretKeyAccessKey, region string, useSSL bool, logger *zap.Logger) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessTokenID, secretKeyAccessKey, ""),
		Region: region,
		Secure: useSSL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "s3 connection error")
	}
	isExists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		logger.Error("error checking bucket", zap.Error(err))
		return nil, err
	}
	if !isExists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			logger.Error("error creating bucket", zap.Error(err))
			return nil, err
		}
	}
	return &MinioClient{client: client, logger: logger, close: make(chan struct{})}, nil
}

func (m *MinioClient) GetUploadURL(ctx context.Context, id string) (string, error) {
	link, err := m.client.PresignedPutObject(ctx, bucketName, id, presignedExpiration)
	if err != nil {
		return "", err
	}
	return link.String(), nil
}

func (m *MinioClient) GetDownloadURL(ctx context.Context, id, name string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", name)
	link, err := m.client.PresignedGetObject(ctx, bucketName, id, presignedExpiration, reqParams)
	if err != nil {
		return "", err
	}
	return link.String(), nil
}

// HandleUploadNotification receive and handle notifications from S3 storage (warning! locking method)
func (m *MinioClient) HandleUploadNotification(h func(ctx context.Context, id string) error) {
	msgs := m.client.ListenBucketNotification(context.Background(), bucketName, "", "", []string{"s3:ObjectCreated:*"})
	for {
		select {
		case <-m.close:
			return
		case info := <-msgs:
			for _, msg := range info.Records {
				ctx, closefn := context.WithTimeout(context.Background(), ctxTimeout)
				m.logger.Info(
					"notification received",
					zap.String("s3-event", msg.EventName),
					zap.String("s3-object-key", msg.S3.Object.Key),
				)
				if err := h(ctx, msg.S3.Object.Key); err != nil {
					m.logger.Warn("error handling upload notification", zap.String("object-name", msg.S3.Object.Key))
				}
				closefn()
			}
		}
	}
}

func (m *MinioClient) Close() {
	close(m.close)
}
