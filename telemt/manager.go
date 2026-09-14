package telemt

import (
	"sync"
	"time"
)

// Status describes the proxy process for the panel.
type Status struct {
	Installed bool      `json:"installed"`
	Running   bool      `json:"running"`
	StartedAt time.Time `json:"startedAt"`
	LastError string    `json:"lastError"`
}

// Manager owns the telemt process and keeps it in line with the panel's
// desired state. It is safe for concurrent use.
type Manager struct {
	binary     string
	configPath string
	workDir    string

	mu         sync.Mutex
	proc       *process
	runningKey string
	wanted     bool
	lastErr    string
}

// NewManager returns a manager for the bundled binary and default paths.
func NewManager() *Manager {
	return newManager(GetBinaryPath(), GetConfigPath(), GetWorkDir())
}

func newManager(binary, configPath, workDir string) *Manager {
	return &Manager{binary: binary, configPath: configPath, workDir: workDir}
}

// Apply makes the process match the desired state. A disabled proxy, or one
// without users, is stopped. When only users changed, the running process
// reloads its config in place so existing Telegram sessions stay up; a change
// to port, TLS domain, public host or API credentials restarts it.
func (m *Manager) Apply(enabled bool, s Settings, users []User) error {
	return m.apply(enabled, s, users, false)
}

// Restart is Apply that always restarts a running process.
func (m *Manager) Restart(enabled bool, s Settings, users []User) error {
	return m.apply(enabled, s, users, true)
}

func (m *Manager) apply(enabled bool, s Settings, users []User, force bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !enabled || len(users) == 0 {
		m.wanted = false
		m.stopLocked()
		m.lastErr = ""
		return nil
	}

	data, err := BuildConfig(s, users)
	if err != nil {
		return err
	}
	m.wanted = true

	if !force && m.proc != nil && m.proc.running() && m.runningKey == s.restartKey() {
		if err := writeConfig(m.configPath, data); err != nil {
			return err
		}
		return m.proc.reload()
	}

	m.stopLocked()
	if err := writeConfig(m.configPath, data); err != nil {
		m.lastErr = err.Error()
		return err
	}
	proc, err := startProcess(m.binary, m.configPath, m.workDir)
	if err != nil {
		m.lastErr = err.Error()
		return err
	}
	m.proc = proc
	m.runningKey = s.restartKey()
	m.lastErr = ""
	return nil
}

// Stop terminates the process, for example on panel shutdown.
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wanted = false
	m.stopLocked()
}

func (m *Manager) stopLocked() {
	if m.proc != nil {
		m.proc.stop()
	}
}

// Crashed reports whether the process should be running but is not.
func (m *Manager) Crashed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.wanted && (m.proc == nil || !m.proc.running())
}

// Running reports whether the process is up.
func (m *Manager) Running() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.proc != nil && m.proc.running()
}

// Status returns the current process state.
func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()
	st := Status{Installed: fileExists(m.binary), LastError: m.lastErr}
	if m.proc == nil {
		return st
	}
	if m.proc.running() {
		st.Running = true
		st.StartedAt = m.proc.startedAt
	} else if m.wanted && st.LastError == "" && m.proc.exitErr != nil {
		st.LastError = m.proc.exitErr.Error()
	}
	return st
}

// Logs returns the recent output of the current or last process.
func (m *Manager) Logs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.proc == nil {
		return nil
	}
	return m.proc.logs.snapshot()
}
