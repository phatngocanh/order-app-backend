package serviceimplement

import (
	"context"
	"github.com/pna/order-app-backend/internal/bean"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
)

type HelloWorldService struct {
	helloWorldRepository repository.HelloWorldRepository
	passwordEncoder      bean.PasswordEncoder
}

func NewHelloWorldService(helloWorldRepository repository.HelloWorldRepository, passwordEncoder bean.PasswordEncoder) service.HelloWorldService {
	return &HelloWorldService{
		helloWorldRepository: helloWorldRepository,
		passwordEncoder:      passwordEncoder,
	}
}

func (s HelloWorldService) HelloWorld(ctx context.Context) (*model.HelloWorldResponse, string) {
	helloWorldEntity, err := s.helloWorldRepository.GetHelloWorldQuery(ctx, nil)
	if err != nil {
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	return &model.HelloWorldResponse{
		Message: helloWorldEntity.Message,
	}, ""
}
