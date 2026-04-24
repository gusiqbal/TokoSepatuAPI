package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

var jwtTestSecret = []byte("test-secret")

func makeTestToken(userID uuid.UUID, secret []byte) string {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     9999999999,
		"iss":     "toko-sepatu-api",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := token.SignedString(secret)
	return str
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/", JWTAuth(jwtTestSecret), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_NoBearerPrefix(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/", JWTAuth(jwtTestSecret), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token some-token")
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/", JWTAuth(jwtTestSecret), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_WrongSecret(t *testing.T) {
	userID := uuid.New()
	token := makeTestToken(userID, []byte("other-secret"))

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/", JWTAuth(jwtTestSecret), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_ValidToken_SetsContextAndCallsNext(t *testing.T) {
	userID := uuid.New()
	token := makeTestToken(userID, jwtTestSecret)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	nextCalled := false
	r.GET("/", JWTAuth(jwtTestSecret), func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !nextCalled {
		t.Error("expected next handler to be called")
	}
}
