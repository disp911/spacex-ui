package site

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"runtime"
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

func TestCopyFollowsTheStoredCertificate(t *testing.T) {
	folder := t.TempDir()
	c := &Certs{domain: "docs.example.com", folder: folder}
	if cert, key := CopyFiles(folder); cert != "" || key != "" {
		t.Fatalf("a copy before any certificate: %q %q", cert, key)
	}
	store := func(name, chain, key string) {
		dir := filepath.Join(folder, "certificates", "acme", name)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		os.WriteFile(filepath.Join(dir, name+".crt"), []byte(chain), 0o600)
		os.WriteFile(filepath.Join(dir, name+".key"), []byte(key), 0o600)
	}
	store("docs.example.com", "chain one", "key one")
	c.copyStored("certificates/acme/docs.example.com/docs.example.com.crt", "certificates/acme/docs.example.com/docs.example.com.key")
	cert, key := CopyFiles(folder)
	if b, _ := os.ReadFile(cert); string(b) != "chain one" {
		t.Fatalf("copy %q", b)
	}
	if b, _ := os.ReadFile(key); string(b) != "key one" {
		t.Fatalf("key %q", b)
	}
	if info, _ := os.Stat(key); info.Mode().Perm()&0o077 != 0 && runtime.GOOS != "windows" {
		t.Fatalf("key mode %v", info.Mode())
	}

	// A renewal replaces the copy; a missing stored file leaves it alone.
	store("docs.example.com", "chain two", "key two")
	c.copyStored("certificates/acme/docs.example.com/docs.example.com.crt", "certificates/acme/docs.example.com/docs.example.com.key")
	c.copyStored("certificates/acme/gone/gone.crt", "certificates/acme/gone/gone.key")
	if b, _ := os.ReadFile(cert); string(b) != "chain two" {
		t.Fatalf("renewed copy %q", b)
	}
	if b, _ := os.ReadFile(key); string(b) != "key two" {
		t.Fatalf("renewed key %q", b)
	}
	if left, _ := filepath.Glob(filepath.Join(folder, ".copy-*")); len(left) != 0 {
		t.Fatalf("left behind %v", left)
	}
}
