package auth

import (
	"User_Service/internal/erorrs"
	"User_Service/internal/logger"
	"User_Service/internal/repository/model"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type DBAuth struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewDBAuth(db *sql.DB, logger *logger.Logger) *DBAuth {
	return &DBAuth{
		db:     db,
		logger: logger,
	}
}

func (r *DBAuth) Create(ctx context.Context, user model.User) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction",
			zap.String("operation", "sql.Auth.CreateUser"),
			zap.Error(err))
		return 0, err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash) VALUES($1, $2) RETURNING id", model.UserTable)

	var id int
	err = tx.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&id)
	if err != nil {
		r.logger.Error("failed to query row",
			zap.String("operation", "sql.Auth.CreateUser"),
			zap.Error(err))
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505 — код ошибки "unique_violation"
			return 0, erorrs.ErrUniqueConstraint
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("Failed to commit the operation",
			zap.String("operation", "sql.Auth.CreateUser"),
			zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (r *DBAuth) GetUser(ctx context.Context, email, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, email, password_hash FROM %s WHERE email = $1 AND password_hash = $2", model.UserTable)

	err := r.db.QueryRowContext(ctx, query, email, password).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		r.logger.Error("failed to query row", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("Error no rows",
				zap.String("operation", "sql.Auth.GetUser"))
			zap.Error(err)
			return user, erorrs.ErrNoRows
		}
		return user, err
	}

	return user, nil
}
