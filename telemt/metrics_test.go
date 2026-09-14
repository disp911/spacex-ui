package telemt

import (
	"strings"
	"testing"
)

const sampleMetrics = `# HELP telemt_user_connections_total Per-user total connections
# TYPE telemt_user_connections_total counter
telemt_ip_reservation_rollback_total{reason="tcp_limit"} 0
telemt_user_connections_total{user="alice"} 12
telemt_user_connections_current{user="alice"} 2
telemt_user_octets_from_client_total{user="alice"} 1024
telemt_user_octets_to_client_total{user="alice"} 2048
telemt_user_msgs_from_client_total{user="alice"} 7
telemt_user_connections_current{user="bob.pc"} 0
telemt_user_octets_from_client_total{user="bob.pc"} 5e+06
telemt_user_octets_to_client_total{user="bob.pc"} 1.2e+07
telemt_user_unique_ips_current{user="alice"} 1
telemt_user_octets_to_client_total{user="evil",x="y"} 99
telemt_user_octets_to_client_total{user="broken"} not-a-number
`

func TestParseUserCounters(t *testing.T) {
	got, err := parseUserCounters(strings.NewReader(sampleMetrics))
	if err != nil {
		t.Fatal(err)
	}
	if a := got["alice"]; a != (UserCounters{Up: 1024, Down: 2048, Connections: 2}) {
		t.Errorf("alice = %+v", a)
	}
	if b := got["bob.pc"]; b != (UserCounters{Up: 5000000, Down: 12000000}) {
		t.Errorf("bob.pc = %+v", b)
	}
	if _, ok := got["evil"]; ok {
		t.Error("samples with extra labels must be ignored")
	}
	if _, ok := got["broken"]; ok {
		t.Error("samples with unparsable values must be ignored")
	}
}
