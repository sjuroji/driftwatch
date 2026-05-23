package redact_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/redact"
)

func TestMap_RedactsSensitiveKeys(t *testing.T) {
	input := map[string]string{
		"username":    "admin",
		"password":    "s3cr3t",
		"db_host":     "localhost",
		"api_key":     "abc123",
		"description": "main service",
	}

	got := redact.Map(input, redact.Options{})

	if got["username"] != "admin" {
		t.Errorf("username: want %q, got %q", "admin", got["username"])
	}
	if got["password"] != "[REDACTED]" {
		t.Errorf("password: want [REDACTED], got %q", got["password"])
	}
	if got["api_key"] != "[REDACTED]" {
		t.Errorf("api_key: want [REDACTED], got %q", got["api_key"])
	}
	if got["db_host"] != "localhost" {
		t.Errorf("db_host: want %q, got %q", "localhost", got["db_host"])
	}
}

func TestMap_CaseInsensitiveMatch(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD": "hunter2",
		"Auth_Token":  "tok",
		"Region":      "us-east-1",
	}

	got := redact.Map(input, redact.Options{})

	if got["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD: want [REDACTED], got %q", got["DB_PASSWORD"])
	}
	if got["Auth_Token"] != "[REDACTED]" {
		t.Errorf("Auth_Token: want [REDACTED], got %q", got["Auth_Token"])
	}
	if got["Region"] != "us-east-1" {
		t.Errorf("Region: want %q, got %q", "us-east-1", got["Region"])
	}
}

func TestMap_CustomSensitiveKeys(t *testing.T) {
	input := map[string]string{
		"password": "should-not-redact",
		"pin_code": "1234",
		"name":     "svc",
	}

	opts := redact.Options{SensitiveKeys: []string{"pin"}}
	got := redact.Map(input, opts)

	if got["password"] != "should-not-redact" {
		t.Errorf("password should not be redacted with custom keys")
	}
	if got["pin_code"] != "[REDACTED]" {
		t.Errorf("pin_code: want [REDACTED], got %q", got["pin_code"])
	}
}

func TestMap_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"secret": "mysecret"}
	_ = redact.Map(input, redact.Options{})
	if input["secret"] != "mysecret" {
		t.Error("Map must not mutate the original map")
	}
}

func TestValue_SensitiveKey(t *testing.T) {
	got := redact.Value("auth_token", "Bearer xyz", redact.Options{})
	if got != "[REDACTED]" {
		t.Errorf("want [REDACTED], got %q", got)
	}
}

func TestValue_SafeKey(t *testing.T) {
	got := redact.Value("image", "nginx:latest", redact.Options{})
	if got != "nginx:latest" {
		t.Errorf("want %q, got %q", "nginx:latest", got)
	}
}
