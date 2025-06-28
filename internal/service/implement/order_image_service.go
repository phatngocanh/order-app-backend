package serviceimplement

import (
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/bean"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type OrderImageService struct {
	orderImageRepo repository.OrderImageRepository
	s3Service      bean.S3Service
}

func NewOrderImageService(
	orderImageRepo repository.OrderImageRepository,
	s3Service bean.S3Service,
) service.OrderImageService {
	return &OrderImageService{
		orderImageRepo: orderImageRepo,
		s3Service:      s3Service,
	}
}

func (s *OrderImageService) UploadImage(ctx *gin.Context, orderID int, file io.Reader, fileName string) (model.UploadOrderImageResponse, string) {
	// Upload image to S3 and get the S3 key
	s3Key, err := s.s3Service.UploadImage(ctx, file, fileName)
	if err != nil {
		log.Error("OrderImageService.UploadImage Error uploading to S3: " + err.Error())
		return model.UploadOrderImageResponse{}, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Generate a signed URL for immediate access
	signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, s3Key, 1*time.Hour)
	if err != nil {
		log.Error("OrderImageService.UploadImage Error generating signed URL: " + err.Error())
		return model.UploadOrderImageResponse{}, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Create order image entity
	orderImage := &entity.OrderImage{
		OrderID:  orderID,
		ImageURL: signedURL, // Store signed URL temporarily
		S3Key:    s3Key,     // Store S3 key for future signed URL generation
	}

	// Save to database
	err = s.orderImageRepo.CreateCommand(ctx, orderImage, nil)
	if err != nil {
		log.Error("OrderImageService.UploadImage Error saving to database: " + err.Error())
		return model.UploadOrderImageResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response model
	response := model.UploadOrderImageResponse{
		OrderImage: model.OrderImage{
			ID:       orderImage.ID,
			OrderID:  orderImage.OrderID,
			ImageURL: orderImage.ImageURL,
			S3Key:    orderImage.S3Key,
		},
	}

	return response, ""
}

func (s *OrderImageService) GetImagesByOrderID(ctx context.Context, orderID int) (model.GetOrderImagesResponse, string) {
	// Get images from database
	orderImages, err := s.orderImageRepo.GetAllByOrderIDQuery(ctx, orderID, nil)
	if err != nil {
		log.Error("OrderImageService.GetImagesByOrderID Error: " + err.Error())
		return model.GetOrderImagesResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response model and generate fresh signed URLs
	var responseImages []model.OrderImage
	if len(orderImages) > 0 {
		for _, img := range orderImages {
			// Generate a fresh signed URL for each image
			signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, img.S3Key, 1*time.Hour)
			if err != nil {
				log.Error("OrderImageService.GetImagesByOrderID Error generating signed URL: " + err.Error())
				// Continue with other images even if one fails
				signedURL = ""
			}

			responseImages = append(responseImages, model.OrderImage{
				ID:       img.ID,
				OrderID:  img.OrderID,
				ImageURL: signedURL,
				S3Key:    img.S3Key,
			})
		}
		return model.GetOrderImagesResponse{
			OrderImages: responseImages,
		}, ""
	} else {
		return model.GetOrderImagesResponse{
			OrderImages: make([]model.OrderImage, 0),
		}, ""
	}
}

func (s *OrderImageService) DeleteImage(ctx *gin.Context, imageID int) string {
	// First, get the image record to retrieve the S3 key
	orderImage, err := s.orderImageRepo.GetOneByIDQuery(ctx, imageID, nil)
	if err != nil {
		log.Error("OrderImageService.DeleteImage Error getting image record: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if orderImage == nil {
		return error_utils.ErrorCode.NOT_FOUND
	}

	// Delete from S3 using the S3 key
	err = s.s3Service.DeleteImage(ctx, orderImage.S3Key)
	if err != nil {
		log.Error("OrderImageService.DeleteImage Error deleting from S3: " + err.Error())
		// Continue with database deletion even if S3 deletion fails
		// This prevents orphaned database records
	}

	// Delete from database
	err = s.orderImageRepo.DeleteByIDCommand(ctx, imageID, nil)
	if err != nil {
		log.Error("OrderImageService.DeleteImage Error deleting from database: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	return ""
}

func (s *OrderImageService) GenerateSignedUploadURL(ctx *gin.Context, orderID int, fileName string, contentType string) (model.GenerateSignedUploadURLResponse, string) {
	// Generate signed upload URL and S3 key
	signedURL, s3Key, err := s.s3Service.GenerateSignedUploadURL(ctx, fileName, contentType)
	if err != nil {
		log.Error("OrderImageService.GenerateSignedUploadURL Error generating signed upload URL: " + err.Error())
		return model.GenerateSignedUploadURLResponse{}, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Create order image entity with S3 key (no signed URL yet since file hasn't been uploaded)
	orderImage := &entity.OrderImage{
		OrderID: orderID,
		S3Key:   s3Key,
	}

	// Save to database
	err = s.orderImageRepo.CreateCommand(ctx, orderImage, nil)
	if err != nil {
		log.Error("OrderImageService.GenerateSignedUploadURL Error saving to database: " + err.Error())
		return model.GenerateSignedUploadURLResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	response := model.GenerateSignedUploadURLResponse{
		SignedURL: signedURL,
		S3Key:     s3Key,
	}

	return response, ""
}
