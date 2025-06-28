package serviceimplement

import (
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
		ImageID:   orderImage.ID,
	}

	return response, ""
}
