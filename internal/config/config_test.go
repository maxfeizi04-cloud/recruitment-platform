package config_test

import (
	"os"
	"testing"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/config"
)

func TestLoad_ValidConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `
server:
  port: "9090"
  mode: "release"
database:
  host: "db.example.com"
  port: "5432"
  user: "admin"
  password: "secret"
  name: "testdb"
  sslmode: "require"
redis:
  addr: "cache:6379"
  password: "redispass"
  db: 1
jwt:
  secret: "test-secret"
  expire_hours: 72
sms:
  secret_id: "sms-id"
  secret_key: "sms-key"
  sdk_app_id: "14000001"
  template_id: "123456"
  sign_name: "测试"
cos:
  secret_id: "cos-id"
  secret_key: "cos-key"
  bucket_url: "https://test.cos.ap-guangzhou.myqcloud.com"
  region: "ap-guangzhou"
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	cfg, err := config.Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("port = %s, want 9090", cfg.Server.Port)
	}
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("db host = %s, want db.example.com", cfg.Database.Host)
	}
	if cfg.JWT.ExpireHours != 72 {
		t.Errorf("expire_hours = %d, want 72", cfg.JWT.ExpireHours)
	}
	wantDSN := "postgres://admin:secret@db.example.com:5432/testdb?sslmode=require"
	if dsn := cfg.Database.DSN(); dsn != wantDSN {
		t.Errorf("DSN = %s, want %s", dsn, wantDSN)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestJWTConfig_ExpireDuration(t *testing.T) {
	cfg := config.JWTConfig{ExpireHours: 24}
	if d := cfg.ExpireDuration(); d.Hours() != 24 {
		t.Errorf("duration = %v, want 24h", d)
	}
}
