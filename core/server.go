package core

import (
	"io"
	"net"
	"os"
	"strconv"

	"github.com/cirocosta/gupload/messaging"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Server struct {
	logger zerolog.Logger
	server *grpc.Server
	port   int
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

func (s *Server) Listen() (err error) {
	var listener net.Listener

	listener, err = net.Listen("tcp", ":"+strconv.Itoa(s.port))
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
	var (
		in *messaging.Chunk
	)

	for {
		in, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}

			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return
		}

		s.logger.Info().
			Interface("chunk", in.GetContent()).
			Msg("message received")
	}

	return
}

func (s *Server) Close() {
	if s.server != nil {
		s.server.Stop()
	}

	return
}
