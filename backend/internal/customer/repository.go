package customer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("customer not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const baseColumns = `id, company, last_name, first_name, email_address, job_title,
	business_phone, home_phone, mobile_phone, fax_number, address, city,
	state_province, zip_postal_code, country_region, web_page, notes`

func scanCustomer(row interface{ Scan(dest ...any) error }) (*Customer, error) {
	var c Customer
	err := row.Scan(
		&c.ID, &c.Company, &c.LastName, &c.FirstName, &c.EmailAddress, &c.JobTitle,
		&c.BusinessPhone, &c.HomePhone, &c.MobilePhone, &c.FaxNumber, &c.Address, &c.City,
		&c.StateProvince, &c.ZipPostalCode, &c.CountryRegion, &c.WebPage, &c.Notes,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) List(ctx context.Context) ([]Customer, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+baseColumns+` FROM customers ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	var result []Customer
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		result = append(result, *c)
	}
	return result, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int64) (*Customer, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+baseColumns+` FROM customers WHERE id = ?`, id)
	c, err := scanCustomer(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return c, nil
}
