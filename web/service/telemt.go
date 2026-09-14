package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/logger"
	"github.com/mhsanaei/3x-ui/v2/telemt"
	"github.com/mhsanaei/3x-ui/v2/util/random"

	"gorm.io/gorm"
)

var (
	// telemtManager is created on first use, after the panel has loaded the
	// environment that decides where the bin folder is.
	telemtManager = sync.OnceValue(telemt.NewManager)
	// telemtAPIToken authorizes the panel against telemt's loopback API. A
	// fresh token per panel process is enough: telemt runs as a child of the
	// panel and gets the token in the config written when it starts.
	telemtAPIToken = random.Seq(32)
)

// telemtAPICallTimeout bounds the live statistics lookup of one page load.
const telemtAPICallTimeout = 3 * time.Second

// TelemtService manages the bundled Telegram MTProto proxy: its settings,
// its users and the telemt process.
type TelemtService struct {
	settingService SettingService
}

// TelemtSettings are the proxy options edited on the Telegram page.
type TelemtSettings struct {
	Enable     bool   `json:"enable" form:"enable"`
	Port       int    `json:"port" form:"port"`
	TlsDomain  string `json:"tlsDomain" form:"tlsDomain"`
	PublicHost string `json:"publicHost" form:"publicHost"`
}

// TelemtUserView is a user row enriched with live proxy statistics.
type TelemtUserView struct {
	model.TelemtUser
	Link        string   `json:"link"`
	Connections uint64   `json:"connections"`
	ActiveIps   []string `json:"activeIps"`
	Traffic     uint64   `json:"traffic"`
}

// TelemtOverview is everything the Telegram page shows.
type TelemtOverview struct {
	Settings  TelemtSettings   `json:"settings"`
	Installed bool             `json:"installed"`
	Running   bool             `json:"running"`
	Version   string           `json:"version"`
	StartedAt int64            `json:"startedAt"`
	LastError string           `json:"lastError"`
	Users     []TelemtUserView `json:"users"`
}

// GetSettings reads the proxy settings.
func (s *TelemtService) GetSettings() (TelemtSettings, error) {
	var ts TelemtSettings
	var err error
	if ts.Enable, err = s.settingService.getBool("telemtEnable"); err != nil {
		return ts, err
	}
	if ts.Port, err = s.settingService.getInt("telemtPort"); err != nil {
		return ts, err
	}
	if ts.TlsDomain, err = s.settingService.getString("telemtTlsDomain"); err != nil {
		return ts, err
	}
	if ts.PublicHost, err = s.settingService.getString("telemtPublicHost"); err != nil {
		return ts, err
	}
	return ts, nil
}

// SaveSettings validates and stores the proxy settings, then applies them.
func (s *TelemtService) SaveSettings(ts TelemtSettings) error {
	ts.TlsDomain = strings.ToLower(strings.TrimSpace(ts.TlsDomain))
	ts.PublicHost = strings.TrimSpace(ts.PublicHost)

	if ts.Enable && !telemt.IsInstalled() {
		return errors.New("telemt is not bundled for this platform")
	}
	if ts.Enable || ts.TlsDomain != "" {
		rs, err := s.runtimeSettings(ts)
		if err != nil {
			return err
		}
		if err := rs.Validate(); err != nil {
			return err
		}
	}
	if err := s.checkPortFree(ts.Port); err != nil {
		return err
	}

	if err := s.settingService.setBool("telemtEnable", ts.Enable); err != nil {
		return err
	}
	if err := s.settingService.setInt("telemtPort", ts.Port); err != nil {
		return err
	}
	if err := s.settingService.setString("telemtTlsDomain", ts.TlsDomain); err != nil {
		return err
	}
	if err := s.settingService.setString("telemtPublicHost", ts.PublicHost); err != nil {
		return err
	}
	return s.Apply()
}

// checkPortFree rejects a proxy port the panel already uses elsewhere.
func (s *TelemtService) checkPortFree(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid proxy port %d", port)
	}
	if webPort, err := s.settingService.GetPort(); err == nil && webPort == port {
		return fmt.Errorf("port %d is used by the panel", port)
	}
	if subEnable, err := s.settingService.GetSubEnable(); err == nil && subEnable {
		if subPort, err := s.settingService.GetSubPort(); err == nil && subPort == port {
			return fmt.Errorf("port %d is used by the subscription server", port)
		}
	}
	if apiPort, err := s.settingService.getInt("telemtApiPort"); err == nil && apiPort == port {
		return fmt.Errorf("port %d is reserved for the proxy API", port)
	}
	var count int64
	if err := database.GetDB().Model(&model.Inbound{}).Where("port = ?", port).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("port %d is used by an inbound", port)
	}
	return nil
}

