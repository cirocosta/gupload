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
	server *http.Server
	logger zerolog.Logger
}

type ServerH2Config struct {
	Port int
}

func NewServerH2(cfg ServerH2Config) (s ServerH2, err error) {
	if cfg.Port == 0 {
		err = errors.Errorf("Port must be non-zero")
		return
	}

	s.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "server_h2").
		Logger()

	s.server = &http.Server{
		Addr: ":" + strconv.Itoa(cfg.Port),
	}

	http2.ConfigureServer(s.server, nil)
	http.HandleFunc("/upload", s.Upload)

	return
}

func (s *ServerH2) Listen() (err error) {
	err = s.server.ListenAndServe()
	if err != nil {
		err = errors.Wrapf(err, "failed during server listen and serve")
		return
	}

	return
}

func (s *ServerH2) Upload(w http.ResponseWriter, r *http.Request) {
	// just receives the content and prints to stdout
	return
}

func (s *ServerH2) Close() {
	return
}
