package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "")
	t.Setenv("POSTGRES_USER", "")
	t.Setenv("POSTGRES_PASSWORD", "")
	t.Setenv("POSTGRES_DB", "")
	t.Setenv("POSTGRES_HOST", "")
	t.Setenv("POSTGRES_PORT", "")
	t.Setenv("POSTGRES_SSLMODE", "")
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_TTL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config load without error, got %v", err)
	}

	if cfg.HTTP.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.HTTP.Port)
	}

	if cfg.DB.Host == "" || cfg.DB.DSN == "" {
		t.Fatalf("expected db config and dsn to be populated")
	}

	if cfg.JWT.Secret == "" {
		t.Fatalf("expected jwt secret to be populated")
	}
}
