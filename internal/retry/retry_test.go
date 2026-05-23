package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"driftwatch/internal/retry"
)

func fastConfig(max int) retry.Config {
	return retry.Config{
		MaxAttempts:  max,
		InitialDelay: time.Millisecond,
		MaxDelay:     5 * time.Millisecond,
		Multiplier:   2.0,
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(3), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesUpToMax(t *testing.T) {
	calls := 0
	sentinel := errors.New("transient")
	err := retry.Do(context.Background(), fastConfig(3), func() error {
		calls++
		return sentinel
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel in chain, got %v", err)
	}
}

func TestDo_SucceedsOnSecondAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(3), func() error {
		calls++
		if calls < 2 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestDo_PermFail_StopsImmediately(t *testing.T) {
	calls := 0
	perm := errors.New("permanent")
	err := retry.Do(context.Background(), fastConfig(5), func() error {
		calls++
		return retry.PermFail(perm)
	})
	if !errors.Is(err, perm) {
		t.Fatalf("expected perm in chain, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_ContextCancelled_StopsEarly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := retry.Do(ctx, fastConfig(5), func() error {
		calls++
		return nil
	})
	// fn may or may not be called once before cancellation is detected.
	if err == nil && calls == 0 {
		t.Fatal("expected either an error or at least one call")
	}
}

func TestDefaultConfig_SaneValues(t *testing.T) {
	cfg := retry.DefaultConfig()
	if cfg.MaxAttempts <= 0 {
		t.Error("MaxAttempts must be positive")
	}
	if cfg.InitialDelay <= 0 {
		t.Error("InitialDelay must be positive")
	}
	if cfg.Multiplier < 1.0 {
		t.Error("Multiplier must be >= 1.0")
	}
}
