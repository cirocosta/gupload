package core

import (
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/http2"
)

type ServerH2 struct {
	server      *http.Server
	logger      zerolog.Logger
	certificate string
	key         string
}

type ServerH2Config struct {
	Port        int
	Certificate string
	Key         string
}

func NewServerH2(cfg ServerH2Config) (s ServerH2, err error) {
	if cfg.Port == 0 {
		err = errors.Errorf("Port must be non-zero")
		return
	}

	if cfg.Certificate == "" {
		err = errors.Errorf("Certificate must be specified")
		return
	}

	if cfg.Key == "" {
		err = errors.Errorf("Key must be specified")
		return
	}

	s.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "server_h2").
		Logger()

	s.server = &http.Server{
		Addr: ":" + strconv.Itoa(cfg.Port),
	}

	s.certificate = cfg.Certificate
	s.key = cfg.Key

	http2.ConfigureServer(s.server, nil)
	http.HandleFunc("/upload", s.Upload)

	return
}

func (s *ServerH2) Listen() (err error) {
	err = s.server.ListenAndServeTLS(
		s.certificate, s.key)
	if err != nil {
		err = errors.Wrapf(err, "failed during server listen and serve")
		return
	}

	return
}

func (s *ServerH2) Upload(w http.ResponseWriter, r *http.Request) {
	// just receives the content and prints to stdout
	// read the body.

	s.logger.Info().Msg("upload received")

	return
}

func (s *ServerH2) Close() {
	return
}
