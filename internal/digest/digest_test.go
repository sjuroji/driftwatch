package digest_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/digest"
)

type fakeState struct {
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
}

func TestCompute_Deterministic(t *testing.T) {
	s := fakeState{Image: "nginx:1.25", Replicas: 3}

	r1, err := digest.Compute("web", s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := digest.Compute("web", s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r1.Hash != r2.Hash {
		t.Errorf("expected identical hashes, got %q and %q", r1.Hash, r2.Hash)
	}
}

func TestCompute_DifferentValuesProduceDifferentHashes(t *testing.T) {
	a := fakeState{Image: "nginx:1.25", Replicas: 2}
	b := fakeState{Image: "nginx:1.26", Replicas: 2}

	ra, _ := digest.Compute("web", a)
	rb, _ := digest.Compute("web", b)

	if ra.Hash == rb.Hash {
		t.Error("expected different hashes for different states")
	}
}

func TestChanged_DetectsChange(t *testing.T) {
	old, _ := digest.Compute("svc", fakeState{Image: "app:1", Replicas: 1})
	new_, _ := digest.Compute("svc", fakeState{Image: "app:2", Replicas: 1})

	if !digest.Changed(old, new_) {
		t.Error("expected Changed to return true")
	}
}

func TestChanged_SameState(t *testing.T) {
	s := fakeState{Image: "app:1", Replicas: 1}
	r1, _ := digest.Compute("svc", s)
	r2, _ := digest.Compute("svc", s)

	if digest.Changed(r1, r2) {
		t.Error("expected Changed to return false for identical state")
	}
}

func TestComputeAll_ReturnsResultPerEntry(t *testing.T) {
	states := map[string]any{
		"auth":    fakeState{Image: "auth:1", Replicas: 2},
		"gateway": fakeState{Image: "gw:3", Replicas: 1},
	}

	results, err := digest.ComputeAll(states)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestIndex_LookupByService(t *testing.T) {
	r, _ := digest.Compute("auth", fakeState{Image: "auth:1", Replicas: 1})
	idx := digest.Index([]digest.Result{r})

	got, ok := idx["auth"]
	if !ok {
		t.Fatal("expected 'auth' in index")
	}
	if got.Hash != r.Hash {
		t.Errorf("hash mismatch: got %q want %q", got.Hash, r.Hash)
	}
}
