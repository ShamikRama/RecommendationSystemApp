package repository

import (
	"Product_Service/internal/domain"
	"Product_Service/internal/erorrs"
	"Product_Service/internal/logger"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=RepositoryProduct
type RepositoryProduct interface {
	CreateCartItem(ctx context.Context, cart *domain.Cart) error
	UpdateCartItem(ctx context.Context, cart *domain.Cart) error
	DeleteCartItem(ctx context.Context, userID, productID int) error
	GetAllProducts(ctx context.Context, pageSize int, page int) ([]domain.Product, int, error)
}

type repositoryProduct struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewProductRepository(db *sql.DB, logger *logger.Logger) *repositoryProduct {
	return &repositoryProduct{
		db:     db,
		logger: logger,
	}
}

func (r *repositoryProduct) CreateCartItem(ctx context.Context, cart *domain.Cart) error {
	query := `INSERT INTO cart(user_id, product_id, quantity) VALUES ($1,$2,$3)`

	r.logger.Info("Executing SQL query",
		zap.String("query", query),
		zap.Int("userID", cart.UserId),
		zap.Int("productID", cart.ProductId),
		zap.Int("quantity", cart.Quantity))

	_, err := r.db.ExecContext(ctx, query, cart.UserId, cart.ProductId, cart.Quantity)
	if err != nil {
		r.logger.Error("Error executing", zap.String("operation", "sql.Product.CreateCartItem"), zap.Error(err))
		return err
	}

	r.logger.Info("Cart item added successfully", zap.Int("user_id", cart.UserId), zap.Int("product_id", cart.ProductId))
	return nil
}

func (r *repositoryProduct) UpdateCartItem(ctx context.Context, cart *domain.Cart) error {
	query := `UPDATE carts SET quantity = $1 WHERE user_id = $2 AND product_id = $3`

	r.logger.Info("Executing SQL query", zap.String("query", query), zap.Any("params", []interface{}{
		cart.Quantity, cart.UserId, cart.ProductId}))

	result, err := r.db.ExecContext(ctx, query, cart.Quantity, cart.UserId, cart.ProductId)
	if err != nil {
		r.logger.Error("Error executing",
			zap.String("operation", "sql.Product.UpdateCartItem"),
			zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Error affecting rows", zap.String("operation", "sql.Product.UpdateCartItem"), zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		r.logger.Error("Cart item for user not found",
			zap.String("operation", "sql.Product.UpdateCartItem"),
			zap.Int("userId", cart.UserId),
			zap.Int("productId", cart.ProductId),
			zap.Error(err),
		)
		return erorrs.ErrRowsNull
	}

	r.logger.Info("Cart item updated successfully", zap.Int("user_id", cart.UserId), zap.Int("product_id", cart.ProductId))
	return nil
}

func (r *repositoryProduct) DeleteCartItem(ctx context.Context, userID, productID int) error {
	query := `DELETE FROM carts WHERE user_id = $1 AND product_id = $2`

	r.logger.Info("Executing SQL query", zap.String("query", query), zap.Any("params", []interface{}{userID, productID}))

	result, err := r.db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		r.logger.Error("Failed to execute SQL query",
			zap.String("operation", "sql.Product.DeleteCartItem"),
			zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected",
			zap.String("operation", "sql.Product.DeleteCartItem"),
			zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		r.logger.Error("Cart item for user not found",
			zap.String("operation", "sql.Product.UpdateCartItem"),
			zap.Int("userId", userID),
			zap.Int("productId", productID),
			zap.Error(err),
		)
		return erorrs.ErrRowsNull
	}

	r.logger.Info("Cart item deleted successfully",
		zap.Int("user_id", userID),
		zap.Int("product_id", productID))
	return nil
}

func (r *repositoryProduct) GetAllProducts(ctx context.Context, pageSize int, page int) ([]domain.Product, int, error) {
	var total int
	count := `SELECT COUNT(*) FROM products`
	err := r.db.QueryRowContext(ctx, count).Scan(&total)
	if err != nil {
		r.logger.Error("Error to get total products",
			zap.String("operation", "sql.Product.GetTotalProducts"),
			zap.Error(err))
		return nil, 0, err
	}

	query := `SELECT * FROM products`
	var args []interface{}

	if pageSize > 0 && page > 0 {
		offset := (page - 1) * pageSize
		query += ` LIMIT $1 OFFSET $2`
		args = append(args, pageSize, offset)
	}

	r.logger.Info("Executing SQL query", zap.String("query", query), zap.Int("page_size", pageSize), zap.Int("page", page))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Error to execute SQL query",
			zap.String("operation", "sql.Product.GetAllProducts"),
			zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Category, &product.Price); err != nil {
			r.logger.Error("Error to scan the rows",
				zap.String("operation", "sql.Product.GetAllProducts"),
				zap.Error(err))
			return nil, 0, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating the rows",
			zap.String("operation", "sql.Product.GetAllProducts"),
			zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Products retrieved successfully", zap.Int("count", len(products)))
	return products, total, nil
}
