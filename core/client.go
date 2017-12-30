package core

import (
	"io"
	"os"

	"github.com/cirocosta/gupload/messaging"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	logger zerolog.Logger
	conn   *grpc.ClientConn
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

func (c *Client) UploadFile(ctx context.Context, f string) (err error) {
	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
	)

	file, err = os.Open(f)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to open file %s",
			f)
		return
	}
	defer file.Close()

	stream, err := c.client.Upload(ctx)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create upload stream for file %s",
			f)
		return
	}
	defer stream.CloseSend()

	buf = make([]byte, 0, 100)
	for writing {
		n, err = file.Read(buf[:cap(buf)])
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}

			err = errors.Wrapf(err,
				"errored while copying from file to buf")
			return
		}

		buf = buf[:n]
		err = stream.Send(&messaging.Chunk{
			Content: buf,
		})
		if err != nil {
			err = errors.Wrapf(err,
				"failed to send chunk via stream")
			return
		}
	}

	return
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
