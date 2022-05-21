package ftp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	serverlib "github.com/fclairamb/ftpserverlib"

	"context"

	cosservice "github.com/Littlefisher619/cosdisk/service"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cosdisk         *cosservice.CosDisk
	logger          *logrus.Entry
	config          *Config
	nbClients       uint32
	nbClientsSync   sync.Mutex
	zeroClientEvent chan error
}

// ErrTimeout is returned when an operation timeouts
var ErrTimeout = errors.New("timeout")

// ErrNotImplemented is returned when we're using something that has not been implemented yet
// var ErrNotImplemented = errors.New("not implemented")

// ErrNotEnabled is returned when a feature hasn't been enabled
var ErrNotEnabled = errors.New("not enabled")

// NewServer creates a server instance
func NewServer(config *Config, cosdisk *cosservice.CosDisk, logger *logrus.Entry) (*Server, error) {
	return &Server{
		config:  config,
		cosdisk: cosdisk,
		logger:  logger,
	}, nil
}

// GetSettings returns some general settings around the server setup
func (s *Server) GetSettings() (*serverlib.Settings, error) {
	conf := s.config

	var portRange *serverlib.PortRange

	if conf.PassiveTransferPortRange != nil {
		portRange = &serverlib.PortRange{
			Start: conf.PassiveTransferPortRange.Start,
			End:   conf.PassiveTransferPortRange.End,
		}
	}

	return &serverlib.Settings{
		ListenAddr:               conf.ListenAddress,
		PublicHost:               conf.PublicHost,
		DisableActiveMode:        false,
		PassiveTransferPortRange: portRange,
	}, nil
}

// ClientConnected is called to send the very first welcome message
func (s *Server) ClientConnected(cc serverlib.ClientContext) (string, error) {
	s.nbClientsSync.Lock()
	defer s.nbClientsSync.Unlock()
	s.nbClients++
	s.logger.WithFields(
		logrus.Fields{
			"clientId":   cc.ID(),
			"remoteAddr": cc.RemoteAddr(),
			"nbClients":  s.nbClients,
		},
	).Info(
		"Client connected",
	)

	if !s.config.Logging.FtpExchanges {
		cc.SetDebug(true)
	}

	return "ftpserver", nil
}

// ClientDisconnected is called when the user disconnects, even if he never authenticated
func (s *Server) ClientDisconnected(cc serverlib.ClientContext) {
	s.nbClientsSync.Lock()
	defer s.nbClientsSync.Unlock()

	s.nbClients--

	s.logger.WithFields(
		logrus.Fields{
			"clientId":   cc.ID(),
			"remoteAddr": cc.RemoteAddr(),
			"nbClients":  s.nbClients,
		},
	).Info(
		"Client disconnected",
	)

	s.considerEnd()
}

// Stop will trigger a graceful stop of the server. All currently connected clients won't be disconnected instantly.
func (s *Server) Stop() {
	s.nbClientsSync.Lock()
	defer s.nbClientsSync.Unlock()
	s.zeroClientEvent = make(chan error, 1)
	s.considerEnd()
}

// WaitGracefully allows to gracefully wait for all currently connected clients before disconnecting
func (s *Server) WaitGracefully(timeout time.Duration) error {
	s.logger.Info("Waiting for last client to disconnect...")

	defer func() { s.zeroClientEvent = nil }()

	select {
	case err := <-s.zeroClientEvent:
		return err
	case <-time.After(timeout):
		return ErrTimeout
	}
}

func (s *Server) considerEnd() {
	if s.nbClients == 0 && s.zeroClientEvent != nil {
		s.zeroClientEvent <- nil
		close(s.zeroClientEvent)
	}
}

// GetTLSConfig returns a TLS Certificate to use
// The certificate could frequently change if we use something like "let's encrypt"
func (s *Server) GetTLSConfig() (*tls.Config, error) {
	// The function is called every single time a control or transfer connection requires a TLS connection. As such
	// it's important to cache it.
	return nil, nil
}

// AuthUser authenticates the user and selects an handling driver
func (s *Server) AuthUser(cc serverlib.ClientContext, user, pass string) (serverlib.ClientDriver, error) {

	gotUser, err := s.cosdisk.UserLogin(context.Background(), user, pass)
	if err != nil {
		return nil, err
	}

	logger := s.logger.WithFields(
		logrus.Fields{
			"userName":   gotUser.Id,
			"clientId":   cc.ID(),
			"remoteAddr": cc.RemoteAddr(),
		},
	)

	return &FileSystem{
		cosdisk:   s.cosdisk,
		userId:    fmt.Sprintf("%d", gotUser.Id),
		logger:    logger,
		openfiles: map[string]*File{},
		cwd:       "/",
	}, nil
}
