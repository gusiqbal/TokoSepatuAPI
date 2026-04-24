package helpers

import (
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(http.StatusNotFound, "not found")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Code != http.StatusNotFound {
		t.Errorf("expected code %d, got %d", http.StatusNotFound, err.Code)
	}
	if err.Message != "not found" {
		t.Errorf("expected message %q, got %q", "not found", err.Message)
	}
}

func TestNewError_ZeroValues(t *testing.T) {
	err := NewError(0, "")
	if err.Code != 0 {
		t.Errorf("expected code 0, got %d", err.Code)
	}
	if err.Message != "" {
		t.Errorf("expected empty message, got %q", err.Message)
	}
}

func TestAppError_Error(t *testing.T) {
	err := NewError(http.StatusBadRequest, "bad request")
	if err.Error() != "bad request" {
		t.Errorf("expected %q, got %q", "bad request", err.Error())
	}
}

func TestAppError_ImplementsError(t *testing.T) {
	var _ error = NewError(500, "test")
}
