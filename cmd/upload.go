package cmd

import (
	"errors"

	"github.com/cirocosta/gupload/core"
	"golang.org/x/net/context"
	"gopkg.in/urfave/cli.v2"
)

var Upload = cli.Command{
	Name:   "upload",
	Usage:  "uploads a file",
	Action: uploadAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Value: "localhost:1313",
		},
		&cli.IntFlag{
			Name:  "chunk-size",
			Value: (1 << 12),
		},
		&cli.StringFlag{
			Name: "file",
		},
		&cli.BoolFlag{
			Name: "http2",
		},
	},
}

func uploadAction(c *cli.Context) (err error) {
	var (
		address   = c.String("address")
		file      = c.String("file")
		chunkSize = c.Int("chunk-size")
		http2     = c.Bool("http2")
		client    core.Client
	)

	if address == "" {
		must(errors.New("address"))
	}

	if file == "" {
		must(errors.New("file must be set"))
	}

	switch {
	case http2:
		http2Client, err := core.NewClientH2(core.ClientH2Config{
			Address: address,
		})
		must(err)
		client = &http2Client
	default:
		grpcClient, err := core.NewClientGRPC(core.ClientGRPCConfig{
			Address:   address,
			ChunkSize: chunkSize,
		})
		must(err)
		client = &grpcClient
	}

	err = client.UploadFile(context.Background(), file)
	must(err)
	defer client.Close()

	return
}
