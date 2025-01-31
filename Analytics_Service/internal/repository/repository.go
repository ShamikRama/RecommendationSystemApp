package repository

import (
	"Analytics_Service/internal/logger"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=AnaliticRepository
type AnaliticRepository interface {
	SaveUserEvent(userId uint32, email string, event string) error
	SaveUserAction(userId uint32, productId uint32, action string, productName string, category string) error
	UpdateProductStats(productID uint32, column string) error
}

type analRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewAnalRepository(db *sql.DB, logger *logger.Logger) *analRepository {
	return &analRepository{
		db:     db,
		logger: logger,
	}
}

func (r *analRepository) SaveUserAction(userId uint32, productId uint32, action string, productName string, category string) error {

	query := `
			INSERT INTO user_actions(user_id, product_id, action_type, product_name, category) VALUES($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, userId, productId, action, productName, category)
	if err != nil {
		r.logger.Error("Error to execute SQL query",
			zap.String("operation", "sql.Repository.SaveUserAction"),
			zap.Error(err))
		return err
	}

	return nil
}

func (r *analRepository) SaveUserEvent(userId uint32, email string, event string) error {

	query := `
			INSERT INTO user_events(user_id, event_type, email) VALUES($1, $2, $3)`

	_, err := r.db.Exec(query, userId, event, email)
	if err != nil {
		r.logger.Error("Error to execute SQL query",
			zap.String("operation", "sql.Repository.SaveUserEvent"),
			zap.Error(err))
		return err
	}

	return nil
}

func (r *analRepository) UpdateProductStats(productID uint32, column string) error {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM product_stats WHERE product_id = $1)`, productID).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check if product exist",
			zap.String("operation", "sql.Repository.UpdateProductStats"),
			zap.Error(err))
		return err
	}

	if !exists {
		query := fmt.Sprintf(`INSERT INTO product_stats (product_id, %s) VALUES ($1, 1)`, column)
		_, err := r.db.Exec(query, productID)
		if err != nil {
			r.logger.Error("Error to insert new products SQL query",
				zap.String("operation", "sql.Repository.UpdateProductStats"),
				zap.Error(err))
			return err
		}
	} else {
		query := fmt.Sprintf(`UPDATE product_stats SET %s = %s + 1 WHERE product_id = $1`, column, column)
		_, err := r.db.Exec(query, productID)
		if err != nil {
			r.logger.Error("Failed to update product stats",
				zap.String("operation", "sql.Repository.UpdateProductStats"),
				zap.Error(err))
			return err
		}
	}
	return nil
}
