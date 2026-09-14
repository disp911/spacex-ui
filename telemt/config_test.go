package telemt

import (
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
)

func testSettings() Settings {
	return Settings{Port: 8443, TLSDomain: "www.example.com", APIPort: 9091, APIToken: "tok"}
}

func TestBuildConfigRendersUsersAndLimits(t *testing.T) {
	expiry := time.Date(2026, 12, 31, 23, 59, 59, 0, time.FixedZone("MSK", 3*3600))
	users := []User{
		{Name: "alice", Secret: "00112233445566778899aabbccddeeff", Enabled: true, ExpiresAt: expiry, QuotaBytes: 1 << 30, MaxUniqueIPs: 3},
		{Name: "bob.phone", Secret: "0123456789abcdef0123456789abcdef", Enabled: false},
	}
	s := testSettings()
	s.PublicHost = "proxy.example.com"
	data, err := BuildConfig(s, users)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := toml.Unmarshal(data, &got); err != nil {
		t.Fatalf("generated config is not valid TOML: %v\n%s", err, data)
	}
	get := func(path string) any {
		var cur any = got
		for _, key := range strings.Split(path, "/") {
			m, ok := cur.(map[string]any)
			if !ok {
				return nil
			}
			cur = m[key]
		}
		return cur
	}
	want := map[string]any{
		"general/modes/tls":                        true,
		"general/modes/classic":                    false,
		"general/modes/secure":                     false,
		"general/links/public_host":                "proxy.example.com",
		"general/disable_colors":                   true,
		"server/port":                              int64(8443),
		"server/api/listen":                        "127.0.0.1:9091",
		"server/api/auth_header":                   "Bearer tok",
		"censorship/tls_domain":                    "www.example.com",
		"access/users/alice":                       "00112233445566778899aabbccddeeff",
		"access/users/bob.phone":                   "0123456789abcdef0123456789abcdef",
		"access/user_enabled/bob.phone":            false,
		"access/user_expirations/alice":            "2026-12-31T20:59:59Z",
		"access/user_data_quota/alice":             int64(1 << 30),
		"access/user_max_unique_ips/alice":         int64(3),
		"access/user_enabled/alice":                nil,
		"access/user_data_quota/bob.phone":         nil,
		"access/user_max_unique_ips/bob.phone":     nil,
		"access/user_expirations/bob.phone":        nil,
		"server/api/whitelist":                     []any{"127.0.0.1/32"},
		"censorship/mask":                          true,
		"general/quota_state_path":                 "telemt.limit.json",
		"censorship/tls_front_dir":                 "tlsfront",
		"server/api/enabled":                       true,
		"censorship/tls_emulation":                 true,
		"general/log_level":                        "normal",
		"access/users/does-not-exist-sanity-check": nil,
	}
	for path, value := range want {
		gotValue := get(path)
		if !equalValue(gotValue, value) {
			t.Errorf("%s = %#v, want %#v", path, gotValue, value)
		}
	}
}

func equalValue(a, b any) bool {
	as, aok := a.([]any)
	bs, bok := b.([]any)
	if aok || bok {
		if len(as) != len(bs) {
			return false
		}
		for i := range as {
			if as[i] != bs[i] {
				return false
			}
		}
		return true
	}
	return a == b
}

func TestBuildConfigOmitsEmptyPublicHost(t *testing.T) {
	data, err := BuildConfig(testSettings(), []User{{Name: "a", Secret: strings.Repeat("a", 32), Enabled: true}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "public_host") {
		t.Fatalf("empty public host must not be written:\n%s", data)
	}
}

func TestBuildConfigRejectsInvalidInput(t *testing.T) {
	ok := User{Name: "a", Secret: strings.Repeat("a", 32), Enabled: true}
	cases := map[string]struct {
		settings Settings
		users    []User
	}{
		"no users":           {testSettings(), nil},
		"bad port":           {Settings{Port: 0, TLSDomain: "example.com", APIPort: 9091}, []User{ok}},
		"api port clash":     {Settings{Port: 9091, TLSDomain: "example.com", APIPort: 9091}, []User{ok}},
		"empty domain":       {Settings{Port: 8443, APIPort: 9091}, []User{ok}},
		"domain with path":   {Settings{Port: 8443, TLSDomain: "example.com/x", APIPort: 9091}, []User{ok}},
		"domain with quote":  {Settings{Port: 8443, TLSDomain: `exa"mple.com`, APIPort: 9091}, []User{ok}},
		"host with space":    {Settings{Port: 8443, TLSDomain: "example.com", PublicHost: "a b", APIPort: 9091}, []User{ok}},
		"bad user name":      {testSettings(), []User{{Name: "bad name", Secret: ok.Secret}}},
		"quote in user name": {testSettings(), []User{{Name: `a"b`, Secret: ok.Secret}}},
		"short secret":       {testSettings(), []User{{Name: "a", Secret: "abc"}}},
		"uppercase secret":   {testSettings(), []User{{Name: "a", Secret: strings.Repeat("A", 32)}}},
		"duplicate user":     {testSettings(), []User{ok, ok}},
	}
	for name, tc := range cases {
		if _, err := BuildConfig(tc.settings, tc.users); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestLinkMatchesTelemtFormat(t *testing.T) {
	got := Link("1.2.3.4", 8443, "00112233445566778899aabbccddeeff", "www.telia.se")
	want := "tg://proxy?server=1.2.3.4&port=8443&secret=ee00112233445566778899aabbccddeeff7777772e74656c69612e7365"
	if got != want {
		t.Fatalf("Link() = %s, want %s", got, want)
	}
}

func TestRestartKeyIgnoresUsers(t *testing.T) {
	a := testSettings()
	b := testSettings()
	if a.restartKey() != b.restartKey() {
		t.Fatal("identical settings must share a restart key")
	}
	for name, mutate := range map[string]func(*Settings){
		"port":        func(s *Settings) { s.Port++ },
		"domain":      func(s *Settings) { s.TLSDomain = "other.example.com" },
		"public host": func(s *Settings) { s.PublicHost = "h" },
		"api port":    func(s *Settings) { s.APIPort++ },
		"api token":   func(s *Settings) { s.APIToken = "other" },
	} {
		c := testSettings()
		mutate(&c)
		if c.restartKey() == a.restartKey() {
			t.Errorf("changing %s must change the restart key", name)
		}
	}
}
