package handler

import (
	"encoding/json"
	"net/http"
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
