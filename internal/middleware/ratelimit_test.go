package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, time.Minute)
	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
	if rl.limit != 10 {
		t.Errorf("expected limit 10, got %d", rl.limit)
	}
	if rl.window != time.Minute {
		t.Errorf("expected window 1m, got %v", rl.window)
	}
	if rl.entries == nil {
		t.Error("expected initialized entries map")
	}
}

func makeRateLimiterRouter(limit int, window time.Duration) (*gin.Engine, *RateLimiter) {
	rl := NewRateLimiter(limit, window)
	r := gin.New()
	r.GET("/", rl.Middleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r, rl
}

func TestRateLimiter_AllowsRequestsUnderLimit(t *testing.T) {
	r, _ := makeRateLimiterRouter(5, time.Minute)

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimiter_BlocksRequestsOverLimit(t *testing.T) {
	r, _ := makeRateLimiterRouter(3, time.Minute)

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "5.6.7.8:1234"
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "5.6.7.8:1234"
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	r, _ := makeRateLimiterRouter(1, 50*time.Millisecond)

	sendReq := func() int {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "9.9.9.9:1234"
		r.ServeHTTP(w, req)
		return w.Code
	}

	if code := sendReq(); code != http.StatusOK {
		t.Fatalf("first request: expected 200, got %d", code)
	}
	if code := sendReq(); code != http.StatusTooManyRequests {
		t.Fatalf("second request: expected 429, got %d", code)
	}

	time.Sleep(60 * time.Millisecond)

	if code := sendReq(); code != http.StatusOK {
		t.Errorf("after window reset: expected 200, got %d", code)
	}
}

func TestRateLimiter_DifferentIPsIndependent(t *testing.T) {
	r, _ := makeRateLimiterRouter(1, time.Minute)

	for _, ip := range []string{"10.0.0.1:1", "10.0.0.2:1", "10.0.0.3:1"} {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = ip
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("ip %s: expected 200, got %d", ip, w.Code)
		}
	}
}
