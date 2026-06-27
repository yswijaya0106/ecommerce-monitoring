package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("product not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const baseColumns = `id, product_code, product_name, description, standard_cost, list_price,
	reorder_level, target_level, quantity_per_unit, discontinued, minimum_reorder_quantity, category`

func scanProduct(row interface{ Scan(dest ...any) error }) (*Product, error) {
	var p Product
	err := row.Scan(
		&p.ID, &p.ProductCode, &p.ProductName, &p.Description, &p.StandardCost, &p.ListPrice,
		&p.ReorderLevel, &p.TargetLevel, &p.QuantityPerUnit, &p.Discontinued, &p.MinimumReorderQuantity, &p.Category,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List returns products, optionally filtered by category. Pass an empty
// string to skip the filter.
func (r *Repository) List(ctx context.Context, category string) ([]Product, error) {
	query := `SELECT ` + baseColumns + ` FROM products`
	args := []any{}
	if category != "" {
		query += ` WHERE category = ?`
		args = append(args, category)
	}
	query += ` ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var result []Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		result = append(result, *p)
	}
	return result, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int64) (*Product, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+baseColumns+` FROM products WHERE id = ?`, id)
	p, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return p, nil
}
