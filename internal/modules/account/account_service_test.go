package account

import (
	"context"
	"errors"
	"learnapirest/helpers"
	"learnapirest/internal/config"
	"testing"

	"github.com/google/uuid"
)

// mockAccountRepo is a test double implementing IAccountRepository.
type mockAccountRepo struct {
	createAccountFn     func(ctx context.Context, req *RegisterUserRequest) error
	getUserByUserNameFn func(ctx context.Context, username, password string) (*User, error)
	getUserByIDFn       func(ctx context.Context, userID uuid.UUID) (*User, error)
	updateUserFn        func(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) error
}

func (m *mockAccountRepo) CreateAccount(ctx context.Context, req *RegisterUserRequest) error {
	if m.createAccountFn != nil {
		return m.createAccountFn(ctx, req)
	}
	return nil
}
func (m *mockAccountRepo) GetUserByUserName(ctx context.Context, username, password string) (*User, error) {
	if m.getUserByUserNameFn != nil {
		return m.getUserByUserNameFn(ctx, username, password)
	}
	return nil, nil
}
func (m *mockAccountRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	if m.getUserByIDFn != nil {
		return m.getUserByIDFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockAccountRepo) UpdateUser(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) error {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, userID, req)
	}
	return nil
}

func newTestAccountService(repo IAccountRepository) *AccountService {
	return NewAccountService(repo, &config.Config{JWTSecret: []byte("test-secret")})
}

// --- CreateAccount ---

