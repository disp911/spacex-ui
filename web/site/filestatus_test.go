package site

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

func TestFileStatusReadsTheAdministratorsCertificate(t *testing.T) {
	cert, err := selfSigned("docs.example.com")
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "cert.pem")
	if err := os.WriteFile(file, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]}), 0o600); err != nil {
		t.Fatal(err)
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	s, err := FileStatus(file)
	if err != nil {
		t.Fatal(err)
	}
	if s.Domain != "docs.example.com" || s.Issuer == "" || !s.NotAfter.Equal(leaf.NotAfter) {
		t.Fatalf("status %+v", s)
	}
	if _, err := FileStatus(filepath.Join(t.TempDir(), "none.pem")); err == nil {
		t.Fatal("a missing file was read")
	}
}
