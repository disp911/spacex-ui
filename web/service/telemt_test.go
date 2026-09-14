package service

import (
	"strings"
	"testing"

	"github.com/disp911/spacex-ui/v2/database/model"
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
