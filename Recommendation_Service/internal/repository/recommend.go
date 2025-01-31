package repository

import (
	"Recommendation_Service/internal/domain"
	"Recommendation_Service/internal/logger"
	"database/sql"

	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ProductRepository
type ProductRepository interface {
	SaveRecommendation(userId int, productId int, name string, category string) error
	Get(category string, productId int) ([]domain.Product, error)
	SaveDataProduct(productId int, name string, category string) error
	SaveUserActions(userId int, productId int, name string, category string) error
	GetUserProducts(userId int) ([]domain.Product, error)
}

type productRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewProductRepository(db *sql.DB, logger *logger.Logger) *productRepository {
	return &productRepository{
		db:     db,
		logger: logger,
	}
}

func (r *productRepository) Get(category string, productId int) ([]domain.Product, error) {
	query := `SELECT product_id, name, category FROM products WHERE category = $1 AND product_id != $2 LIMIT 3`
	rows, err := r.db.Query(query, category, productId)
	if err != nil {
		r.logger.Error("Error querying",
			zap.String("operation", "sql.Rec.Get"),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ProductId, &p.Name, &p.Category); err != nil {
			r.logger.Error("Error to scan row",
				zap.String("operation", "sql.Rec.Get"),
				zap.Error(err))
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iteration rows",
			zap.String("operation", "sql.Rec.Get"),
			zap.Error(err))
		return nil, err
	}

	return products, nil
}

func (r *productRepository) SaveDataProduct(productId int, name string, category string) error {
	query := `INSERT INTO products(product_id, name, category) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, productId, name, category)
	if err != nil {
		r.logger.Error("Error executing query",
			zap.String("operation", "sql.Rec.SaveDataProducts"),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *productRepository) SaveUserActions(userId int, productId int, name string, category string) error {
	query := `INSERT INTO user_actions(user_id, product_id, name, category) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, userId, productId, name, category)
	if err != nil {
		r.logger.Error("Error executing query",
			zap.String("operation", "sql.Rec.SaveUserActions"),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *productRepository) GetUserProducts(userId int) ([]domain.Product, error) {
	query := `SELECT product_id, name, category FROM user_actions WHERE user_id = $1`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		r.logger.Error("Error querying",
			zap.String("operation", "sql.Rec.GetUserProducts"),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ProductId, &p.Name, &p.Category); err != nil {
			r.logger.Error("Error scanning the row",
				zap.String("operation", "sql.Rec.GetUserProducts"),
				zap.Error(err))
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating the rows",
			zap.String("operation", "sql.Rec.GetUserProducts"),
			zap.Error(err))
		return nil, err
	}

	return products, nil
}

func (r *productRepository) SaveRecommendation(userId int, productId int, name string, category string) error {
	query := `INSERT INTO recommendations(user_id, product_id, name, category) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, userId, productId, name, category)
	if err != nil {
		r.logger.Error("Error executing the query",
			zap.String("operation", "sql.Rec.GetUserProducts"),
			zap.Error(err))
		return err
	}
	return nil
}
