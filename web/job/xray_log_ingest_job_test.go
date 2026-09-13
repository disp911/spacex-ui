package job

import (
	"testing"
	"time"
)

func TestParseXrayLogLine_ParsesAcceptedConnection(t *testing.T) {
	line := "2026/09/13 00:04:15.123456 from 5.142.43.223:11691 accepted tcp:api.tiktokv.com:443 [inbound-27015 -> US-gateway] email: iPhone-Stepan"

	entry, ok := parseXrayLogLine(line)
	if !ok {
		t.Fatalf("expected line to parse")
	}

	want, _ := time.ParseInLocation(xrayLogTimeFormat, "2026/09/13 00:04:15.123456", time.Local)
	if entry.Timestamp != want.UnixMicro() {
		t.Errorf("timestamp: got %d, want %d (microseconds must survive, the dedupe index relies on them)", entry.Timestamp, want.UnixMicro())
	}
	if entry.Day != "2026-09-13" {
		t.Errorf("day: got %q, want %q", entry.Day, "2026-09-13")
	}
	if entry.FromAddress != "5.142.43.223:11691" {
		t.Errorf("from: got %q", entry.FromAddress)
	}
	if entry.ToAddress != "tcp:api.tiktokv.com:443" {
		t.Errorf("to: got %q", entry.ToAddress)
	}
	if entry.Inbound != "inbound-27015" {
		t.Errorf("inbound: got %q", entry.Inbound)
	}
	if entry.Outbound != "US-gateway" {
		t.Errorf("outbound: got %q", entry.Outbound)
	}
	if entry.Email != "iPhone-Stepan" {
		t.Errorf("email: got %q", entry.Email)
	}
}

func TestParseXrayLogLine_KeepsLinesWithoutEmail(t *testing.T) {
	line := "2026/09/13 00:07:32.000001 from 5.142.43.223:8923 accepted tcp:xp.apple.com:443 [inbound-27015 -> direct]"

	entry, ok := parseXrayLogLine(line)
	if !ok {
		t.Fatalf("a line without an email is still a connection and must be kept")
	}
	if entry.Email != "" || entry.Outbound != "direct" {
		t.Errorf("got email %q outbound %q", entry.Email, entry.Outbound)
	}
}

func TestParseXrayLogLine_SkipsPanelApiTraffic(t *testing.T) {
	line := "2026/09/13 00:07:32.000001 from 127.0.0.1:50000 accepted tcp:127.0.0.1:62789 [api -> api]"

	if _, ok := parseXrayLogLine(line); ok {
		t.Fatalf("the panel's own api traffic should not be stored")
	}
}

func TestParseXrayLogLine_SkipsLinesWithoutTimestamp(t *testing.T) {
	for _, line := range []string{"", "   ", "garbage", "not a date from 1.2.3.4:5 accepted tcp:x:443 [a -> b]"} {
		if _, ok := parseXrayLogLine(line); ok {
			t.Errorf("expected %q to be skipped", line)
		}
	}
}
