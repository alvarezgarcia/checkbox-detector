package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetectHandler_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/detect", nil)
	w := httptest.NewRecorder()

	detectHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestDetectHandler_MissingImage(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/detect", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()

	detectHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDetectHandler_NotAnImage(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "file.txt")
	part.Write([]byte("this is not an image"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/detect", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	detectHandler(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestDetectHandler_ErrorsAreJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/detect", nil)
	w := httptest.NewRecorder()

	detectHandler(w, req)

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Errorf("response is not valid JSON: %v", err)
	}
	if _, ok := body["error"]; !ok {
		t.Error("response JSON missing 'error' key")
	}
}
