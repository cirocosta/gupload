package cmd

import (
	"errors"

	"github.com/cirocosta/gupload/core"
	"gopkg.in/urfave/cli.v2"
)

var Serve = cli.Command{
	Name:   "serve",
	Usage:  "initiates a gRPC server",
	Action: serveAction,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Usage: "port to bind to",
			Value: 1313,
		},
		&cli.BoolFlag{
			Name:  "http2",
			Usage: "whether or not to serve via HTTP2 instead of gRPC",
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "path to TLS certificate",
		},
		&cli.StringFlag{
			Name:  "certificate",
			Usage: "path to TLS certificate",
		},
		&cli.StringFlag{
			Name:  "compress",
			Usage: "whether or not to enable compression",
		},
	},
}

func serveAction(c *cli.Context) (err error) {
	var (
		port        = c.Int("port")
		http2       = c.Bool("http2")
		key         = c.String("key")
		certificate = c.String("certificate")
		server      core.Server
	)

	switch {
	case http2:
		if key == "" || certificate == "" {
			must(errors.New(
				"http2 requires key and certificate to be specified"))
		}

		http2Server, err := core.NewServerH2(core.ServerH2Config{
			Port:        port,
			Certificate: certificate,
			Key:         key,
		})
		must(err)
		server = &http2Server
	default:
		grpcServer, err := core.NewServerGRPC(core.ServerGRPCConfig{
			Port:        port,
			Certificate: certificate,
			Key:         key,
		})
		must(err)
		server = &grpcServer
	}

	err = server.Listen()
	must(err)
	defer server.Close()

	return
}
