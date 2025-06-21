package serviceimplement

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/bean"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/constants"
	"github.com/pna/order-app-backend/internal/utils/env"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	"github.com/pna/order-app-backend/internal/utils/jwt"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	userRepository  repository.UserRepository
	passwordEncoder bean.PasswordEncoder
}

func NewUserService(userRepository repository.UserRepository, passwordEncoder bean.PasswordEncoder) service.UserService {
	return &UserService{
		userRepository:  userRepository,
		passwordEncoder: passwordEncoder,
	}
}

func (s *UserService) Login(ctx *gin.Context, request model.LoginRequest) (*model.LoginResponse, string) {
	// Find user by username
	user, err := s.userRepository.FindByUsernameQuery(ctx, request.Username, nil)
	if err != nil {
		log.Error("UserService.Login Error when get user: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if user == nil {
		return nil, error_utils.ErrorCode.USERNAME_NOT_FOUND
	}

	// Verify password
	isValid := s.passwordEncoder.Compare(user.Password, request.Password)
	if !isValid {
		return nil, error_utils.ErrorCode.UNAUTHORIZED
	}

	// Generate JWT token
	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		log.Error("UserService.Login Error when get JWT secret: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	token, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
	if err != nil {
		log.Error("UserService.Login Error when generate token: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return &model.LoginResponse{
		Token:    token,
		Username: user.Username,
	}, ""
}
