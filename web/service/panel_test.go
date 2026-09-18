package service

import "testing"

func TestIsNewerVersion(t *testing.T) {
	cases := []struct {
		latest  string
		current string
		want    bool
	}{
		{"v2.9.4", "2.9.3", true},
		{"v2.10.0", "2.9.9", true},
		{"v2.9.3", "2.9.3", false},
		{"v2.9.2", "2.9.3", false},
		{"v3.0.0", "2.9.3", true},
		// A test build ahead of the stable line is not offered the older
		// stable release, and is offered the final release it leads up to.
		{"v2.9.12", "2.10.0-beta.1", false},
		{"v2.9.13", "2.10.0-beta.1", false},
		{"v2.10.0", "2.10.0-beta.1", true},
		{"v2.10.0-beta.2", "2.10.0-beta.1", true},
		{"v2.10.0-beta.10", "2.10.0-beta.2", true},
		{"v2.10.0-rc.1", "2.10.0-beta.10", true},
		{"v2.10.0-beta.1", "2.10.0", false},
		{"v2.10.0-beta.1", "2.10.0-beta.1", false},
	}

	for _, tc := range cases {
		if got := isNewerVersion(tc.latest, tc.current); got != tc.want {
			t.Fatalf("isNewerVersion(%q, %q) = %v, want %v", tc.latest, tc.current, got, tc.want)
		}
	}
}

func TestCompareVersionStringsRejectsUnexpectedFormats(t *testing.T) {
	if _, ok := compareVersionStrings("latest", "2.9.3"); ok {
		t.Fatal("expected non-semver latest tag to be rejected")
	}
	if _, ok := compareVersionStrings("v2.9", "2.9.3"); ok {
		t.Fatal("expected short version to be rejected")
	}
	if _, ok := compareVersionStrings("v2.10.0-", "2.9.3"); ok {
		t.Fatal("expected an empty pre-release suffix to be rejected")
	}
}

func TestShellQuote(t *testing.T) {
	if got := shellQuote("/usr/bin/curl"); got != "'/usr/bin/curl'" {
		t.Fatalf("unexpected quote result: %s", got)
	}
	if got := shellQuote("/tmp/a'b"); got != "'/tmp/a'\\''b'" {
		t.Fatalf("unexpected quote result with single quote: %s", got)
	}
}
