package handler

import (
	"net/http"
	"trinity/internal/model"
)

type Handler struct {
	db *model.Model
}

func NewHandler(db model.Model) *Handler {
	return &Handler{
		db: &db,
	}
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("web/html")).ServeHTTP(w, r)
}

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("web/"))).ServeHTTP(w, r)
}
