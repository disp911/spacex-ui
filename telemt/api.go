package telemt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// apiTimeout bounds one request to the loopback control API.
const apiTimeout = 3 * time.Second

var apiClient = &http.Client{Timeout: apiTimeout}

// UserIPs are the source addresses telemt sees for one user.
type UserIPs struct {
	// Active are the addresses with a live connection right now.
	Active []string
	// Recent are the addresses seen within telemt's recent-IP window,
	// including the active ones.
	Recent []string
}

// ErrNotRunning is returned for an inbound whose process is not up.
var ErrNotRunning = errors.New("telemt is not running for this inbound")

// UserIPs asks the control API of an inbound's process which addresses a
// telemt user connects from. A user telemt does not know has no addresses.
func (m *Manager) UserIPs(ctx context.Context, inboundID int, user string) (UserIPs, error) {
	m.mu.Lock()
	inst := m.instances[inboundID]
	var port int
	var token string
	if inst != nil && inst.proc != nil && inst.proc.running() {
		port, token = inst.apiPort, inst.apiToken
	}
	m.mu.Unlock()
	if port == 0 {
		return UserIPs{}, ErrNotRunning
	}
	return fetchUserIPs(ctx, port, token, user)
}

func fetchUserIPs(ctx context.Context, port int, token, user string) (UserIPs, error) {
	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()
	endpoint := fmt.Sprintf("http://127.0.0.1:%d/v1/users/%s", port, url.PathEscape(user))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return UserIPs{}, err
	}
	req.Header.Set("Authorization", token)
	resp, err := apiClient.Do(req)
	if err != nil {
		return UserIPs{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return UserIPs{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return UserIPs{}, fmt.Errorf("telemt API: HTTP %d", resp.StatusCode)
	}
	var body struct {
		OK   bool `json:"ok"`
		Data struct {
			Active []string `json:"active_unique_ips_list"`
			Recent []string `json:"recent_unique_ips_list"`
		} `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body); err != nil {
		return UserIPs{}, err
	}
	if !body.OK {
		return UserIPs{}, errors.New("telemt API: request failed")
	}
	return UserIPs{Active: body.Data.Active, Recent: body.Data.Recent}, nil
}

// anyEnabled reports whether at least one user may connect.
func anyEnabled(users []User) bool {
	for _, u := range users {
		if !u.Disabled {
			return true
		}
	}
	return false
}
