package usecase

import (
	"context"

	"github.com/SilentPlaces/basicauth.git/internal/application/port"
	"github.com/SilentPlaces/basicauth.git/internal/dto/user"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
)

type UserUseCase struct {
	userService port.UserReader
	logger      appLogger.Logger
}

func NewUserUseCase(userService port.UserReader, logger appLogger.Logger) *UserUseCase {
	return &UserUseCase{userService: userService, logger: logger}
}

func (u *UserUseCase) GetUserByID(ctx context.Context, userID string) (*dto.UserResponseDTO, error) {
	u.logger.Info(ctx, "user fetch requested", map[string]interface{}{"user_id": userID})
	return u.userService.GetUser(userID)
}
