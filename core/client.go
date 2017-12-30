package core

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/cirocosta/gupload/messaging"
	"google.golang.org/grpc"
)

type Client struct {
	logger zerolog.Logger
	conn *grpc.ClientConn
	client messaging.GuploadServiceClient
}

type ClientConfig struct {
	Address string
}

func NewClient(cfg ClientConfig) (c Client, err error) {
	if cfg.Address == "" {
		err = errors.Errorf("address must be specified")
		return
	}

	c.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "client").
		Logger()

	c.conn, err = grpc.Dial(cfg.Address, grpc.WithInsecure())
	if err != nil {
		err = errors.Wrapf(err,
			"failed to start grpc connection with address %s",
			cfg.Address)
		return
	}

	c.client = messaging.NewGuploadServiceClient(c.conn)

	return
}

func (c *Client) Close () {
	if c.conn != nil {
		c.conn.Close()
	}
}

