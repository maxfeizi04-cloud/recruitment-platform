package auth_test

import (
	"testing"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/config"
	pkgauth "github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/auth"

	"github.com/google/uuid"
)

func TestJWTManager_GenerateAndValidate(t *testing.T) {
	cfg := config.JWTConfig{Secret: "test-secret-key", ExpireHours: 1}
	manager := pkgauth.NewJWTManager(cfg)

	userID := uuid.New()
	token, err := manager.Generate(userID, "candidate")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	claims, err := manager.Validate(token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if claims.UserID != userID.String() {
		t.Errorf("UserID = %s, want %s", claims.UserID, userID.String())
	}
	if claims.Role != "candidate" {
		t.Errorf("Role = %s, want candidate", claims.Role)
	}
}

func TestJWTManager_Validate_TamperedToken(t *testing.T) {
	cfg := config.JWTConfig{Secret: "test-secret-key", ExpireHours: 1}
	manager := pkgauth.NewJWTManager(cfg)

	userID := uuid.New()
	token, _ := manager.Generate(userID, "candidate")

	tampered := token + "x"
	_, err := manager.Validate(tampered)
	if err == nil {
		t.Error("expected error for tampered token, got nil")
	}
}

func TestJWTManager_Validate_WrongSecret(t *testing.T) {
	cfg1 := config.JWTConfig{Secret: "secret-a", ExpireHours: 1}
	cfg2 := config.JWTConfig{Secret: "secret-b", ExpireHours: 1}
	m1 := pkgauth.NewJWTManager(cfg1)
	m2 := pkgauth.NewJWTManager(cfg2)

	userID := uuid.New()
	token, _ := m1.Generate(userID, "candidate")

	_, err := m2.Validate(token)
	if err == nil {
		t.Error("expected error with wrong secret, got nil")
	}
}
