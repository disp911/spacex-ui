package site

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"
)

func writeFile(t *testing.T, name, data string) {
	t.Helper()
	if err := os.WriteFile(name, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestFetchEAB(t *testing.T) {
	var gotEmail string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		gotEmail = r.PostForm.Get("email")
		if gotEmail == "bad@example.com" {
			w.Write([]byte(`{"success":false,"error":{"code":2900,"type":"invalid_email"}}`))
			return
		}
		w.Write([]byte(`{"success":true,"eab_kid":"kid-1","eab_hmac_key":"mac-1"}`))
	}))
	defer srv.Close()

	eab, err := fetchEAB(context.Background(), srv.URL, "admin@example.com")
	if err != nil || eab.KeyID != "kid-1" || eab.MACKey != "mac-1" || gotEmail != "admin@example.com" {
		t.Fatalf("got %+v, %v, email %q", eab, err, gotEmail)
	}
	if _, err := fetchEAB(context.Background(), srv.URL, "bad@example.com"); err == nil || err.Error() != "ZeroSSL account binding: invalid_email" {
		t.Fatalf("got %v", err)
	}
}

func TestCertsIssuers(t *testing.T) {
	without := NewCerts(t.TempDir(), "example.com", "")
	defer without.Stop()
	if n := len(without.cfg.Issuers); n != 1 {
		t.Fatalf("without an email: %d issuers", n)
	}
	with := NewCerts(t.TempDir(), "example.com", "admin@example.com")
	defer with.Stop()
	if n := len(with.cfg.Issuers); n != 2 {
		t.Fatalf("with an email: %d issuers", n)
	}
	if key := with.cfg.Issuers[1].IssuerKey(); key != "acme.zerossl.com-v2-DV90" {
		t.Fatalf("ZeroSSL issuer key %q", key)
	}
}

func TestCertsTLSConfig(t *testing.T) {
	c := NewCerts(t.TempDir(), "example.com", "")
	defer c.Stop()
	cfg := c.TLSConfig()
	if !slices.Equal(cfg.NextProtos, []string{"h2", "http/1.1", "acme-tls/1"}) || cfg.MinVersion != tls.VersionTLS12 || cfg.GetCertificate == nil {
		t.Fatalf("got %v, min %x", cfg.NextProtos, cfg.MinVersion)
	}
	if s := c.Status(); s.Domain != "example.com" || !s.NotAfter.IsZero() || s.Error != "" {
		t.Fatalf("status %+v", s)
	}
}
