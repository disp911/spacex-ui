package telemt

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// metricsTimeout bounds one scrape of the loopback metrics endpoint.
const metricsTimeout = 3 * time.Second

var metricsClient = &http.Client{Timeout: metricsTimeout}

// UserCounters are the per-user counters telemt reports.
type UserCounters struct {
	// Up counts bytes received from the client, Down bytes sent to it.
	Up, Down uint64
	// Connections is the number of live connections.
	Connections uint64
}

// FetchUserCounters scrapes telemt's Prometheus metrics on the loopback port
// and returns the counters of every user it reports.
func FetchUserCounters(ctx context.Context, port int) (map[string]UserCounters, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/metrics", port), nil)
	if err != nil {
		return nil, err
	}
	resp, err := metricsClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telemt metrics: HTTP %d", resp.StatusCode)
	}
	return parseUserCounters(io.LimitReader(resp.Body, 64<<20))
}

// parseUserCounters reads the per-user series of the Prometheus text format.
func parseUserCounters(r io.Reader) (map[string]UserCounters, error) {
	counters := map[string]UserCounters{}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64<<10), 1<<20)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "telemt_user_") {
			continue
		}
		name, user, value, ok := parseUserSample(line)
		if !ok {
			continue
		}
		c := counters[user]
		switch name {
		case "telemt_user_octets_from_client_total":
			c.Up = value
		case "telemt_user_octets_to_client_total":
			c.Down = value
		case "telemt_user_connections_current":
			c.Connections = value
		default:
			continue
		}
		counters[user] = c
	}
	return counters, scanner.Err()
}

// parseUserSample parses `metric{user="name"} 123`.
func parseUserSample(line string) (name, user string, value uint64, ok bool) {
	open := strings.IndexByte(line, '{')
	end := strings.LastIndexByte(line, '}')
	if open <= 0 || end < open {
		return "", "", 0, false
	}
	name = line[:open]
	labels := line[open+1 : end]
	const prefix = `user="`
	if !strings.HasPrefix(labels, prefix) || !strings.HasSuffix(labels, `"`) || strings.Count(labels, `"`) != 2 {
		return "", "", 0, false
	}
	user = labels[len(prefix) : len(labels)-1]
	fields := strings.Fields(line[end+1:])
	if len(fields) == 0 {
		return "", "", 0, false
	}
	f, err := strconv.ParseFloat(fields[0], 64)
	if err != nil || f < 0 {
		return "", "", 0, false
	}
	return name, user, uint64(f), true
}

// freeLoopbackPort asks the kernel for an unused loopback port.
func freeLoopbackPort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port, nil
}
