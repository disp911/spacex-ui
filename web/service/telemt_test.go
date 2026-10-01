package service

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/disp911/spacex-ui/v2/database/model"
	"github.com/disp911/spacex-ui/v2/xray"
)

func TestValidateMTProtoClients(t *testing.T) {
	good := model.Client{Email: "a", ID: strings.Repeat("0f", 16)}
	if err := validateMTProtoClients(model.MTProto, []model.Client{good}); err != nil {
		t.Fatalf("valid secret rejected: %v", err)
	}
	for _, id := range []string{"", "0f0f", strings.Repeat("0F", 16), strings.Repeat("zz", 16), "8e3a3c8e-0f57-4a57-9a6e-4b3c7f0f5e11"} {
		bad := model.Client{Email: "b", ID: id}
		if err := validateMTProtoClients(model.MTProto, []model.Client{good, bad}); err == nil {
			t.Errorf("secret %q must be rejected", id)
		}
	}
	// Other protocols keep their own ID rules.
	uuidClient := model.Client{Email: "c", ID: "8e3a3c8e-0f57-4a57-9a6e-4b3c7f0f5e11"}
	if err := validateMTProtoClients(model.VLESS, []model.Client{uuidClient}); err != nil {
		t.Fatalf("vless clients must not be checked as mtproto: %v", err)
	}
}

func TestValidateMTProtoInboundIgnoresOtherProtocols(t *testing.T) {
	inbound := &model.Inbound{Protocol: model.VLESS, Settings: `not json`}
	if err := validateMTProtoInbound(inbound, nil); err != nil {
		t.Fatalf("non-mtproto inbounds must pass unchanged: %v", err)
	}
}

func TestIsXrayProtocol(t *testing.T) {
	if model.IsXrayProtocol(model.MTProto) {
		t.Fatal("mtproto must not be treated as an Xray protocol")
	}
	for _, p := range []model.Protocol{model.VMESS, model.VLESS, model.Trojan, model.Shadowsocks, model.Hysteria, model.WireGuard} {
		if !model.IsXrayProtocol(p) {
			t.Errorf("%s must be an Xray protocol", p)
		}
	}
}

func mtprotoInbound(t *testing.T, enable bool, settings map[string]any) *model.Inbound {
	t.Helper()
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	return &model.Inbound{Id: 7, Protocol: model.MTProto, Enable: enable, Settings: string(data)}
}

