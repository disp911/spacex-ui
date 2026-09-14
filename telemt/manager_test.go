package telemt

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// The test binary doubles as a fake telemt: when TELEMT_FAKE_DIR is set it
// records what happens to it in that directory, serves the metrics file from
// there, and does not run tests.
func TestMain(m *testing.M) {
	if dir := os.Getenv("TELEMT_FAKE_DIR"); dir != "" {
		runFakeTelemt(dir)
		return
	}
	os.Exit(m.Run())
}

var metricsListenRe = regexp.MustCompile(`metrics_listen = '([^']+)'`)

func runFakeTelemt(dir string) {
	events, err := os.OpenFile(filepath.Join(dir, "events"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		os.Exit(2)
	}
	configPath := os.Args[len(os.Args)-1]
	sigs := make(chan os.Signal, 4)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGTERM)
	cwd, _ := os.Getwd()
	fmt.Fprintf(events, "start %d %s %s\n", os.Getpid(), configPath, cwd)
	fmt.Println("\x1b[32mINFO\x1b[0m fake telemt listening")
	if os.Getenv("TELEMT_FAKE_CRASH") == "1" {
		fmt.Fprintln(os.Stderr, "ERROR bind failed")
		os.Exit(3)
	}
	cfg, _ := os.ReadFile(configPath)
	if m := metricsListenRe.FindSubmatch(cfg); m != nil {
		ln, err := net.Listen("tcp", string(m[1]))
		if err == nil {
			go http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := os.ReadFile(filepath.Join(dir, "metrics.txt"))
				w.Write(body)
			}))
		}
	}
	for sig := range sigs {
		if sig == syscall.SIGHUP {
			fmt.Fprintf(events, "hup %d\n", os.Getpid())
			continue
		}
		os.Exit(0)
	}
}

type fakeEnv struct {
	dir string
	mgr *Manager
}

func newFakeEnv(t *testing.T) *fakeEnv {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("TELEMT_FAKE_DIR", dir)
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	mgr := newManager(exe, func(id int) string { return filepath.Join(dir, fmt.Sprintf("inbound-%d", id)) })
	t.Cleanup(mgr.StopAll)
	return &fakeEnv{dir: dir, mgr: mgr}
}

func (e *fakeEnv) lines(kind string) []string {
	data, _ := os.ReadFile(filepath.Join(e.dir, "events"))
	var out []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, kind+" ") {
			out = append(out, line)
		}
	}
	return out
}

func (e *fakeEnv) count(kind string) int { return len(e.lines(kind)) }

