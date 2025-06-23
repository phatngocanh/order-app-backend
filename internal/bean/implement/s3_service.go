package beanimplement

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pna/order-app-backend/internal/bean"
	log "github.com/sirupsen/logrus"
)

type S3Service struct {
	s3Client   *s3.Client
	uploader   *manager.Uploader
	bucketName string
	prefix     string
	enabled    bool
}

func NewS3Service() bean.S3Service {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Errorf("Failed to load AWS config: %v", err)
		return &S3Service{enabled: false}
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Create uploader
	uploader := manager.NewUploader(s3Client)

	// Get configuration from environment variables
	bucketName := os.Getenv("AWS_S3_BUCKET")
	if bucketName == "" {
		log.Error("AWS_S3_BUCKET environment variable is required")
		return &S3Service{enabled: false}
	}

	prefix := os.Getenv("AWS_S3_ORDER_IMAGES_PREFIX")
	if prefix == "" {
		prefix = "order-images/"
	}

	log.Infof("S3 service initialized with bucket: %s, prefix: %s", bucketName, prefix)

	return &S3Service{
		s3Client:   s3Client,
		uploader:   uploader,
		bucketName: bucketName,
		prefix:     prefix,
		enabled:    true,
	}
}

func (s *S3Service) UploadImage(ctx context.Context, file io.Reader, fileName string) (string, error) {
	if !s.enabled {
		return "", fmt.Errorf("S3 service is not enabled - check AWS configuration")
	}

	// Generate unique filename with timestamp
	timestamp := time.Now().Unix()
	ext := filepath.Ext(fileName)
	uniqueFileName := fmt.Sprintf("%s%d%s", filepath.Base(fileName[:len(fileName)-len(ext)]), timestamp, ext)

	// Create the full key (path) for the file
	key := s.prefix + uniqueFileName

	// Upload the file
	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return the S3 URL
	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, key)
	return s3URL, nil
}

func (s *S3Service) DeleteImage(ctx context.Context, imageURL string) error {
	if !s.enabled {
		return fmt.Errorf("S3 service is not enabled - check AWS configuration")
	}

	// Extract the key from the S3 URL
	// URL format: https://bucket-name.s3.amazonaws.com/key
	// We need to extract the key part after the bucket name

	// For now, we'll implement a simple extraction
	// In production, you might want to store the key separately or parse the URL more robustly
	key := imageURL
	if len(imageURL) > len("https://") {
		// Remove the protocol and bucket part
		// This is a simplified approach - you might want to use URL parsing for production
		key = imageURL[len("https://"+s.bucketName+".s3.amazonaws.com/"):]
	}

	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}
