package service

import (
	"User_Service/internal/logger"
	"User_Service/internal/repository"
	"User_Service/internal/service/auth"
	"User_Service/internal/service/model"
	"User_Service/internal/service/user"
	"User_Service/pkg/kafka"
	"context"
)

type Service struct {
	Auth
	User
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Auth
type Auth interface {
	CreateUser(ctx context.Context, login model.Login) (int, error)
	GenerateJwtToken(ctx context.Context, email string, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=User
type User interface {
	UpdateUser(ctx context.Context, userID int, updateInfo model.UpdateInfoUser) (int, error)
}

func NewService(repo repository.Repo, logger *logger.Logger, kafka kafka.Kafka) *Service {
	return &Service{
		Auth: auth.NewAuthService(repo.Auth, logger, kafka),
		User: user.NewUserService(repo.User, logger, kafka),
	}
}