func (e *fakeEnv) setMetrics(t *testing.T, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(e.dir, "metrics.txt"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func eventually(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

var (
	alice = User{Name: "alice", Secret: strings.Repeat("a", 32)}
	bob   = User{Name: "bob", Secret: strings.Repeat("b", 32)}
)

func testInstance(id, port int, users ...User) Instance {
	emails := map[string]string{}
	for _, u := range users {
		emails[u.Name] = u.Name + "@mail"
	}
	return Instance{InboundID: id, Tag: fmt.Sprintf("inbound-%d", port), Port: port, TLSDomain: "www.example.com", Users: users, Emails: emails}
}

func TestManagerReloadsUserChangesAndRestartsOnSettingChanges(t *testing.T) {
	env := newFakeEnv(t)

	if err := env.mgr.Sync([]Instance{testInstance(1, 8443, alice)}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "first start", func() bool { return env.count("start") == 1 })
	if !env.mgr.Running(1) {
		t.Fatal("inbound 1 should be running")
	}
	cfgPath := filepath.Join(env.dir, "inbound-1", "telemt.toml")
	cfg, _ := os.ReadFile(cfgPath)
	if !strings.Contains(string(cfg), "alice") {
		t.Fatalf("config should contain alice:\n%s", cfg)
	}

	// Unchanged state does nothing.
	if err := env.mgr.Sync([]Instance{testInstance(1, 8443, alice)}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if env.count("start") != 1 || env.count("hup") != 0 {
		t.Fatalf("an unchanged sync must not touch the process: %v %v", env.lines("start"), env.lines("hup"))
	}

	// A user change is written in place and reloaded, not restarted.
	if err := env.mgr.Sync([]Instance{testInstance(1, 8443, alice, bob)}); err != nil {
		t.Fatal(err)
	}
	cfg, _ = os.ReadFile(cfgPath)
	if !strings.Contains(string(cfg), "bob") {
		t.Fatalf("config should contain bob after the change:\n%s", cfg)
	}
	if runtime.GOOS != "windows" {
		eventually(t, "reload signal", func() bool { return env.count("hup") == 1 })
	}
	if n := env.count("start"); n != 1 {
		t.Fatalf("user change restarted the process (%d starts)", n)
	}

	// A port change needs a new process.
	if err := env.mgr.Sync([]Instance{testInstance(1, 8444, alice, bob)}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "restart after port change", func() bool { return env.count("start") == 2 })

	// A second inbound gets its own process and directory.
	if err := env.mgr.Sync([]Instance{testInstance(1, 8444, alice, bob), testInstance(2, 9443, alice)}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "second inbound", func() bool { return env.count("start") == 3 })
	if !env.mgr.Running(1) || !env.mgr.Running(2) {
		t.Fatal("both inbounds should run")
	}

	// Dropping an inbound stops it and removes its directory; the other stays.
	if err := env.mgr.Sync([]Instance{testInstance(2, 9443, alice)}); err != nil {
		t.Fatal(err)
	}
	if env.mgr.Running(1) || !env.mgr.Running(2) {
		t.Fatal("inbound 1 must stop and inbound 2 keep running")
	}
	if _, err := os.Stat(filepath.Join(env.dir, "inbound-1")); !os.IsNotExist(err) {
		t.Fatalf("directory of a dropped inbound must be removed: %v", err)
	}

	logs := strings.Join(env.mgr.Logs(2), "\n")
	if !strings.Contains(logs, "fake telemt listening") || strings.Contains(logs, "\x1b[") {
		t.Fatalf("logs should hold the output without ANSI escapes: %q", logs)
	}
}

func TestManagerResolvesRelativePaths(t *testing.T) {
	env := newFakeEnv(t)
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	// Mirror the production layout: relative "bin" paths, with the process
	// working directory nested inside the bin folder.
	binDir := filepath.Dir(exe)
	t.Chdir(filepath.Dir(binDir))
	relBin := filepath.Base(binDir)
	env.mgr = newManager(filepath.Join(relBin, filepath.Base(exe)), func(id int) string {
		return filepath.Join(relBin, "telemt-test", fmt.Sprintf("inbound-%d", id))
	})
	t.Cleanup(env.mgr.StopAll)
	t.Cleanup(func() { os.RemoveAll(filepath.Join(binDir, "telemt-test")) })

	if err := env.mgr.Sync([]Instance{testInstance(7, 8443, alice)}); err != nil {
		t.Fatalf("relative binary and config paths must work: %v", err)
	}
	eventually(t, "start", func() bool { return env.count("start") == 1 })
	if line := env.lines("start")[0]; !strings.Contains(line, filepath.Join(binDir, "telemt-test", "inbound-7", "telemt.toml")) {
		t.Fatalf("telemt must get an absolute config path: %s", line)
	}
}

func TestManagerStopsInboundWithoutUsers(t *testing.T) {
	env := newFakeEnv(t)
	if err := env.mgr.Sync([]Instance{testInstance(1, 8443, alice)}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "start", func() bool { return env.count("start") == 1 })
	if err := env.mgr.Sync([]Instance{testInstance(1, 8443)}); err != nil {
		t.Fatalf("an inbound without active clients must stop, not fail: %v", err)
	}
	if env.mgr.Running(1) {
		t.Fatal("inbound without users must be stopped")
	}
}

func TestManagerRestartsCrashedProcessOnSync(t *testing.T) {
	env := newFakeEnv(t)
	t.Setenv("TELEMT_FAKE_CRASH", "1")
	var logMu sync.Mutex
	var logged []string
	env.mgr.OnLog = func(id int, line string) {
		logMu.Lock()
		defer logMu.Unlock()
		logged = append(logged, fmt.Sprintf("%d:%s", id, line))
	}
	if err := env.mgr.Sync([]Instance{testInstance(3, 8443, alice)}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "crash", func() bool { return !env.mgr.Running(3) && env.mgr.LastError(3) != "" })
	if logs := strings.Join(env.mgr.Logs(3), "\n"); !strings.Contains(logs, "bind failed") {
		t.Fatalf("logs of the crashed process should be kept: %q", logs)
	}

	os.Unsetenv("TELEMT_FAKE_CRASH")
	if err := env.mgr.Sync([]Instance{testInstance(3, 8443, alice)}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "restart", func() bool { return env.count("start") == 2 && env.mgr.Running(3) })
	joined := strings.Join(logged, "\n")
	if !strings.Contains(joined, "3:ERROR bind failed") || !strings.Contains(joined, "3:INFO fake telemt listening") {
		t.Fatalf("OnLog should receive every output line with the inbound id: %q", joined)
	}
}

func TestManagerRejectsInvalidConfigWithoutTouchingProcess(t *testing.T) {
	env := newFakeEnv(t)
	if err := env.mgr.Sync([]Instance{testInstance(1, 8443, alice)}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "start", func() bool { return env.count("start") == 1 })
	bad := testInstance(1, 8443, alice)
	bad.TLSDomain = ""
	if err := env.mgr.Sync([]Instance{bad}); err == nil {
		t.Fatal("invalid settings must be rejected")
	}
	if !env.mgr.Running(1) || env.count("start") != 1 || env.mgr.LastError(1) == "" {
		t.Fatal("a rejected change must leave the running proxy alone and record the error")
	}
}

func TestManagerCollectsTrafficDeltas(t *testing.T) {
	env := newFakeEnv(t)
	env.setMetrics(t, `telemt_user_octets_from_client_total{user="alice"} 100
telemt_user_octets_to_client_total{user="alice"} 1000
telemt_user_connections_current{user="alice"} 1
telemt_user_octets_from_client_total{user="stranger"} 5
`)
	if err := env.mgr.Sync([]Instance{testInstance(1, 8443, alice, bob)}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	var first []InboundTraffic
	eventually(t, "metrics", func() bool {
		first = env.mgr.CollectTraffic(ctx)
		return len(first) == 1
	})
	it := first[0]
	if it.Tag != "inbound-8443" || it.Up != 100 || it.Down != 1000 || len(it.Clients) != 1 {
		t.Fatalf("first collection = %+v", it)
	}
	if c := it.Clients[0]; c != (ClientTraffic{Email: "alice@mail", Up: 100, Down: 1000, Online: true}) {
		t.Fatalf("alice = %+v", c)
	}

	env.setMetrics(t, `telemt_user_octets_from_client_total{user="alice"} 150
telemt_user_octets_to_client_total{user="alice"} 1300
telemt_user_connections_current{user="alice"} 0
telemt_user_octets_from_client_total{user="bob"} 7
telemt_user_octets_to_client_total{user="bob"} 9
`)
	second := env.mgr.CollectTraffic(ctx)
	if len(second) != 1 {
		t.Fatalf("second collection = %+v", second)
	}
	got := map[string]ClientTraffic{}
	for _, c := range second[0].Clients {
		got[c.Email] = c
	}
	if got["alice@mail"] != (ClientTraffic{Email: "alice@mail", Up: 50, Down: 300}) {
		t.Errorf("alice delta = %+v", got["alice@mail"])
	}
	if got["bob@mail"] != (ClientTraffic{Email: "bob@mail", Up: 7, Down: 9}) {
		t.Errorf("bob delta = %+v", got["bob@mail"])
	}
	if second[0].Up != 57 || second[0].Down != 309 {
		t.Errorf("inbound delta = %d/%d", second[0].Up, second[0].Down)
	}

	// Counters that went backwards mean a fresh process: count from zero.
	env.setMetrics(t, `telemt_user_octets_from_client_total{user="alice"} 10
telemt_user_octets_to_client_total{user="alice"} 20
`)
	third := env.mgr.CollectTraffic(ctx)
	if len(third) != 1 || len(third[0].Clients) != 1 || third[0].Clients[0].Up != 10 || third[0].Clients[0].Down != 20 {
		t.Fatalf("reset counters = %+v", third)
	}
}

func TestLogBufferKeepsLastLines(t *testing.T) {
	b := newLogBuffer(3)
	fmt.Fprint(b, "one\ntwo\nthr")
	fmt.Fprint(b, "ee\nfour\n\nfive")
	if got := strings.Join(b.snapshot(), ","); got != "three,four,five" {
		t.Fatalf("snapshot = %q", got)
	}
}
