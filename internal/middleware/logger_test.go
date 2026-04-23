package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogger_CallsNext(t *testing.T) {
	r := gin.New()
	r.Use(Logger())
	nextCalled := false
	r.GET("/log", func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/log", nil)
	r.ServeHTTP(w, req)

	if !nextCalled {
		t.Error("expected next handler to be called")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestLogger_PreservesResponseStatus(t *testing.T) {
	r := gin.New()
	r.Use(Logger())
	r.GET("/notfound", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/notfound", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
