package handler

import (
	"encoding/json"
	"net/http"
	"trinity/internal/models"
	"trinity/internal/repository"
	"trinity/pkg/kuznechik"
	"trinity/pkg/rsa"
	"trinity/pkg/stribog"
)

type Handler struct {
	db models.Model
}

func NewHandler(db models.Model) *Handler {
	h := new(Handler)
	h.db = db
	return h
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": msg,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

type Data struct {
	Message string `json:"message"`
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("web/html")).ServeHTTP(w, r)
}

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("web/"))).ServeHTTP(w, r)
}
