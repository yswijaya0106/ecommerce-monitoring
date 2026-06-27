package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Open connects to MySQL with otelsql instrumentation so every query
// produces a span and a set of DB client metrics automatically.
func Open(dsn string) (*sql.DB, error) {
	db, err := otelsql.Open("mysql", dsn,
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{Ping: true}),
	)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if _, err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemMySQL)); err != nil {
		return nil, fmt.Errorf("register db stats metrics: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// WaitForReady pings the database until it answers or the context expires.
// MySQL containers can take a few seconds after "started" before accepting
// connections, so the API retries instead of crashing on boot.
func WaitForReady(ctx context.Context, db *sql.DB) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