func (s *TelemtService) runtimeSettings(ts TelemtSettings) (telemt.Settings, error) {
	apiPort, err := s.settingService.getInt("telemtApiPort")
	if err != nil {
		return telemt.Settings{}, err
	}
	return telemt.Settings{
		Port:       ts.Port,
		TLSDomain:  ts.TlsDomain,
		PublicHost: ts.PublicHost,
		APIPort:    apiPort,
		APIToken:   telemtAPIToken,
	}, nil
}

func (s *TelemtService) getUsers() ([]model.TelemtUser, error) {
	var users []model.TelemtUser
	err := database.GetDB().Order("id asc").Find(&users).Error
	return users, err
}

// desiredState loads everything the proxy process should run with.
func (s *TelemtService) desiredState() (bool, telemt.Settings, []telemt.User, error) {
	ts, err := s.GetSettings()
	if err != nil {
		return false, telemt.Settings{}, nil, err
	}
	rs, err := s.runtimeSettings(ts)
	if err != nil {
		return false, telemt.Settings{}, nil, err
	}
	rows, err := s.getUsers()
	if err != nil {
		return false, telemt.Settings{}, nil, err
	}
	users := make([]telemt.User, 0, len(rows))
	for _, row := range rows {
		u := telemt.User{
			Name:         row.Username,
			Secret:       row.Secret,
			Enabled:      row.Enable,
			QuotaBytes:   row.TotalBytes,
			MaxUniqueIPs: row.LimitIp,
		}
		if row.ExpiryTime > 0 {
			u.ExpiresAt = time.UnixMilli(row.ExpiryTime)
		}
		users = append(users, u)
	}
	return ts.Enable && telemt.IsInstalled(), rs, users, nil
}

// Apply brings the proxy process in line with the stored settings and users.
func (s *TelemtService) Apply() error {
	enabled, rs, users, err := s.desiredState()
	if err != nil {
		return err
	}
	return telemtManager().Apply(enabled, rs, users)
}

// Restart restarts the proxy process with the stored settings and users.
func (s *TelemtService) Restart() error {
	enabled, rs, users, err := s.desiredState()
	if err != nil {
		return err
	}
	return telemtManager().Restart(enabled, rs, users)
}

// Stop stops the proxy process, for panel shutdown.
func (s *TelemtService) Stop() {
	telemtManager().Stop()
}

// RestartIfCrashed restarts the proxy when it should run but has exited.
func (s *TelemtService) RestartIfCrashed() {
	if !telemtManager().Crashed() {
		return
	}
	logger.Warning("telemt is not running, restarting it")
	if err := s.Restart(); err != nil {
		logger.Error("restart telemt failed:", err)
	}
}

// Logs returns the recent output of the proxy process.
func (s *TelemtService) Logs() []string {
	return telemtManager().Logs()
}

// GetOverview returns settings, process state and users with live stats.
func (s *TelemtService) GetOverview() (*TelemtOverview, error) {
	ts, err := s.GetSettings()
	if err != nil {
		return nil, err
	}
	rows, err := s.getUsers()
	if err != nil {
		return nil, err
	}
	st := telemtManager().Status()
	ov := &TelemtOverview{
		Settings:  ts,
		Installed: st.Installed,
		Running:   st.Running,
		LastError: st.LastError,
		Users:     make([]TelemtUserView, 0, len(rows)),
	}
	if st.Running {
		ov.StartedAt = st.StartedAt.UnixMilli()
	}

	stats := map[string]telemt.UserStats{}
	if st.Running {
		rs, err := s.runtimeSettings(ts)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), telemtAPICallTimeout)
			defer cancel()
			client := telemt.NewAPIClient(rs)
			if info, err := client.SystemInfo(ctx); err == nil {
				ov.Version = info.Version
			}
			if list, err := client.Users(ctx); err == nil {
				for _, u := range list {
					stats[u.Username] = u
				}
			} else {
				logger.Debug("telemt users API:", err)
			}
		}
	}

	for _, row := range rows {
		view := TelemtUserView{TelemtUser: row}
		if u, ok := stats[row.Username]; ok {
			view.Connections = u.CurrentConnections
			view.ActiveIps = u.ActiveIPs
			view.Traffic = u.TotalOctets
			if len(u.Links.TLS) > 0 {
				view.Link = u.Links.TLS[0]
			}
		}
		if view.Link == "" && ts.PublicHost != "" && ts.TlsDomain != "" {
			view.Link = telemt.Link(ts.PublicHost, ts.Port, row.Secret, ts.TlsDomain)
		}
		ov.Users = append(ov.Users, view)
	}
	return ov, nil
}

