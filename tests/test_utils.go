package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func PerformRequest(
	r http.Handler,
	method, path string,
	body any,
	token string,
) *httptest.ResponseRecorder {

	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}

	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ParseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	var res map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	return res
}
