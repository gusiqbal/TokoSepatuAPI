package account

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"learnapirest/helpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockAccountService implements IAccountService for handler tests.
type mockAccountService struct {
	createAccountFn func(ctx context.Context, account *RegisterUserRequest) error
	loginFn         func(ctx context.Context, username, password string) (*TokenResponse, error)
	logoutFn        func(ctx context.Context, refreshToken string) error
	refreshTokenFn  func(ctx context.Context, refreshToken string) (*TokenResponse, error)
	getProfileFn    func(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	updateProfileFn func(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) error
}

func (m *mockAccountService) CreateAccount(ctx context.Context, req *RegisterUserRequest) error {
	if m.createAccountFn != nil {
		return m.createAccountFn(ctx, req)
	}
	return nil
}
func (m *mockAccountService) Login(ctx context.Context, username, password string) (*TokenResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, username, password)
	}
	return &TokenResponse{}, nil
}
func (m *mockAccountService) Logout(ctx context.Context, refreshToken string) error {
	if m.logoutFn != nil {
		return m.logoutFn(ctx, refreshToken)
	}
	return nil
}
func (m *mockAccountService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if m.refreshTokenFn != nil {
		return m.refreshTokenFn(ctx, refreshToken)
	}
	return &TokenResponse{}, nil
}
func (m *mockAccountService) GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	if m.getProfileFn != nil {
		return m.getProfileFn(ctx, userID)
	}
	return &UserResponse{}, nil
}
func (m *mockAccountService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) error {
	if m.updateProfileFn != nil {
		return m.updateProfileFn(ctx, userID, req)
	}
	return nil
}

func newTestAccountController(svc IAccountService) *AccountController {
	return NewAccountController(svc)
}

func postJSON(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

// --- CreateAccount handler ---

func TestCreateAccountHandler_Success(t *testing.T) {
	ctrl := newTestAccountController(&mockAccountService{})
	r := gin.New()
	r.POST("/account/create", ctrl.CreateAccount)

	w := postJSON(r, "/account/create", gin.H{
		"name": "Alice", "userName": "alice",
		"email": "a@b.com", "password": "pass123", "phoneNumber": "08001",
	})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateAccountHandler_InvalidJSON(t *testing.T) {
	ctrl := newTestAccountController(&mockAccountService{})
	r := gin.New()
	r.POST("/account/create", ctrl.CreateAccount)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/account/create", bytes.NewBufferString("notjson"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateAccountHandler_MissingRequiredFields(t *testing.T) {
	ctrl := newTestAccountController(&mockAccountService{})
	r := gin.New()
	r.POST("/account/create", ctrl.CreateAccount)

	w := postJSON(r, "/account/create", gin.H{"name": "only name"})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}
}

func TestCreateAccountHandler_ServiceError(t *testing.T) {
	svc := &mockAccountService{
		createAccountFn: func(_ context.Context, _ *RegisterUserRequest) error {
			return errors.New("email already exists")
		},
	}
	ctrl := newTestAccountController(svc)
	r := gin.New()
	r.POST("/account/create", ctrl.CreateAccount)

	w := postJSON(r, "/account/create", gin.H{
		"name": "Alice", "userName": "alice",
		"email": "a@b.com", "password": "pass123", "phoneNumber": "08001",
	})
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- Login handler ---

func TestLoginHandler_Success(t *testing.T) {
	svc := &mockAccountService{
		loginFn: func(_ context.Context, _, _ string) (*TokenResponse, error) {
			return &TokenResponse{AccessToken: "access", RefreshToken: "refresh"}, nil
		},
	}
	ctrl := newTestAccountController(svc)
	r := gin.New()
	r.POST("/account/login", ctrl.Login)

	w := postJSON(r, "/account/login", gin.H{"username": "alice", "password": "secret"})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	ctrl := newTestAccountController(&mockAccountService{})
	r := gin.New()
	r.POST("/account/login", ctrl.Login)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/account/login", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLoginHandler_AppError(t *testing.T) {
	svc := &mockAccountService{
		loginFn: func(_ context.Context, _, _ string) (*TokenResponse, error) {
			return nil, helpers.NewError(http.StatusUnauthorized, "Invalid credentials")
		},
	}
	ctrl := newTestAccountController(svc)
	r := gin.New()
	r.POST("/account/login", ctrl.Login)

	w := postJSON(r, "/account/login", gin.H{"username": "alice", "password": "wrong"})
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLoginHandler_GenericError(t *testing.T) {
	svc := &mockAccountService{
		loginFn: func(_ context.Context, _, _ string) (*TokenResponse, error) {
			return nil, errors.New("db error")
		},
	}
	ctrl := newTestAccountController(svc)
	r := gin.New()
	r.POST("/account/login", ctrl.Login)

	w := postJSON(r, "/account/login", gin.H{"username": "alice", "password": "pass"})
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- GetProfile handler ---

func TestGetProfileHandler_Unauthorized(t *testing.T) {
	ctrl := newTestAccountController(&mockAccountService{})
	r := gin.New()
	r.GET("/account/profile", ctrl.GetProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/account/profile", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetProfileHandler_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockAccountService{
		getProfileFn: func(_ context.Context, id uuid.UUID) (*UserResponse, error) {
			return &UserResponse{ID: id.String(), Name: "Alice"}, nil
		},
	}
	ctrl := newTestAccountController(svc)
	r := gin.New()
	r.GET("/account/profile", func(c *gin.Context) {
		c.Set("userId", userID)
	}, ctrl.GetProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/account/profile", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// --- UpdateProfile handler ---

func TestUpdateProfileHandler_Unauthorized(t *testing.T) {
	ctrl := newTestAccountController(&mockAccountService{})
	r := gin.New()
	r.PUT("/account/profile", ctrl.UpdateProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/account/profile", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestUpdateProfileHandler_Success(t *testing.T) {
	userID := uuid.New()
	ctrl := newTestAccountController(&mockAccountService{})
	r := gin.New()
	r.PUT("/account/profile", func(c *gin.Context) {
		c.Set("userId", userID)
	}, ctrl.UpdateProfile)

	b, _ := json.Marshal(gin.H{"name": "New Name"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/account/profile", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Logout handler ---

func TestLogoutHandler_Success(t *testing.T) {
	ctrl := newTestAccountController(&mockAccountService{})
	r := gin.New()
	r.POST("/account/logout", ctrl.Logout)

	w := postJSON(r, "/account/logout", gin.H{"refreshToken": "sometoken"})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// --- RefreshToken handler ---

func TestRefreshTokenHandler_Success(t *testing.T) {
	svc := &mockAccountService{
		refreshTokenFn: func(_ context.Context, _ string) (*TokenResponse, error) {
			return &TokenResponse{AccessToken: "new-access", RefreshToken: "new-refresh"}, nil
		},
	}
	ctrl := newTestAccountController(svc)
	r := gin.New()
	r.POST("/account/refresh", ctrl.RefreshToken)

	w := postJSON(r, "/account/refresh", gin.H{"refreshToken": "old-refresh"})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRefreshTokenHandler_InvalidToken(t *testing.T) {
	svc := &mockAccountService{
		refreshTokenFn: func(_ context.Context, _ string) (*TokenResponse, error) {
			return nil, helpers.NewError(http.StatusUnauthorized, "invalid token")
		},
	}
	ctrl := newTestAccountController(svc)
	r := gin.New()
	r.POST("/account/refresh", ctrl.RefreshToken)

	w := postJSON(r, "/account/refresh", gin.H{"refreshToken": "bad"})
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
