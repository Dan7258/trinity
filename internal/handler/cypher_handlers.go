package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"trinity/internal/models"
	"trinity/pkg/jwt"
	"trinity/pkg/kuznechik"
	"trinity/pkg/rsa"
	"trinity/pkg/stribog"
)

func (h *Handler) Encode(w http.ResponseWriter, r *http.Request) {
	algorithm := r.PathValue("algorithm")
	if algorithm == "" {
		jsonError(w, http.StatusBadRequest, "algorithm is required")
		return
	}
	data := new(Data)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	var resp []byte
	var encodeData any
	switch algorithm {
	case "kuznechik":
		encodeData = kuznechik.EncryptText(data.Message)
	case "rsa":
		encodeData, err = rsa.EncodeData(data.Message)
	case "stribog":
		encodeData = stribog.HashingText(data.Message)
	default:
		jsonError(w, http.StatusBadRequest, "unknown algorithm")
		return
	}
	claims, ok := r.Context().Value("user").(jwt.Claims)
	if ok {
		userID := claims.ID
		err = h.addDataToDatabase(encodeData, userID)
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
	}
	resp, err = json.Marshal(encodeData)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) Decode(w http.ResponseWriter, r *http.Request) {
	algorithm := r.PathValue("algorithm")
	if algorithm == "" {
		jsonError(w, http.StatusBadRequest, "algorithm is required")
		return
	}
	var err error
	var text *string
	switch algorithm {
	case "kuznechik":
		decodeData := new(kuznechik.EncryptedData)
		err = json.NewDecoder(r.Body).Decode(decodeData)
		text = kuznechik.DecryptText(*decodeData)
	case "rsa":
		decodeData := new(rsa.EncryptedData)
		err = json.NewDecoder(r.Body).Decode(decodeData)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}
		text, err = rsa.DecodeData(decodeData)
	default:
		jsonError(w, http.StatusBadRequest, "unknown algorithm")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	data := new(Data)
	data.Message = *text
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*data)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	algorithm := r.PathValue("algorithm")
	if algorithm == "" {
		jsonError(w, http.StatusBadRequest, "algorithm is required")
		return
	}
	var err error
	data := make([]any, 0)
	claims, ok := r.Context().Value("user").(jwt.Claims)
	if !ok {
		jsonError(w, http.StatusForbidden, "failed to get user claims")
		log.Println("failed to get user claims: ", err)
		return
	}
	id := claims.ID
	switch algorithm {
	case "kuznechik":
		kuznechik := make([]models.Kuznechik, 0)
		kuznechik, err = h.db.GetKuznechikListByUserID(id)
		for _, k := range kuznechik {
			data = append(data, k)
		}
	case "rsa":
		rsaList := make([]models.RSA, 0)
		rsaList, err = h.db.GetRSAListByUserID(id)
		for _, r := range rsaList {
			data = append(data, r)
		}
	case "stribog":
		stribogList := make([]models.Stribog, 0)
		stribogList, err = h.db.GetStribogListByUserID(id)
		for _, s := range stribogList {
			data = append(data, s)
		}
	default:
		jsonError(w, http.StatusBadRequest, "unknown algorithm")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) addDataToDatabase(data interface{}, userID uint) error {
	switch d := data.(type) {
	case *kuznechik.EncryptedData:
		m := new(models.Kuznechik)
		m.EncryptedMessage = d.EncryptedMessage
		m.UserID = userID
		m.Key = d.Key
		return h.db.CreateKuznechik(m)
	case *rsa.EncryptedData:
		m := new(models.RSA)
		m.EncryptedMessage = d.EncryptedMessage
		m.UserID = userID
		m.D = d.D
		m.N = d.N
		return h.db.CreateRSA(m)
	case *stribog.EncryptedData:
		m := new(models.Stribog)
		m.UserID = userID
		m.Hash = d.EncryptedMessage
		return h.db.CreateStribog(m)
	default:
		return errors.New("unknown data type")
	}
}

func (h *Handler) GetHistoryByLogin(w http.ResponseWriter, r *http.Request) {
	algorithm := r.PathValue("algorithm")
	if algorithm == "" {
		jsonError(w, http.StatusBadRequest, "algorithm is required")
		return
	}
	login := r.PathValue("login")
	if login == "" {
		jsonError(w, http.StatusBadRequest, "login is required")
		return
	}
	var err error
	data := make([]any, 0)
	switch algorithm {
	case "kuznechik":
		kuznechik := make([]models.Kuznechik, 0)
		kuznechik, err = h.db.GetKuznechikListByLogin(login)
		for _, k := range kuznechik {
			data = append(data, k)
		}
	case "rsa":
		rsaList := make([]models.RSA, 0)
		rsaList, err = h.db.GetRSAListByLogin(login)
		for _, r := range rsaList {
			data = append(data, r)
		}
	case "stribog":
		stribogList := make([]models.Stribog, 0)
		stribogList, err = h.db.GetStribogListByLogin(login)
		for _, s := range stribogList {
			data = append(data, s)
		}
	default:
		jsonError(w, http.StatusBadRequest, "unknown algorithm")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
