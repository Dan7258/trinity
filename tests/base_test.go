package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"trinity/internal/handler"
)

func baseTest(t *testing.T) {
	if 1 == 2 {
		t.Error("Base test failed")
	}
}

func TestEncode_Kuznechik(t *testing.T) {
	h := &handler.Handler{}

	data := handler.Data{
		Message: "test message",
	}
	body, _ := json.Marshal(data)

	req := httptest.NewRequest(http.MethodPost, "/encode/kuznechik", bytes.NewReader(body))
	req.SetPathValue("algorithm", "kuznechik")

	w := httptest.NewRecorder()

	h.Encode(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	if w.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}
