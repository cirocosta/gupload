package core

import (
	"os"
	"strconv"
	"net"

	"github.com/cirocosta/gupload/messaging"
	"google.golang.org/grpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Server struct {
	logger zerolog.Logger
	server *grpc.Server
	port int
}

type ServerConfig struct {
	Port int
}


func NewServer(cfg ServerConfig) (s Server, err error) {
	s.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "server").
		Logger()

	if cfg.Port == 0 {
		err = errors.Errorf("Port must be specified")
		return
	}

	s.port = cfg.Port

	return
}

func (s *Server) Listen () (err error) {
	var listener net.Listener

	listener, err = net.Listen("tcp", ":" + strconv.Itoa(s.port))
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on port %d",
			s.port)
		return
	}

	s.server = grpc.NewServer()
	messaging.RegisterGuploadServiceServer(s.server, s)

	err = s.server.Serve(listener)
	if err != nil {
		err = errors.Wrapf(err, "errored listening for grpc connections")
		return
	}

	return
}

func (s *Server) Upload(stream messaging.GuploadService_UploadServer) (err error) {
	return
}

func (s *Server) Close () {
	if s.server != nil {
		s.server.Stop()
	}

	return
}

