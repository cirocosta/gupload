package cmd

import (
	"errors"
	"strings"

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
		&cli.StringFlag{
			Name: "root-certificate",
		},
		&cli.BoolFlag{
			Name: "http2",
		},
	},
}

func uploadAction(c *cli.Context) (err error) {
	var (
		address         = c.String("address")
		file            = c.String("file")
		chunkSize       = c.Int("chunk-size")
		http2           = c.Bool("http2")
		rootCertificate = c.String("root-certificate")
		client          core.Client
	)

	if address == "" {
		must(errors.New("address"))
	}

	if file == "" {
		must(errors.New("file must be set"))
	}

	switch {
	case http2:
		if rootCertificate == "" {
			must(errors.New("http2 requires root-certificate to be supplied"))
		}

		if !strings.HasPrefix(address, "https://") {
			address = "https://" + address
		}

		http2Client, err := core.NewClientH2(core.ClientH2Config{
			Address:         address,
			RootCertificate: rootCertificate,
		})
		must(err)
		client = &http2Client
	default:
		grpcClient, err := core.NewClientGRPC(core.ClientGRPCConfig{
			Address:         address,
			RootCertificate: rootCertificate,
			ChunkSize:       chunkSize,
		})
		must(err)
		client = &grpcClient
	}

	err = client.UploadFile(context.Background(), file)
	must(err)
	defer client.Close()

	return
}
