package bean

import (
	"context"
	"io"
	"time"
)

type S3Service interface {
	UploadImage(ctx context.Context, file io.Reader, fileName string) (string, error)
	DeleteImage(ctx context.Context, imageURL string) error
	GenerateSignedUploadURL(ctx context.Context, fileName string, contentType string) (string, string, error)
	GenerateSignedDownloadURL(ctx context.Context, s3Key string, expiresIn time.Duration) (string, error)
}
