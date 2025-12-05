// Copyright 2025 cloudeng llc. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

//go:build linux

package keyrings_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"cloudeng.io/linux/keyrings"
)

func TestKeyrings(t *testing.T) {
	ctx := context.Background()
	joinSessionKeyring(t)
	kr := keyrings.New()

	name := fmt.Sprintf("test-key-%d", time.Now().UnixNano())
	data := []byte("secret-data")

	// Test Write
	if err := kr.WriteFileCtx(ctx, name, data); err != nil {
		t.Fatalf("WriteFileCtx %v failed: %v", name, err)
	}

	// Test Read
	got, err := kr.ReadFileCtx(ctx, name)
	if err != nil {
		t.Fatalf("ReadFileCtx failed: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Errorf("ReadFileCtx got %q, want %q", got, data)
	}

	// Test Delete
	if err := kr.Delete(ctx, name); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Test Not Exist
	_, err = kr.ReadFileCtx(ctx, name)
	if err == nil || !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("ReadFileCtx expected error for deleted key, got %v", err)
	}
}

func runKeyctl(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), "keyctl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("runKeyctl: %v failed: %s: %v", strings.Join(cmd.Args, " "), out, err)
	}
	return string(bytes.TrimSpace(out))
}

func joinSessionKeyring(t *testing.T) {
	// KEYCTL_JOIN_SESSION_KEYRING = 1
	// Passing an empty string (or 0 arguments) creates a new session keyring
	// and attaches this Go process to it.
	// All exec.Command() calls will now inherit this keyring.
	_, _, errno := syscall.Syscall(syscall.SYS_KEYCTL, 1, 0, 0)
	if errno != 0 {
		t.Fatalf("failed to join session keyring: %v", errno)
	}
}

func TestKeyctlInterop(t *testing.T) {
	ctx := t.Context()
	joinSessionKeyring(t)

	then := time.Now().UnixNano()
	name := fmt.Sprintf("test-key-%d", then)
	data := "secret-data"

	// Test Write
	out := runKeyctl(t, "add", "user", name, data, "@s")
	defer func(key string) {
		runKeyctl(t, "unlink", key, "@s")
	}(out)

	kr := keyrings.New()
	got, err := kr.ReadFileCtx(ctx, name)
	if err != nil {
		t.Fatalf("ReadFileCtx failed: %v", err)
	}
	if string(got) != data {
		t.Errorf("ReadFileCtx got %q, want %q", got, data)
	}

	then += 1
	name = fmt.Sprintf("test-key-%d", then)
	data = "secret-data"

	if err := kr.WriteFileCtx(ctx, name, []byte(data)); err != nil {
		t.Fatalf("WriteFileCtx %v failed: %v", name, err)
	}

	out = runKeyctl(t, "search", "@s", "user", name)
	defer func() {
		kr.Delete(ctx, name)
	}()

	out = runKeyctl(t, "print", out)
	if string(out) != data {
		t.Errorf("runKeyctl: print got %q, want %q", out, data)
	}

}
