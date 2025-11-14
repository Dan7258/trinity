package handler

import (
	"encoding/json"
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
		if err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "stribog":
		encodeData = stribog.HashingText(data.Message)
	default:
		jsonError(w, http.StatusBadRequest, "unknown algorithm")
		return
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
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}
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
	m, err := json.Marshal(data)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	algorithm := r.PathValue("algorithm")
	if algorithm == "" {
		jsonError(w, http.StatusBadRequest, "algorithm is required")
		return
	}
	var err error
	data := make([]interface{}, 0)
	claims := r.Context().Value("user").(jwt.Claims)
	id := claims.ID
	switch algorithm {
	case "kuznechik":
		kuznechik := make([]models.Kuznechik, 0)
		kuznechik, err = h.db.GetKuznechikListByUserID(id)
		data = append(data, kuznechik)
	case "rsa":
		rsaList := make([]models.RSA, 0)
		rsaList, err = h.db.GetRSAListByUserID(id)
		data = append(data, rsaList)
	case "stribog":
		stribogList := make([]models.Stribog, 0)
		stribogList, err = h.db.GetStribogListByUserID(id)
		data = append(data, stribogList)
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
