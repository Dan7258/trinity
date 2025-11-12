package handler

import (
	"encoding/json"
	"net/http"
	"trinity/internal/model"
	"trinity/pkg/rsa"
)

type Handler struct {
	db *model.Model
}

func NewHandler(db model.Model) *Handler {
	return &Handler{
		db: &db,
	}
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

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("web/html")).ServeHTTP(w, r)
}

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("web/"))).ServeHTTP(w, r)
}

func (h *Handler) EncodeRSA(w http.ResponseWriter, r *http.Request) {
	data := new(rsa.Data)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	encodeData, err := rsa.EncodeData(data)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	m, err := json.Marshal(encodeData)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func (h *Handler) DecodeRSA(w http.ResponseWriter, r *http.Request) {
	decodeData := new(rsa.EncryptedData)
	err := json.NewDecoder(r.Body).Decode(decodeData)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	data := new(rsa.Data)
	data, err = rsa.DecodeData(decodeData)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	m, err := json.Marshal(data)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}
