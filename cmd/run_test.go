// Copyright 2020 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/open-policy-agent/opa/test/e2e"
)

func TestRunServerBase(t *testing.T) {
	params := newRunParams()
	params.rt.Addrs = &[]string{":0"}
	params.rt.DiagnosticAddrs = &[]string{}
	params.serverMode = true
	ctx, cancel := context.WithCancel(context.Background())

	rt := initRuntime(ctx, params, nil)

	testRuntime := e2e.WrapRuntime(ctx, cancel, rt)

	done := make(chan bool)
	go func() {
		startRuntime(ctx, rt, true)
		done <- true
	}()

	err := testRuntime.WaitForServer()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	validateBasicServe(t, testRuntime)

	cancel()
	<-done
}

func TestRunServerWithDiagnosticAddr(t *testing.T) {
	params := newRunParams()
	params.rt.Addrs = &[]string{":0"}
	params.rt.DiagnosticAddrs = &[]string{":0"}
	ctx, cancel := context.WithCancel(context.Background())

	rt := initRuntime(ctx, params, nil)

	testRuntime := e2e.WrapRuntime(ctx, cancel, rt)

	done := make(chan bool)
	go func() {
		startRuntime(ctx, rt, true)
		done <- true
	}()

	err := testRuntime.WaitForServer()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	validateBasicServe(t, testRuntime)

	diagURL, err := testRuntime.AddrToURL(rt.DiagnosticAddrs()[0])
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if err := testRuntime.HealthCheck(diagURL); err != nil {
		t.Error(err)
	}

	cancel()
	<-done
}

func validateBasicServe(t *testing.T, runtime *e2e.TestRuntime) {
	t.Helper()

	err := runtime.UploadData(bytes.NewBufferString(`{"x": 1}`))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	resp := struct {
		Result int `json:"result"`
	}{}
	err = runtime.GetDataWithInputTyped("x", nil, &resp)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if resp.Result != 1 {
		t.Fatalf("Expected x to be 1, got %v", resp)
	}
}
