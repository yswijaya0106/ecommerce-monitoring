package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"ecommerce-backend/internal/customer"
	"ecommerce-backend/internal/order"
	"ecommerce-backend/internal/product"
)

// New builds the application router. Every route is wrapped by otelhttp so
// each HTTP request produces a server span with route, method and status
// attributes, propagating trace context from/to callers automatically.
func New(productHandler *product.Handler, customerHandler *customer.Handler, orderHandler *order.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	productHandler.Routes(r)
	customerHandler.Routes(r)
	orderHandler.Routes(r)

	return otelhttp.NewHandler(r, "http.server",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}

// corsMiddleware allows the Laravel frontend (running on a different
// origin/container) to call this API directly from the browser if needed.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
