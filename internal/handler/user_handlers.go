package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	user.Role = "user"
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
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	codeEmail, ok := data["code_email"]
	if !ok {
		jsonError(w, http.StatusBadRequest, "email code is required")
		return
	}
	codeTg, ok := data["code_telegram"]
	if !ok {
		jsonError(w, http.StatusBadRequest, "telegram code is required")
		return
	}
	user.Login = data["email"].(string)
	user.Password = data["password"].(string)

	codeRE, err := h.rdb.GetData(user.Login)
	if err != nil || string(codeRE) != codeEmail {
		jsonError(w, http.StatusBadRequest, "email code is invalid")
		return
	}
	codeRT, err := h.rdb.GetData("tg_" + user.Login)
	if err != nil || string(codeRT) != codeTg {
		jsonError(w, http.StatusBadRequest, "telegram code is invalid")
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
	h.rdb.DeleteData(user.Login)
	h.rdb.DeleteData("tg_" + user.Login)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(jwt.Claims)
	idString := r.PathValue("id")
	id, _ := strconv.Atoi(idString)
	var err error
	if claims.Role != "user" {
		if idString == "" {
			err = h.db.DeleteUser(claims.ID)
		} else if uint(id) == claims.ID {
			jsonError(w, http.StatusTeapot, "Самый умный!? Удаляй себя с главной страницы!")
		} else {
			err = h.db.DeleteUser(uint(id))
		}
	} else {
		if idString != "" && claims.ID != uint(id) {
			jsonError(w, http.StatusForbidden, "permission denied")
			return
		}
		err = h.db.DeleteUser(claims.ID)
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
