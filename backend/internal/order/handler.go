package order

import (
	"encoding/json"
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
	r.Get("/orders", h.list)
	r.Post("/orders", h.create)
	r.Get("/orders/{id}", h.get)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.List(r.Context())
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpjson.Write(w, http.StatusOK, orders)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := httpjson.IDFromPath(r)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}
	o, err := h.repo.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httpjson.WriteError(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpjson.Write(w, http.StatusOK, o)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}
	created, err := h.repo.Create(r.Context(), req)
	switch {
	case errors.Is(err, ErrEmptyOrder), errors.Is(err, ErrProductNotFound), errors.Is(err, ErrDiscontinued):
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	case err != nil:
		httpjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpjson.Write(w, http.StatusCreated, created)
}