func TestPrepareMTProtoRelayKeepsTheStoredRelay(t *testing.T) {
	// A new route through Xray gets a relay with loopback port and credentials.
	inbound := mtprotoInbound(t, true, map[string]any{"tlsDomain": "example.com", "outboundTag": " warp ", "clients": []any{}})
	if err := prepareMTProtoRelay(inbound, "", nil); err != nil {
		t.Fatal(err)
	}
	var first mtprotoSettings
	if err := json.Unmarshal([]byte(inbound.Settings), &first); err != nil {
		t.Fatal(err)
	}
	if first.OutboundTag != "warp" || !first.XrayRelay.valid() {
		t.Fatalf("relay not generated: %+v", first)
	}
	key := mtprotoRelayKey(inbound)
	if key == "" {
		t.Fatal("an enabled inbound with a relay must have a relay key")
	}

	// The form never sends the relay back; saving again keeps the stored one,
	// so neither Xray nor telemt restarts.
	edited := mtprotoInbound(t, true, map[string]any{"tlsDomain": "example.org", "outboundTag": "warp", "clients": []any{}})
	if err := prepareMTProtoRelay(edited, inbound.Settings, nil); err != nil {
		t.Fatal(err)
	}
	if mtprotoRelayKey(edited) != key {
		t.Fatalf("an unrelated edit changed the relay: %s vs %s", mtprotoRelayKey(edited), key)
	}

	// Switching back to a direct connection drops the relay.
	direct := mtprotoInbound(t, true, map[string]any{"tlsDomain": "example.org", "outboundTag": "", "clients": []any{}})
	if err := prepareMTProtoRelay(direct, edited.Settings, nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(direct.Settings, "xrayRelay") || strings.Contains(direct.Settings, "outboundTag") || mtprotoRelayKey(direct) != "" {
		t.Fatalf("a direct inbound must carry no relay: %s", direct.Settings)
	}

	// A disabled inbound needs nothing from Xray.
	edited.Enable = false
	if mtprotoRelayKey(edited) != "" {
		t.Fatal("a disabled inbound must not ask Xray for a relay")
	}
}

func TestFreeRelayPortAvoidsTakenPorts(t *testing.T) {
	first, err := freeRelayPort(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Whatever the kernel offers next, a port already taken is never chosen.
	for range 5 {
		port, err := freeRelayPort(map[int]bool{first: true})
		if err != nil {
			t.Fatal(err)
		}
		if port == first {
			t.Fatalf("freeRelayPort returned the taken port %d", port)
		}
	}
}

func TestAddMTProtoRelaysRoutesAheadOfOtherRules(t *testing.T) {
	relay := map[string]any{"port": 20000, "user": "u1", "pass": "p1"}
	routed := mtprotoInbound(t, true, map[string]any{"outboundTag": "warp", "xrayRelay": relay})
	gone := mtprotoInbound(t, true, map[string]any{"outboundTag": "removed", "xrayRelay": relay})
	gone.Id = 8
	disabled := mtprotoInbound(t, false, map[string]any{"outboundTag": "warp", "xrayRelay": relay})
	disabled.Id = 9
	cfg := &xray.Config{
		OutboundConfigs: []byte(`[{"tag":"direct","protocol":"freedom"},{"tag":"warp","protocol":"wireguard"}]`),
		RouterConfig:    []byte(`{"domainStrategy":"AsIs","rules":[{"type":"field","outboundTag":"blocked","ip":["geoip:private"]}]}`),
	}
	if err := addMTProtoRelays(cfg, []*model.Inbound{routed, gone, disabled}); err != nil {
		t.Fatal(err)
	}
	if len(cfg.InboundConfigs) != 1 {
		t.Fatalf("only the enabled inbound with an existing outbound gets a relay: %+v", cfg.InboundConfigs)
	}
	in := cfg.InboundConfigs[0]
	if in.Tag != "mtproto-relay-7" || in.Port != 20000 || in.Protocol != "socks" || string(in.Listen) != `"127.0.0.1"` ||
		!strings.Contains(string(in.Settings), `"user":"u1"`) || !strings.Contains(string(in.Settings), `"auth":"password"`) {
		t.Fatalf("relay inbound = %+v %s", in, in.Settings)
	}
	var routing struct {
		DomainStrategy string           `json:"domainStrategy"`
		Rules          []map[string]any `json:"rules"`
	}
	if err := json.Unmarshal(cfg.RouterConfig, &routing); err != nil {
		t.Fatal(err)
	}
	if routing.DomainStrategy != "AsIs" || len(routing.Rules) != 2 {
		t.Fatalf("routing = %s", cfg.RouterConfig)
	}
	if first := routing.Rules[0]; first["outboundTag"] != "warp" || first["inboundTag"].([]any)[0] != "mtproto-relay-7" {
		t.Fatalf("the relay rule must come first: %s", cfg.RouterConfig)
	}
}

func TestMTProtoUserCarriesLimits(t *testing.T) {
	expiry := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	c := model.Client{ID: strings.Repeat("0f", 16), LimitIP: 2, MaxConns: 5, RateUp: 3, RateDown: 25, ExpiryTime: expiry.UnixMilli()}
	u := mtprotoUser("alice", c, true)
	if !u.Disabled || u.MaxUniqueIPs != 2 || u.MaxTCPConns != 5 || u.RateUpBps != 3_000_000 || u.RateDownBps != 25_000_000 || !u.Expires.Equal(expiry) {
		t.Fatalf("user = %+v", u)
	}
	// An expiry counted from the first connection is negative and stays with
	// the panel; negative limits mean nothing.
	c.ExpiryTime, c.RateUp, c.MaxConns = -86400000, -1, -1
	u = mtprotoUser("alice", c, false)
	if !u.Expires.IsZero() || u.RateUpBps != 0 || u.MaxTCPConns != 0 {
		t.Fatalf("user = %+v", u)
	}
}
