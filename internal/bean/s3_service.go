package bean

import (
	"context"
	"io"
)

type S3Service interface {
	UploadImage(ctx context.Context, file io.Reader, fileName string) (string, error)
	DeleteImage(ctx context.Context, imageURL string) error
}