func TestCreateAccount_Success(t *testing.T) {
	svc := newTestAccountService(&mockAccountRepo{})
	err := svc.CreateAccount(context.Background(), &RegisterUserRequest{
		Name: "Alice", UserName: "alice", Email: "a@b.com", Password: "pass123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateAccount_HashesPassword(t *testing.T) {
	var capturedReq *RegisterUserRequest
	repo := &mockAccountRepo{
		createAccountFn: func(_ context.Context, req *RegisterUserRequest) error {
			capturedReq = req
			return nil
		},
	}
	svc := newTestAccountService(repo)
	original := "plainpassword"
	_ = svc.CreateAccount(context.Background(), &RegisterUserRequest{
		Name: "Bob", UserName: "bob", Email: "b@c.com", Password: original,
	})
	if capturedReq == nil {
		t.Fatal("repo.CreateAccount was not called")
	}
	if capturedReq.Password == original {
		t.Error("expected password to be hashed before calling repo")
	}
}

func TestCreateAccount_RepoError(t *testing.T) {
	repo := &mockAccountRepo{
		createAccountFn: func(_ context.Context, _ *RegisterUserRequest) error {
			return errors.New("duplicate email")
		},
	}
	svc := newTestAccountService(repo)
	err := svc.CreateAccount(context.Background(), &RegisterUserRequest{
		Name: "C", UserName: "c", Email: "c@d.com", Password: "pass123",
	})
	if err == nil {
		t.Error("expected error from repo to propagate")
	}
}

// --- Login ---

func TestLogin_Success(t *testing.T) {
	userID := uuid.New()
	hashed, _ := helpers.HashPassword("secret")
	repo := &mockAccountRepo{
		getUserByUserNameFn: func(_ context.Context, _, _ string) (*User, error) {
			return &User{ID: userID, PasswordHash: hashed}, nil
		},
	}
	svc := newTestAccountService(repo)
	resp, err := svc.Login(context.Background(), "alice", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Error("expected non-empty tokens in response")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockAccountRepo{
		getUserByUserNameFn: func(_ context.Context, _, _ string) (*User, error) {
			return nil, errors.New("username does not exist")
		},
	}
	svc := newTestAccountService(repo)
	_, err := svc.Login(context.Background(), "nobody", "pass")
	if err == nil {
		t.Error("expected error when user not found")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hashed, _ := helpers.HashPassword("correctpass")
	repo := &mockAccountRepo{
		getUserByUserNameFn: func(_ context.Context, _, _ string) (*User, error) {
			return &User{ID: uuid.New(), PasswordHash: hashed}, nil
		},
	}
	svc := newTestAccountService(repo)
	_, err := svc.Login(context.Background(), "alice", "wrongpass")
	if err == nil {
		t.Error("expected error for wrong password")
	}
}

// --- Logout ---

func TestLogout_InvalidToken_ReturnsNil(t *testing.T) {
	svc := newTestAccountService(&mockAccountRepo{})
	// Logout with invalid/expired token is treated as already-logged-out — returns nil
	err := svc.Logout(context.Background(), "invalid.token")
	if err != nil {
		t.Errorf("unexpected error on logout with invalid token: %v", err)
	}
}

func TestLogout_ValidToken_ReturnsNil(t *testing.T) {
	userID := uuid.New()
	svc := newTestAccountService(&mockAccountRepo{
		getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*User, error) {
			return &User{ID: userID}, nil
		},
	})
	token, _ := helpers.GenerateRefreshToken(userID, []byte("test-secret"))
	err := svc.Logout(context.Background(), token)
	if err != nil {
		t.Errorf("unexpected error on logout: %v", err)
	}
}

// --- GetProfile ---

func TestGetProfile_Success(t *testing.T) {
	userID := uuid.New()
	repo := &mockAccountRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*User, error) {
			return &User{ID: id, Name: "Alice", UserName: "alice", Email: "a@b.com"}, nil
		},
	}
	svc := newTestAccountService(repo)
	profile, err := svc.GetProfile(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.Name != "Alice" {
		t.Errorf("expected name Alice, got %q", profile.Name)
	}
	if profile.ID != userID.String() {
		t.Errorf("expected ID %v, got %v", userID.String(), profile.ID)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	repo := &mockAccountRepo{
		getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*User, error) {
			return nil, errors.New("user does not exist")
		},
	}
	svc := newTestAccountService(repo)
	_, err := svc.GetProfile(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error when user not found")
	}
}

// --- UpdateProfile ---

func TestUpdateProfile_Success(t *testing.T) {
	called := false
	repo := &mockAccountRepo{
		updateUserFn: func(_ context.Context, _ uuid.UUID, _ *UpdateProfileRequest) error {
			called = true
			return nil
		},
	}
	svc := newTestAccountService(repo)
	name := "New Name"
	err := svc.UpdateProfile(context.Background(), uuid.New(), &UpdateProfileRequest{Name: &name})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected repo.UpdateUser to be called")
	}
}

func TestUpdateProfile_RepoError(t *testing.T) {
	repo := &mockAccountRepo{
		updateUserFn: func(_ context.Context, _ uuid.UUID, _ *UpdateProfileRequest) error {
			return errors.New("db error")
		},
	}
	svc := newTestAccountService(repo)
	name := "X"
	err := svc.UpdateProfile(context.Background(), uuid.New(), &UpdateProfileRequest{Name: &name})
	if err == nil {
		t.Error("expected error to propagate")
	}
}

// --- RefreshToken ---

func TestRefreshToken_InvalidToken(t *testing.T) {
	svc := newTestAccountService(&mockAccountRepo{})
	_, err := svc.RefreshToken(context.Background(), "bad.token.here")
	if err == nil {
		t.Error("expected error for invalid refresh token")
	}
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	userID := uuid.New()
	token, _ := helpers.GenerateRefreshToken(userID, []byte("test-secret"))
	repo := &mockAccountRepo{
		getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*User, error) {
			return nil, errors.New("user not found")
		},
	}
	svc := newTestAccountService(repo)
	_, err := svc.RefreshToken(context.Background(), token)
	if err == nil {
		t.Error("expected error when user not found during refresh")
	}
}

func TestRefreshToken_Success(t *testing.T) {
	userID := uuid.New()
	token, _ := helpers.GenerateRefreshToken(userID, []byte("test-secret"))
	repo := &mockAccountRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*User, error) {
			return &User{ID: id}, nil
		},
	}
	svc := newTestAccountService(repo)
	resp, err := svc.RefreshToken(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("expected new access token")
	}
}
