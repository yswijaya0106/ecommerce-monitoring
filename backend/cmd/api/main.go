package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/customer"
	"ecommerce-backend/internal/db"
	"ecommerce-backend/internal/httpserver"
	"ecommerce-backend/internal/order"
	"ecommerce-backend/internal/product"
	"ecommerce-backend/internal/telemetry"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownTelemetry, err := telemetry.Setup(ctx, cfg.ServiceName, cfg.ServiceVersion, cfg.OTLPEndpoint)
	if err != nil {
		log.Fatalf("telemetry setup: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTelemetry(shutdownCtx); err != nil {
			log.Printf("telemetry shutdown: %v", err)
		}
	}()

	sqlDB, err := db.Open(cfg.DSN())
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer sqlDB.Close()

	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := db.WaitForReady(waitCtx, sqlDB); err != nil {
		log.Fatalf("db not ready: %v", err)
	}

	productHandler := product.NewHandler(product.NewRepository(sqlDB))
	customerHandler := customer.NewHandler(customer.NewRepository(sqlDB))
	orderHandler := order.NewHandler(order.NewRepository(sqlDB))
	router := httpserver.New(productHandler, customerHandler, orderHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("orders-api listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
}
