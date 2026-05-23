package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"driftwatch/internal/retry"
)

// zeroDelay eliminates sleep overhead so benchmarks measure pure retry logic.
var zeroDelay = retry.Config{
	MaxAttempts:  3,
	InitialDelay: time.Nanosecond,
	MaxDelay:     time.Nanosecond,
	Multiplier:   1.0,
}

func BenchmarkDo_AlwaysSucceeds(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = retry.Do(ctx, zeroDelay, func() error { return nil })
	}
}

func BenchmarkDo_AlwaysFails(b *testing.B) {
	ctx := context.Background()
	sentinel := errors.New("fail")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = retry.Do(ctx, zeroDelay, func() error { return sentinel })
	}
}

func BenchmarkDo_PermFail(b *testing.B) {
	ctx := context.Background()
	perm := errors.New("perm")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = retry.Do(ctx, zeroDelay, func() error { return retry.PermFail(perm) })
	}
}
