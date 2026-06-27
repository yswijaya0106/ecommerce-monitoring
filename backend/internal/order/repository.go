package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound        = errors.New("order not found")
	ErrEmptyOrder      = errors.New("order must have at least one item")
	ErrProductNotFound = errors.New("one or more products do not exist")
	ErrDiscontinued    = errors.New("one or more products are discontinued")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const summaryQuery = `
	SELECT o.id, o.customer_id, o.order_date, o.status_id, o.notes,
		COALESCE(SUM(od.quantity * od.unit_price * (1 - od.discount)), 0) AS total
	FROM orders o
	LEFT JOIN order_details od ON od.order_id = o.id`

func scanSummary(row interface{ Scan(dest ...any) error }) (*OrderSummary, error) {
	var s OrderSummary
	err := row.Scan(&s.ID, &s.CustomerID, &s.OrderDate, &s.StatusID, &s.Notes, &s.Total)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) List(ctx context.Context) ([]OrderSummary, error) {
	rows, err := r.db.QueryContext(ctx, summaryQuery+` GROUP BY o.id ORDER BY o.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var result []OrderSummary
	for rows.Next() {
		s, err := scanSummary(rows)
		if err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		result = append(result, *s)
	}
	return result, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int64) (*Order, error) {
	row := r.db.QueryRowContext(ctx, summaryQuery+` WHERE o.id = ? GROUP BY o.id`, id)
	summary, err := scanSummary(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT od.id, od.product_id, p.product_name, od.quantity, od.unit_price, od.discount
		FROM order_details od
		LEFT JOIN products p ON p.id = od.product_id
		WHERE od.order_id = ?
		ORDER BY od.id`, id)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.ProductID, &it.ProductName, &it.Quantity, &it.UnitPrice, &it.Discount); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &Order{OrderSummary: *summary, Items: items}, nil
}

// Create places a new order: it resolves each product's current list_price
// inside the transaction, inserts the order header and its line items, and
// returns the persisted order with items.
func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Order, error) {
	if len(req.Items) == 0 {
		return nil, ErrEmptyOrder
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	type resolvedItem struct {
		productID int64
		quantity  float64
		unitPrice float64
	}
	resolved := make([]resolvedItem, 0, len(req.Items))
	for _, item := range req.Items {
		var listPrice float64
		var discontinued bool
		err := tx.QueryRowContext(ctx,
			`SELECT list_price, discontinued FROM products WHERE id = ?`, item.ProductID,
		).Scan(&listPrice, &discontinued)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		if err != nil {
			return nil, fmt.Errorf("lookup product %d: %w", item.ProductID, err)
		}
		if discontinued {
			return nil, ErrDiscontinued
		}
		resolved = append(resolved, resolvedItem{productID: item.ProductID, quantity: item.Quantity, unitPrice: listPrice})
	}

	res, err := tx.ExecContext(ctx,
		`INSERT INTO orders (customer_id, order_date, notes, status_id) VALUES (?, NOW(), ?, 0)`,
		req.CustomerID, req.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("insert order: %w", err)
	}
	orderID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("insert order: %w", err)
	}

	for _, item := range resolved {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO order_details (order_id, product_id, quantity, unit_price, discount) VALUES (?, ?, ?, ?, 0)`,
			orderID, item.productID, item.quantity, item.unitPrice,
		); err != nil {
			return nil, fmt.Errorf("insert order item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit order: %w", err)
	}

	return r.Get(ctx, orderID)
}
