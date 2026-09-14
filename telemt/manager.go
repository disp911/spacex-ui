package telemt

import (
	"context"
	"os"
	"path/filepath"
	"sync"
)

// Instance is the desired state of the telemt process serving one inbound.
type Instance struct {
	InboundID int
	Tag       string
	Listen    string
	Port      int
	TLSDomain string
	// Users maps telemt user names to the panel client emails they stand for.
	Users  []User
	Emails map[string]string
}

// ClientTraffic is the traffic of one client since the previous collection.
type ClientTraffic struct {
	Email    string
	Up, Down int64
	Online   bool
}

// InboundTraffic is the traffic of one inbound since the previous collection.
type InboundTraffic struct {
	Tag      string
	Up, Down int64
	Clients  []ClientTraffic
}

// Manager runs one telemt process per mtproto inbound and keeps the set of
// processes in line with the panel's inbounds. It is safe for concurrent use.
type Manager struct {
	binary  string
	baseDir func(inboundID int) string
	// OnLog, when set, receives every output line of every process.
	OnLog func(inboundID int, line string)

	mu        sync.Mutex
	instances map[int]*instance
}

type instance struct {
	proc        *process
	desired     Instance
	metricsPort int
	runningKey  string
	config      []byte
	lastErr     string
	// seen holds the counters of the previous collection, per telemt user.
	seen map[string]UserCounters
}

// NewManager returns a manager for the bundled binary and default paths.
func NewManager() *Manager {
	return newManager(GetBinaryPath(), GetInstanceDir)
}

func newManager(binary string, baseDir func(int) string) *Manager {
	return &Manager{binary: binary, baseDir: baseDir, instances: map[int]*instance{}}
}

// Installed reports whether the telemt binary is present.
func (m *Manager) Installed() bool {
	return fileExists(m.binary)
}

// Sync makes the running processes match want. Inbounds missing from want
// are stopped. For each wanted inbound, a change in users is written to the
// config and hot-reloaded, so Telegram sessions stay up; a change to listen
// address, port or TLS domain restarts the process, as does a process that
// has exited. Errors are recorded per inbound and the first one is returned.
func (m *Manager) Sync(want []Instance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	wanted := make(map[int]bool, len(want))
	var firstErr error
	for _, w := range want {
		wanted[w.InboundID] = true
		if err := m.syncLocked(w); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	for id, inst := range m.instances {
		if !wanted[id] {
			if inst.proc != nil {
				inst.proc.stop()
			}
			delete(m.instances, id)
			_ = os.RemoveAll(m.baseDir(id))
		}
	}
	return firstErr
}

func (m *Manager) syncLocked(w Instance) error {
	inst := m.instances[w.InboundID]
	if inst == nil {
		inst = &instance{}
		m.instances[w.InboundID] = inst
	}
	inst.desired = w

	if len(w.Users) == 0 {
		// telemt cannot run without users; the inbound simply serves nobody.
		if inst.proc != nil {
			inst.proc.stop()
		}
		inst.lastErr = ""
		return nil
	}

	running := inst.proc != nil && inst.proc.running()
	if !running || inst.metricsPort == 0 {
		port, err := freeLoopbackPort()
		if err != nil {
			inst.lastErr = err.Error()
			return err
		}
		inst.metricsPort = port
	}
	s := Settings{Listen: w.Listen, Port: w.Port, TLSDomain: w.TLSDomain, MetricsPort: inst.metricsPort}
	data, err := BuildConfig(s, w.Users)
	if err != nil {
		inst.lastErr = err.Error()
		return err
	}

	dir := m.baseDir(w.InboundID)
	configPath := filepath.Join(dir, "telemt.toml")
	if running && inst.runningKey == s.restartKey() {
		if string(data) == string(inst.config) {
			return nil
		}
		if err := writeConfig(configPath, data); err != nil {
			inst.lastErr = err.Error()
			return err
		}
		inst.config = data
		return inst.proc.reload()
	}

	if inst.proc != nil {
		inst.proc.stop()
	}
	if err := writeConfig(configPath, data); err != nil {
		inst.lastErr = err.Error()
		return err
	}
	var onLine func(string)
	if m.OnLog != nil {
		id, hook := w.InboundID, m.OnLog
		onLine = func(line string) { hook(id, line) }
	}
	proc, err := startProcess(m.binary, configPath, dir, onLine)
	if err != nil {
		inst.lastErr = err.Error()
		return err
	}
	inst.proc = proc
	inst.runningKey = s.restartKey()
	inst.config = data
	inst.seen = nil
	inst.lastErr = ""
	return nil
}

// StopAll stops every process, for panel shutdown.
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, inst := range m.instances {
		if inst.proc != nil {
			inst.proc.stop()
		}
		delete(m.instances, id)
	}
}

// Running reports whether the process of an inbound is up.
func (m *Manager) Running(inboundID int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	inst := m.instances[inboundID]
	return inst != nil && inst.proc != nil && inst.proc.running()
}

// LastError returns the last start or exit error of an inbound's process.
func (m *Manager) LastError(inboundID int) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	inst := m.instances[inboundID]
	if inst == nil {
		return ""
	}
	if inst.lastErr == "" && inst.proc != nil && !inst.proc.running() && inst.proc.exitErr != nil {
		return inst.proc.exitErr.Error()
	}
	return inst.lastErr
}

// Logs returns the recent output of an inbound's current or last process.
func (m *Manager) Logs(inboundID int) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	inst := m.instances[inboundID]
	if inst == nil || inst.proc == nil {
		return nil
	}
	return inst.proc.logs.snapshot()
}

// CollectTraffic scrapes every running process and returns the traffic since
// the previous call. Counters that went backwards (a restarted process) are
// taken as new traffic from zero.
func (m *Manager) CollectTraffic(ctx context.Context) []InboundTraffic {
	type target struct {
		id     int
		port   int
		tag    string
		emails map[string]string
		proc   *process
	}
	m.mu.Lock()
	targets := make([]target, 0, len(m.instances))
	for id, inst := range m.instances {
		if inst.proc != nil && inst.proc.running() {
			targets = append(targets, target{id, inst.metricsPort, inst.desired.Tag, inst.desired.Emails, inst.proc})
		}
	}
	m.mu.Unlock()

	var result []InboundTraffic
	for _, t := range targets {
		scrapeCtx, cancel := context.WithTimeout(ctx, metricsTimeout)
		counters, err := FetchUserCounters(scrapeCtx, t.port)
		cancel()
		if err != nil {
			continue
		}

		m.mu.Lock()
		inst := m.instances[t.id]
		if inst == nil || inst.proc != t.proc {
			m.mu.Unlock()
			continue
		}
		prev := inst.seen
		inst.seen = counters
		m.mu.Unlock()

		it := InboundTraffic{Tag: t.tag}
		for user, c := range counters {
			email, ok := t.emails[user]
			if !ok {
				continue
			}
			old := prev[user]
			up, down := delta(c.Up, old.Up), delta(c.Down, old.Down)
			it.Up += up
			it.Down += down
			it.Clients = append(it.Clients, ClientTraffic{Email: email, Up: up, Down: down, Online: c.Connections > 0})
		}
		result = append(result, it)
	}
	return result
}

func delta(now, before uint64) int64 {
	if now < before {
		return int64(now)
	}
	return int64(now - before)
}
