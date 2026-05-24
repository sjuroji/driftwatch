package normalize_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/normalize"
)

func TestServiceName_KebabCase(t *testing.T) {
	if got := normalize.ServiceName("auth-service"); got != "auth-service" {
		t.Fatalf("expected auth-service, got %q", got)
	}
}

func TestServiceName_SnakeCase(t *testing.T) {
	if got := normalize.ServiceName("auth_service"); got != "auth-service" {
		t.Fatalf("expected auth-service, got %q", got)
	}
}

func TestServiceName_CamelCase(t *testing.T) {
	if got := normalize.ServiceName("AuthService"); got != "auth-service" {
		t.Fatalf("expected auth-service, got %q", got)
	}
}

func TestServiceName_MixedSeparators(t *testing.T) {
	if got := normalize.ServiceName("Auth__Service--v2"); got != "auth-service-v2" {
		t.Fatalf("expected auth-service-v2, got %q", got)
	}
}

func TestServiceName_Empty(t *testing.T) {
	if got := normalize.ServiceName(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestFieldKey_CamelCase(t *testing.T) {
	if got := normalize.FieldKey("ImageTag"); got != "image_tag" {
		t.Fatalf("expected image_tag, got %q", got)
	}
}

func TestFieldKey_KebabCase(t *testing.T) {
	if got := normalize.FieldKey("image-tag"); got != "image_tag" {
		t.Fatalf("expected image_tag, got %q", got)
	}
}

func TestFieldKey_AlreadySnakeCase(t *testing.T) {
	if got := normalize.FieldKey("replica_count"); got != "replica_count" {
		t.Fatalf("expected replica_count, got %q", got)
	}
}

func TestFieldKey_Empty(t *testing.T) {
	if got := normalize.FieldKey(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestLabels_NormalisesKeysAndTrimsValues(t *testing.T) {
	input := map[string]string{
		"EnvName":  "  production  ",
		"app-tier": "backend",
	}
	got := normalize.Labels(input)

	if v, ok := got["env_name"]; !ok || v != "production" {
		t.Fatalf("expected env_name=production, got %v", got)
	}
	if v, ok := got["app_tier"]; !ok || v != "backend" {
		t.Fatalf("expected app_tier=backend, got %v", got)
	}
}

func TestLabels_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"MyKey": "value"}
	_ = normalize.Labels(input)
	if _, ok := input["my_key"]; ok {
		t.Fatal("original map should not be mutated")
	}
	if _, ok := input["MyKey"]; !ok {
		t.Fatal("original key should still be present")
	}
}

func TestLabels_EmptyMap(t *testing.T) {
	got := normalize.Labels(map[string]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}