func newTelemtSecret() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func normalizeTelemtUser(u *model.TelemtUser) error {
	u.Username = strings.TrimSpace(u.Username)
	u.Secret = strings.ToLower(strings.TrimSpace(u.Secret))
	u.Comment = strings.TrimSpace(u.Comment)
	if !telemt.ValidUserName(u.Username) {
		return errors.New("user name must be 1-64 characters: letters, digits, '_', '.' or '-'")
	}
	if u.Secret != "" && !telemt.ValidSecret(u.Secret) {
		return errors.New("secret must be 32 hex characters")
	}
	if u.ExpiryTime < 0 || u.TotalBytes < 0 || u.LimitIp < 0 {
		return errors.New("limits cannot be negative")
	}
	return nil
}

func (s *TelemtService) usernameTaken(username string, exceptID int) (bool, error) {
	var count int64
	err := database.GetDB().Model(&model.TelemtUser{}).
		Where("username = ? AND id <> ?", username, exceptID).Count(&count).Error
	return count > 0, err
}

// AddUser creates a proxy user, generating a secret when none is given.
func (s *TelemtService) AddUser(u *model.TelemtUser) error {
	if err := normalizeTelemtUser(u); err != nil {
		return err
	}
	if taken, err := s.usernameTaken(u.Username, 0); err != nil {
		return err
	} else if taken {
		return fmt.Errorf("user %q already exists", u.Username)
	}
	if u.Secret == "" {
		secret, err := newTelemtSecret()
		if err != nil {
			return err
		}
		u.Secret = secret
	}
	u.Id = 0
	if err := database.GetDB().Create(u).Error; err != nil {
		return err
	}
	return s.Apply()
}

// UpdateUser replaces the editable fields of the user with id.
func (s *TelemtService) UpdateUser(id int, u *model.TelemtUser) error {
	if err := normalizeTelemtUser(u); err != nil {
		return err
	}
	var existing model.TelemtUser
	if err := database.GetDB().First(&existing, id).Error; err != nil {
		return err
	}
	if taken, err := s.usernameTaken(u.Username, id); err != nil {
		return err
	} else if taken {
		return fmt.Errorf("user %q already exists", u.Username)
	}
	if u.Secret == "" {
		u.Secret = existing.Secret
	}
	err := database.GetDB().Model(&existing).
		Updates(map[string]any{
			"username":    u.Username,
			"secret":      u.Secret,
			"enable":      u.Enable,
			"expiry_time": u.ExpiryTime,
			"total_bytes": u.TotalBytes,
			"limit_ip":    u.LimitIp,
			"comment":     u.Comment,
		}).Error
	if err != nil {
		return err
	}
	return s.Apply()
}

// SetUserEnable switches one user on or off.
func (s *TelemtService) SetUserEnable(id int, enable bool) error {
	res := database.GetDB().Model(&model.TelemtUser{}).Where("id = ?", id).Update("enable", enable)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return s.Apply()
}

// DeleteUser removes the user with id.
func (s *TelemtService) DeleteUser(id int) error {
	res := database.GetDB().Delete(&model.TelemtUser{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return s.Apply()
}

// RotateSecret gives the user a new secret, invalidating the old link.
func (s *TelemtService) RotateSecret(id int) error {
	secret, err := newTelemtSecret()
	if err != nil {
		return err
	}
	res := database.GetDB().Model(&model.TelemtUser{}).Where("id = ?", id).Update("secret", secret)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return s.Apply()
}

// ResetTraffic zeroes the user's consumed quota in the running proxy.
func (s *TelemtService) ResetTraffic(id int) error {
	var user model.TelemtUser
	if err := database.GetDB().First(&user, id).Error; err != nil {
		return err
	}
	if !telemtManager().Running() {
		return errors.New("telemt is not running")
	}
	ts, err := s.GetSettings()
	if err != nil {
		return err
	}
	rs, err := s.runtimeSettings(ts)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), telemtAPICallTimeout)
	defer cancel()
	return telemt.NewAPIClient(rs).ResetQuota(ctx, user.Username)
}
