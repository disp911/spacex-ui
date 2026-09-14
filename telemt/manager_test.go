package telemt

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

// The test binary doubles as a fake telemt: when TELEMT_FAKE_DIR is set it
// records what happens to it in that directory instead of running tests.
func TestMain(m *testing.M) {
	if dir := os.Getenv("TELEMT_FAKE_DIR"); dir != "" {
		runFakeTelemt(dir)
		return
	}
	os.Exit(m.Run())
}

func runFakeTelemt(dir string) {
	events, err := os.OpenFile(filepath.Join(dir, "events"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		os.Exit(2)
	}
	sigs := make(chan os.Signal, 4)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGTERM)
	fmt.Fprintf(events, "start %d %s\n", os.Getpid(), os.Args[len(os.Args)-1])
	fmt.Println("\x1b[32mfake telemt listening\x1b[0m")
	if os.Getenv("TELEMT_FAKE_CRASH") == "1" {
		fmt.Fprintln(os.Stderr, "fatal: bind failed")
		os.Exit(3)
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
	mgr := newManager(exe, filepath.Join(dir, "telemt.toml"), filepath.Join(dir, "work"))
	t.Cleanup(mgr.Stop)
	return &fakeEnv{dir: dir, mgr: mgr}
}

func (e *fakeEnv) events() []string {
	data, _ := os.ReadFile(filepath.Join(e.dir, "events"))
	return strings.Fields(strings.ReplaceAll(string(data), "\n", " | "))
}

func (e *fakeEnv) count(kind string) int {
	data, _ := os.ReadFile(filepath.Join(e.dir, "events"))
	n := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, kind+" ") {
			n++
		}
	}
	return n
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
	alice = User{Name: "alice", Secret: strings.Repeat("a", 32), Enabled: true}
	bob   = User{Name: "bob", Secret: strings.Repeat("b", 32), Enabled: true}
)

func TestManagerReloadsUserChangesAndRestartsOnSettingChanges(t *testing.T) {
	env := newFakeEnv(t)
	s := testSettings()

	if err := env.mgr.Apply(true, s, []User{alice}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "first start", func() bool { return env.count("start") == 1 })
	if !env.mgr.Running() || env.mgr.Crashed() {
		t.Fatal("manager should report a running, healthy process")
	}
	if st := env.mgr.Status(); !st.Running || !st.Installed || st.StartedAt.IsZero() {
		t.Fatalf("unexpected status %+v", st)
	}
	cfg, _ := os.ReadFile(env.mgr.configPath)
	if !strings.Contains(string(cfg), "alice") {
		t.Fatalf("config should contain alice:\n%s", cfg)
	}

	// A user change is written in place and reloaded, not restarted.
	if err := env.mgr.Apply(true, s, []User{alice, bob}); err != nil {
		t.Fatal(err)
	}
	cfg, _ = os.ReadFile(env.mgr.configPath)
	if !strings.Contains(string(cfg), "bob") {
		t.Fatalf("config should contain bob after the change:\n%s", cfg)
	}
	if runtime.GOOS != "windows" {
		eventually(t, "reload signal", func() bool { return env.count("hup") == 1 })
	}
	if n := env.count("start"); n != 1 {
		t.Fatalf("user change restarted the process (%d starts): %v", n, env.events())
	}

	// A port change needs a new process.
	s.Port = 8444
	if err := env.mgr.Apply(true, s, []User{alice, bob}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "restart after port change", func() bool { return env.count("start") == 2 })
	if !env.mgr.Running() {
		t.Fatal("process should be running after restart")
	}

	// Restart always restarts.
	if err := env.mgr.Restart(true, s, []User{alice, bob}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "forced restart", func() bool { return env.count("start") == 3 })

	// Disabling stops it and is not a crash.
	if err := env.mgr.Apply(false, s, []User{alice, bob}); err != nil {
		t.Fatal(err)
	}
	if env.mgr.Running() || env.mgr.Crashed() {
		t.Fatal("disabled proxy must be stopped and not reported as crashed")
	}

	logs := strings.Join(env.mgr.Logs(), "\n")
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
	// Mirror the production layout: relative "bin" paths with the process
	// working directory nested inside the bin folder.
	binDir := filepath.Dir(exe)
	t.Chdir(filepath.Dir(binDir))
	relBin := filepath.Base(binDir)
	env.mgr = newManager(filepath.Join(relBin, filepath.Base(exe)), filepath.Join(relBin, "telemt.toml"), filepath.Join(relBin, "telemt-work"))
	t.Cleanup(env.mgr.Stop)
	t.Cleanup(func() {
		os.Remove(filepath.Join(binDir, "telemt.toml"))
		os.RemoveAll(filepath.Join(binDir, "telemt-work"))
	})

	if err := env.mgr.Apply(true, testSettings(), []User{alice}); err != nil {
		t.Fatalf("relative binary and config paths must work: %v", err)
	}
	eventually(t, "start", func() bool { return env.count("start") == 1 })
	events, _ := os.ReadFile(filepath.Join(env.dir, "events"))
	if !strings.Contains(string(events), filepath.Join(filepath.Dir(exe), "telemt.toml")) {
		t.Fatalf("telemt must get an absolute config path: %s", events)
	}
}

func TestManagerStopsWithoutUsers(t *testing.T) {
	env := newFakeEnv(t)
	if err := env.mgr.Apply(true, testSettings(), []User{alice}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "start", func() bool { return env.count("start") == 1 })
	if err := env.mgr.Apply(true, testSettings(), nil); err != nil {
		t.Fatalf("removing the last user must stop the proxy, not fail: %v", err)
	}
	if env.mgr.Running() || env.mgr.Crashed() {
		t.Fatal("proxy without users must be stopped and not reported as crashed")
	}
}

func TestManagerReportsCrash(t *testing.T) {
	env := newFakeEnv(t)
	t.Setenv("TELEMT_FAKE_CRASH", "1")
	if err := env.mgr.Apply(true, testSettings(), []User{alice}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "crash", func() bool { return env.mgr.Crashed() })
	st := env.mgr.Status()
	if st.Running || st.LastError == "" {
		t.Fatalf("crashed status should carry the exit error: %+v", st)
	}
	if logs := strings.Join(env.mgr.Logs(), "\n"); !strings.Contains(logs, "bind failed") {
		t.Fatalf("logs of the crashed process should be kept: %q", logs)
	}
}

func TestManagerRejectsInvalidConfigWithoutTouchingProcess(t *testing.T) {
	env := newFakeEnv(t)
	if err := env.mgr.Apply(true, testSettings(), []User{alice}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "start", func() bool { return env.count("start") == 1 })
	bad := testSettings()
	bad.TLSDomain = ""
	if err := env.mgr.Apply(true, bad, []User{alice}); err == nil {
		t.Fatal("invalid settings must be rejected")
	}
	if !env.mgr.Running() || env.count("start") != 1 {
		t.Fatal("a rejected change must leave the running proxy alone")
	}
}

func TestLogBufferKeepsLastLines(t *testing.T) {
	b := newLogBuffer(3)
	fmt.Fprint(b, "one\ntwo\nthr")
	fmt.Fprint(b, "ee\nfour\n\nfive")
	got := strings.Join(b.snapshot(), ",")
	if got != "three,four,five" {
		t.Fatalf("snapshot = %q", got)
	}
}
