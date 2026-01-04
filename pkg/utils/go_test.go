// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import (
	"context"
	"testing"
	"time"
)

func TestCall(t *testing.T) {
	ctx := context.Background()

	// Test normal execution
	called := false
	Call(ctx, func() {
		called = true
	})
	if !called {
		t.Error("function was not called")
	}

	// Test panic recovery - this will re-panic, so we need to recover in test
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	Call(ctx, func() {
		panic("test panic")
	})
}

func TestCallSafe(t *testing.T) {
	ctx := context.Background()

	// Test normal execution
	called := false
	CallSafe(ctx, func() {
		called = true
	})
	if !called {
		t.Error("function was not called")
	}

	// Test panic recovery - should not re-panic
	CallSafe(ctx, func() {
		panic("test panic")
	})
	// If we reach here, panic was handled
}

func TestGo(t *testing.T) {
	ctx := context.Background()

	// Test normal execution
	done := make(chan bool)
	Go(ctx, func() {
		done <- true
	})
	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Error("goroutine did not complete")
	}
}

func TestPanicIfNotNil(t *testing.T) {
	ctx := context.Background()

	// Test nil - should not panic
	PanicIfNotNil(ctx, nil)

	// Test non-nil - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	PanicIfNotNil(ctx, "test panic")
}

func TestReportPanicIfNotNil(t *testing.T) {
	ctx := context.Background()

	// Test nil
	if ReportPanicIfNotNil(ctx, nil) {
		t.Error("expected false for nil")
	}

	// Test non-nil
	if !ReportPanicIfNotNil(ctx, "test panic") {
		t.Error("expected true for non-nil")
	}
}
