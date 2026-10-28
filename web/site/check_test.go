package site

import (
	"context"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	status := http.StatusOK
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
	srv.EnableHTTP2 = true
	srv.StartTLS()
	defer srv.Close()
	roots := x509.NewCertPool()
	roots.AddCert(srv.Certificate())
	addr := srv.Listener.Addr().String()

	// httptest's certificate is made out to example.com.
	if res := Check(context.Background(), "example.com", addr, roots); !res.OK || res.Error != "" || res.At.IsZero() {
		t.Fatalf("up: %+v", res)
	}
	if res := Check(context.Background(), "other.org", addr, roots); res.OK || !strings.Contains(res.Error, "certificate") {
		t.Fatalf("wrong name: %+v", res)
	}
	if res := Check(context.Background(), "example.com", addr, nil); res.OK {
		t.Fatalf("untrusted: %+v", res)
	}
	status = http.StatusBadGateway
	if res := Check(context.Background(), "example.com", addr, roots); res.OK || res.Error != "the site answered 502 Bad Gateway" {
		t.Fatalf("failing: %+v", res)
	}
}

func TestCheckerReportsChanges(t *testing.T) {
	var c Checker
	steps := []struct {
		ok, changed bool
	}{
		{true, false}, // a first success is no news
		{true, false},
		{false, true},
		{false, false},
		{true, true},
	}
	for i, s := range steps {
		if got := c.Record(CheckResult{OK: s.ok}); got != s.changed {
			t.Fatalf("step %d: changed %v", i, got)
		}
	}
	var first Checker
	if !first.Record(CheckResult{OK: false}) {
		t.Fatal("a first failure is news")
	}
}
