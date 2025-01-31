package repository

import (
	"Analytics_Service/internal/domain"
	"Analytics_Service/internal/logger"
	"database/sql"

	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=RepositoryClient
type RepositoryClient interface {
	GetUserActions() ([]domain.UserAction, error)
	GetUserEvents() ([]domain.UserEvent, error)
	GetProductStats() ([]domain.ProductStat, error)
}

type repositoryClient struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewClientRepo(db *sql.DB, logger *logger.Logger) *repositoryClient {
	return &repositoryClient{
		db:     db,
		logger: logger,
	}
}

func (r *repositoryClient) GetProductStats() ([]domain.ProductStat, error) {
	var stats []domain.ProductStat

	query := `SELECT product_id, cart_add_count FROM product_stats`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to execute SQL query",
			zap.String("operation", "sql.ClientRepository.GetProductsStats"),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat domain.ProductStat
		err := rows.Scan(&stat.ProductID, &stat.CartAddCount)
		if err != nil {
			r.logger.Error("Failed to scan",
				zap.String("operation", "sql.ClientRepository.GetProductsStats"),
				zap.Error(err))
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (r *repositoryClient) GetUserActions() ([]domain.UserAction, error) {
	query := `SELECT user_id, product_id, action_type, product_name, category FROM user_actions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []domain.UserAction
	for rows.Next() {
		var action domain.UserAction
		if err := rows.Scan(&action.UserID, &action.ProductID, &action.ActionType, &action.Name, &action.Category); err != nil {
			r.logger.Error("Failed to scan",
				zap.String("operation", "sql.ClientRepository.GetUserActions"),
				zap.Error(err))
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func (r *repositoryClient) GetUserEvents() ([]domain.UserEvent, error) {
	query := `SELECT user_id, event_type, email FROM user_events`
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query",
			zap.String("operation", "sql.ClientRepository.GetUserEvents"),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var events []domain.UserEvent
	for rows.Next() {
		var event domain.UserEvent
		if err := rows.Scan(&event.UserID, &event.EventType, &event.Email); err != nil {
			r.logger.Error("Failed to scan",
				zap.String("operation", "sql.ClientRepository.GetUserEvents"),
				zap.Error(err))
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
