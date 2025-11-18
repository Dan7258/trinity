package handler

import (
	"encoding/json"
	"net/http"
	"trinity/internal/models"
	"trinity/pkg/hash"
	"trinity/pkg/jwt"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.Password, err = hash.HashPassword(user.Password)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.db.CreateUser(user)
	if err != nil {
		jsonError(w, http.StatusConflict, err.Error())
		return
	}

	claims := jwt.Claims{
		ID:    user.ID,
		Login: user.Login,
		Role:  user.Role,
	}
	resp := new(jwt.JwtResponse)
	resp.Token, err = jwt.GenerateToken(claims)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	getUser, err := h.db.GetUserWithPasswordByLogin(user.Login)
	if err != nil {
		jsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if !hash.CmpPasswordAndHash(user.Password, getUser.Password) {
		jsonError(w, http.StatusUnauthorized, "login failed")
		return
	}

	claims := jwt.Claims{
		ID:    getUser.ID,
		Login: user.Login,
		Role:  getUser.Role,
	}
	resp := new(jwt.JwtResponse)
	resp.Token, err = jwt.GenerateToken(claims)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
