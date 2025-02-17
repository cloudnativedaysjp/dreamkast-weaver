package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestClientIP(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	req.Header.Add("X-Real-IP", "100.100.100.100")
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(ClientIP)

	realIP := ""
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		realIP = ClientIPFromContext(r.Context())
		w.Write([]byte("Hello World"))
	})
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal("Response Code should be 200")
	}

	if realIP != "100.100.100.100" {
		t.Fatal("Test get real IP error.")
	}
}
