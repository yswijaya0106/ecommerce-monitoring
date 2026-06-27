package customer

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
	r.Get("/customers", h.list)
	r.Get("/customers/{id}", h.get)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	customers, err := h.repo.List(r.Context())
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpjson.Write(w, http.StatusOK, customers)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := httpjson.IDFromPath(r)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}
	c, err := h.repo.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httpjson.WriteError(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpjson.Write(w, http.StatusOK, c)
}
