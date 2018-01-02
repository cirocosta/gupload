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

// ClientGRPC provides the implementation of a file
// uploader that streams chunks via protobuf-encoded
// messages.
type ClientGRPC struct {
	logger    zerolog.Logger
	conn      *grpc.ClientConn
	client    messaging.GuploadServiceClient
	chunkSize int
}

type ClientGRPCConfig struct {
	Address   string
	ChunkSize int
}

func NewClientGRPC(cfg ClientGRPCConfig) (c ClientGRPC, err error) {
	if cfg.Address == "" {
		err = errors.Errorf("address must be specified")
		return
	}

	switch {
	case cfg.ChunkSize == 0:
		err = errors.Errorf("ChunkSize must be specified")
		return
	case cfg.ChunkSize > (1 << 20):
		err = errors.Errorf("ChunkSize must be < than 1MB")
		return
	default:
		c.chunkSize = cfg.ChunkSize
	}

	c.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "client").
		Logger()

	// TODO replace this by non-deprecated apis
	c.conn, err = grpc.Dial(cfg.Address,
		grpc.WithInsecure(),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()))
	if err != nil {
		err = errors.Wrapf(err,
			"failed to start grpc connection with address %s",
			cfg.Address)
		return
	}

	c.client = messaging.NewGuploadServiceClient(c.conn)

	return
}

func (c *ClientGRPC) UploadFile(ctx context.Context, f string) (err error) {
	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
		status  *messaging.UploadStatus
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

	status, err = stream.CloseAndRecv()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to receive upstream status response")
		return
	}

	if status.Code != messaging.UploadStatusCode_Ok {
		err = errors.Errorf(
			"upload failed - msg: %s",
			status.Message)
		return
	}

	return
}

func (c *ClientGRPC) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
