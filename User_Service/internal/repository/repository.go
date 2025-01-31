package repository

import (
	"User_Service/internal/logger"
	"User_Service/internal/repository/auth"
	"User_Service/internal/repository/model"
	"User_Service/internal/repository/user"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repo struct {
	Auth
	User
}

func NewRepository(db *sql.DB, logger *logger.Logger) *Repo {
	return &Repo{
		Auth: auth.NewDBAuth(db, logger),
		User: user.NewUserDB(db, logger),
	}
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Auth
type Auth interface {
	Create(ctx context.Context, user model.User) (int, error)
	GetUser(ctx context.Context, username, password string) (model.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=User
type User interface {
	UpdateUser(ctx context.Context, user model.User) (int, error)
}
