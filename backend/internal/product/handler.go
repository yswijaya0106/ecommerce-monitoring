package product

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"ecommerce-backend/internal/httpjson"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/products", h.list)
	r.Get("/products/{id}", h.get)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.List(r.Context(), r.URL.Query().Get("category"))
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpjson.Write(w, http.StatusOK, products)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := httpjson.IDFromPath(r)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}
	p, err := h.repo.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httpjson.WriteError(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpjson.Write(w, http.StatusOK, p)
}
