package serviceimplement

import (
	"context"
	"io"

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
	// Upload image to S3
	imageURL, err := s.s3Service.UploadImage(ctx, file, fileName)
	if err != nil {
		log.Error("OrderImageService.UploadImage Error uploading to S3: " + err.Error())
		return model.UploadOrderImageResponse{}, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Create order image entity
	orderImage := &entity.OrderImage{
		OrderID:  orderID,
		ImageURL: imageURL,
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

	// Convert to response model
	var responseImages []model.OrderImage
	if len(orderImages) > 0 {
		for _, img := range orderImages {
			responseImages = append(responseImages, model.OrderImage{
				ID:       img.ID,
				OrderID:  img.OrderID,
				ImageURL: img.ImageURL,
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
	// First, get the image record to retrieve the S3 URL
	orderImage, err := s.orderImageRepo.GetOneByIDQuery(ctx, imageID, nil)
	if err != nil {
		log.Error("OrderImageService.DeleteImage Error getting image record: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if orderImage == nil {
		return error_utils.ErrorCode.NOT_FOUND
	}

	// Delete from S3
	err = s.s3Service.DeleteImage(ctx, orderImage.ImageURL)
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
