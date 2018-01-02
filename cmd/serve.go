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
			Value: 1313,
		},
		&cli.BoolFlag{
			Name: "http2",
		},
		&cli.StringFlag{
			Name: "key",
		},
		&cli.StringFlag{
			Name: "certificate",
		},
	},
}

func serveAction(c *cli.Context) (err error) {
	var (
		port        = c.Int("port")
		server      core.Server
		http2       = c.Bool("http2")
		key         = c.String("key")
		certificate = c.String("certificate")
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
			Port: port,
		})
		must(err)
		server = &grpcServer
	}

	err = server.Listen()
	must(err)
	defer server.Close()

	return
}
