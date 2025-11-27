package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
	"trinity/internal/models"
	"trinity/internal/repository_redis"
	"trinity/pkg/smtp"
)

type Handler struct {
	db  models.Model
	rdb repository_redis.RedisDB
}

func NewHandler(db models.Model, rdb repository_redis.RedisDB) *Handler {
	h := new(Handler)
	h.db = db
	h.rdb = rdb
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
	http.FileServer(http.Dir("web/html/")).ServeHTTP(w, r)
}
func (h *Handler) GetHistoryPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/admin.html")
}

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("web/"))).ServeHTTP(w, r)
}

func (h *Handler) TwoFA(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	_, err := h.db.GetUserByLogin(email)
	if err != nil {
		jsonError(w, http.StatusNotFound, "User not found")
		return
	}
	code, _ := GenerateRandomNumber()
	err = h.rdb.Set(r.Context(), email, code, time.Minute*10).Err()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = smtp.SendMessage(email, "secret code", code)
	if err != nil {
		h.rdb.Del(r.Context(), email)
		log.Println(err)
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func GenerateRandomNumber() (string, error) {
	bigNum := big.NewInt(100000)
	randomNumber, err := rand.Int(rand.Reader, bigNum)
	if err != nil {
		fmt.Println("Ошибка при генерации случайного числа:", err)
		return "", err
	}
	return fmt.Sprintf("%.6d", randomNumber.Int64()), nil
}
