package user

import (
	"context"

	"github.com/mojganchakeri/whatsapp-manager/internal/init/logger"
	"github.com/mojganchakeri/whatsapp-manager/internal/repository/postgres/userConnections"
)

type UserService interface {
	GetUserStatusByUserID(ctx context.Context, userID string) (string, error)
}

type userService struct {
	log            logger.Logger
	userRepository userConnections.UserConnectionRepository
}

func New(
	log logger.Logger,
	userRepository userConnections.UserConnectionRepository,
) UserService {
	return &userService{
		log:            log,
		userRepository: userRepository,
	}
}

func (us *userService) GetUserStatusByUserID(ctx context.Context, userID string) (string, error) {
	return us.userRepository.GetUserStatusByUserID(ctx, userID)
}
