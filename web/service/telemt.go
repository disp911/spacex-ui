package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/disp911/spacex-ui/v2/database"
	"github.com/disp911/spacex-ui/v2/database/model"
	"github.com/disp911/spacex-ui/v2/logger"
	"github.com/disp911/spacex-ui/v2/telemt"
	"github.com/disp911/spacex-ui/v2/xray"
)

var (
	// telemtManager is created on first use, after the panel has loaded the
	// environment that decides where the bin folder is.
	telemtManager = sync.OnceValue(func() *telemt.Manager {
		m := telemt.NewManager()
		m.OnLog = logTelemtLine
		return m
	})

	telemtState struct {
		sync.Mutex
		online    []string
		lastError map[int]string
	}
)

// TelemtService runs the telemt proxies behind mtproto inbounds and feeds
// their traffic into the same client statistics Xray inbounds use.
type TelemtService struct {
	inboundService InboundService
}

// mtprotoSettings is the settings JSON of an mtproto inbound.
type mtprotoSettings struct {
	TlsDomain string         `json:"tlsDomain"`
	Clients   []model.Client `json:"clients"`
}

// validateMTProtoInbound checks an mtproto inbound before it is saved; other
// protocols pass unchanged.
func validateMTProtoInbound(inbound *model.Inbound, clients []model.Client) error {
	if inbound.Protocol != model.MTProto {
		return nil
	}
	if !telemt.IsInstalled() {
		return errors.New("mtproto is not available: this build has no telemt binary for this platform")
	}
	var settings mtprotoSettings
	if err := json.Unmarshal([]byte(inbound.Settings), &settings); err != nil {
		return err
	}
	if !telemt.ValidTLSDomain(settings.TlsDomain) {
		return fmt.Errorf("invalid masking domain %q", settings.TlsDomain)
	}
	if inbound.Listen != "" && net.ParseIP(inbound.Listen) == nil {
		return fmt.Errorf("mtproto listen address must be an IP, got %q", inbound.Listen)
	}
	return validateMTProtoClients(inbound.Protocol, clients)
}

// validateMTProtoClients checks the secrets of clients of an mtproto inbound.
func validateMTProtoClients(protocol model.Protocol, clients []model.Client) error {
	if protocol != model.MTProto {
		return nil
	}
	for _, c := range clients {
		if !telemt.ValidSecret(c.ID) {
			return fmt.Errorf("client %q: secret must be 32 lowercase hex characters", c.Email)
		}
	}
	return nil
}

// desiredInstances builds the telemt instances every enabled mtproto inbound
// should run, with the clients that are enabled and not depleted.
func (s *TelemtService) desiredInstances() ([]telemt.Instance, error) {
	db := database.GetDB()
	var inbounds []*model.Inbound
	if err := db.Model(model.Inbound{}).Where("protocol = ? AND enable = ?", model.MTProto, true).Find(&inbounds).Error; err != nil {
		return nil, err
	}
	instances := make([]telemt.Instance, 0, len(inbounds))
	for _, inbound := range inbounds {
		var settings mtprotoSettings
		if err := json.Unmarshal([]byte(inbound.Settings), &settings); err != nil {
			logger.Warning("telemt: bad settings of inbound", inbound.Id, err)
			continue
		}
		var traffics []xray.ClientTraffic
		if err := db.Model(xray.ClientTraffic{}).Select("email, enable").Where("inbound_id = ?", inbound.Id).Find(&traffics).Error; err != nil {
			return nil, err
		}
		depleted := make(map[string]bool, len(traffics))
		for _, t := range traffics {
			if !t.Enable {
				depleted[t.Email] = true
			}
		}

		inst := telemt.Instance{
			InboundID: inbound.Id,
			Tag:       inbound.Tag,
			Listen:    inbound.Listen,
			Port:      inbound.Port,
			TLSDomain: settings.TlsDomain,
			Emails:    map[string]string{},
		}
		for _, c := range settings.Clients {
			if !c.Enable || c.Email == "" || depleted[c.Email] || !telemt.ValidSecret(c.ID) {
				continue
			}
			name := telemt.UserName(c.Email)
			if _, dup := inst.Emails[name]; dup {
				continue
			}
			inst.Emails[name] = c.Email
			inst.Users = append(inst.Users, telemt.User{Name: name, Secret: c.ID, MaxUniqueIPs: c.LimitIP})
		}
		instances = append(instances, inst)
	}
	return instances, nil
}

// Sync starts, reloads, restarts or stops telemt processes to match the
// mtproto inbounds in the database. It is cheap when nothing changed.
func (s *TelemtService) Sync() {
	want, err := s.desiredInstances()
	if err != nil {
		logger.Warning("telemt: load mtproto inbounds:", err)
		return
	}
	m := telemtManager()
	if !m.Installed() {
		if len(want) > 0 {
			reportTelemtError(0, "telemt binary is missing; mtproto inbounds cannot run on this platform")
		}
		return
	}
	_ = m.Sync(want)
	for _, w := range want {
		reportTelemtError(w.InboundID, m.LastError(w.InboundID))
	}
}

// reportTelemtError logs a telemt error once per change, not on every sync.
func reportTelemtError(inboundID int, msg string) {
	telemtState.Lock()
	defer telemtState.Unlock()
	if telemtState.lastError == nil {
		telemtState.lastError = map[int]string{}
	}
	if telemtState.lastError[inboundID] == msg {
		return
	}
	telemtState.lastError[inboundID] = msg
	if msg != "" {
		logger.Warningf("telemt (inbound %d): %s", inboundID, msg)
	}
}

// CollectTraffic records the traffic of all telemt processes since the last
// call, then re-syncs if a client ran out of quota or time.
func (s *TelemtService) CollectTraffic() {
	traffics := telemtManager().CollectTraffic(context.Background())
	var inboundTraffics []*xray.Traffic
	var clientTraffics []*xray.ClientTraffic
	online := make([]string, 0)
	for _, it := range traffics {
		if it.Up+it.Down > 0 {
			inboundTraffics = append(inboundTraffics, &xray.Traffic{IsInbound: true, Tag: it.Tag, Up: it.Up, Down: it.Down})
		}
		for _, c := range it.Clients {
			if c.Online {
				online = append(online, c.Email)
			}
			if c.Up+c.Down > 0 {
				clientTraffics = append(clientTraffics, &xray.ClientTraffic{Email: c.Email, Up: c.Up, Down: c.Down})
			}
		}
	}
	telemtState.Lock()
	telemtState.online = online
	telemtState.Unlock()

	if len(clientTraffics) == 0 && len(inboundTraffics) == 0 {
		return
	}
	_, clientsDisabled, err := s.inboundService.AddTelemtTraffic(inboundTraffics, clientTraffics)
	if err != nil {
		logger.Warning("telemt: add traffic:", err)
		return
	}
	if clientsDisabled {
		s.Sync()
	}
}

// OnlineClients returns the emails of clients with a live telemt connection.
func (s *TelemtService) OnlineClients() []string {
	telemtState.Lock()
	defer telemtState.Unlock()
	return append([]string(nil), telemtState.online...)
}

// StopAll stops every telemt process, for panel shutdown.
func (s *TelemtService) StopAll() {
	telemtManager().StopAll()
}

func logTelemtLine(inboundID int, line string) {
	msg := fmt.Sprintf("telemt (inbound %d): %s", inboundID, line)
	if strings.Contains(line, "ERROR") || strings.Contains(line, "WARN") {
		logger.Warning(msg)
	} else {
		logger.Debug(msg)
	}
}
