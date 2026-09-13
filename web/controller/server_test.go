package controller

import (
	"regexp"
	"testing"
)

func TestXrayLogFilename(t *testing.T) {
	cases := []struct {
		day, email, want string
	}{
		{"", "", "xray.log"},
		{"", "iPhone-Stepan", "xray.log"},
		{"2026-09-13", "", "xray-2026-09-13.log"},
		{"2026-09-13", "iPhone-Stepan", "xray-2026-09-13-iPhone-Stepan.log"},
		{"2026-09-13", "user@mail.com", "xray-2026-09-13-user_mail.com.log"},
	}
	for _, c := range cases {
		if got := xrayLogFilename(c.day, c.email); got != c.want {
			t.Errorf("xrayLogFilename(%q, %q) = %q, want %q", c.day, c.email, got, c.want)
		}
	}
}

func TestXrayLogFilename_CannotBreakTheHeader(t *testing.T) {
	// The client name comes from the request and is placed inside a quoted
	// Content-Disposition value; a quote or line break would let it escape.
	got := xrayLogFilename("2026-09-13", "a\"b\r\nSet-Cookie: x=1; ../../etc/passwd")

	if !regexp.MustCompile(`^[a-zA-Z0-9_\-.]+$`).MatchString(got) {
		t.Fatalf("filename %q contains characters outside the safe set", got)
	}
}
