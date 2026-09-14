package telemt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// apiTimeout bounds every control API call; the API is on loopback.
const apiTimeout = 5 * time.Second

// APIClient talks to telemt's control API on 127.0.0.1.
type APIClient struct {
	baseURL string
	token   string
	http    *http.Client
}

// NewAPIClient returns a client for the API described by s.
func NewAPIClient(s Settings) *APIClient {
	return newAPIClient(fmt.Sprintf("http://127.0.0.1:%d", s.APIPort), s.APIToken)
}

func newAPIClient(baseURL, token string) *APIClient {
	return &APIClient{baseURL: baseURL, token: token, http: &http.Client{Timeout: apiTimeout}}
}

// UserStats is the live state telemt reports for one user.
type UserStats struct {
	Username           string   `json:"username"`
	CurrentConnections uint64   `json:"current_connections"`
	ActiveUniqueIPs    int      `json:"active_unique_ips"`
	ActiveIPs          []string `json:"active_unique_ips_list"`
	TotalOctets        uint64   `json:"total_octets"`
	Links              struct {
		TLS []string `json:"tls"`
	} `json:"links"`
}

// SystemInfo is the subset of /v1/system/info the panel shows.
type SystemInfo struct {
	Version       string  `json:"version"`
	UptimeSeconds float64 `json:"uptime_seconds"`
}

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Users returns live statistics and links for every configured user.
func (c *APIClient) Users(ctx context.Context) ([]UserStats, error) {
	var users []UserStats
	err := c.do(ctx, http.MethodGet, "/v1/users", &users)
	return users, err
}

// SystemInfo returns the running binary's version and uptime.
func (c *APIClient) SystemInfo(ctx context.Context) (*SystemInfo, error) {
	info := &SystemInfo{}
	if err := c.do(ctx, http.MethodGet, "/v1/system/info", info); err != nil {
		return nil, err
	}
	return info, nil
}

// ResetQuota zeroes a user's consumed traffic quota.
func (c *APIClient) ResetQuota(ctx context.Context, username string) error {
	return c.do(ctx, http.MethodPost, "/v1/users/"+url.PathEscape(username)+"/reset-quota", nil)
}

func (c *APIClient) do(ctx context.Context, method, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	if c.token != "" {
		req.Header.Set("Authorization", authHeader(c.token))
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return err
	}
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("telemt API %s %s: HTTP %d: %w", method, path, resp.StatusCode, err)
	}
	if !env.OK {
		if env.Error != nil {
			return fmt.Errorf("telemt API %s %s: %s: %s", method, path, env.Error.Code, env.Error.Message)
		}
		return fmt.Errorf("telemt API %s %s: HTTP %d", method, path, resp.StatusCode)
	}
	if out == nil || len(env.Data) == 0 {
		return nil
	}
	return json.Unmarshal(env.Data, out)
}
