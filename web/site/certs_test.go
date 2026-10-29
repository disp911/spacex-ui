package site

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"slices"
	"testing"

	"github.com/caddyserver/certmagic"
)

func writeFile(t *testing.T, name, data string) {
	t.Helper()
	if err := os.WriteFile(name, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCertsComeFromLetsEncryptOnly(t *testing.T) {
	c := NewCerts(t.TempDir(), "example.com")
	defer c.Stop()
	if n := len(c.cfg.Issuers); n != 1 || c.cfg.Issuers[0] != c.le || c.le.CA != certmagic.LetsEncryptProductionCA {
		t.Fatalf("%d issuers, CA %q", n, c.le.CA)
	}
}

func TestCertsTLSConfig(t *testing.T) {
	c := NewCerts(t.TempDir(), "example.com")
	defer c.Stop()
	cfg := c.TLSConfig()
	if !slices.Equal(cfg.NextProtos, []string{"h2", "http/1.1", "acme-tls/1"}) || cfg.MinVersion != tls.VersionTLS12 || cfg.GetCertificate == nil {
		t.Fatalf("got %v, min %x", cfg.NextProtos, cfg.MinVersion)
	}
	if s := c.Status(); s.Domain != "example.com" || !s.NotAfter.IsZero() || s.Error != "" {
		t.Fatalf("status %+v", s)
	}
}

func TestCertsStandInUntilThereIsOne(t *testing.T) {
	for _, name := range []string{"docs.example.com", "203.0.113.7", "2001:db8::7"} {
		c := NewCerts(t.TempDir(), name)
		if ip := net.ParseIP(name); (ip != nil) != (c.le.Profile == "shortlived") {
			t.Errorf("%s: profile %q", name, c.le.Profile)
		}
		a, b := net.Pipe()
		cert, err := c.TLSConfig().GetCertificate(&tls.ClientHelloInfo{Conn: a})
		a.Close()
		b.Close()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		leaf, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			t.Fatal(err)
		}
		if err := leaf.VerifyHostname(name); err != nil {
			t.Errorf("%s: the stand-in is for another name: %v", name, err)
		}
		c.Stop()
	}
}
