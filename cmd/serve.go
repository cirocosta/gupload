package cmd

import (
	"fmt"
	"os"

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
	},
}

func must(err error) {
	if err == nil {
		return
	}

	fmt.Printf("ERROR: %+v\n", err)
	os.Exit(1)
}

func serveAction(c *cli.Context) (err error) {
	var (
		port   = c.Int("port")
		server core.Server
		http2  = c.Bool("http2")
	)

	switch {
	case http2:
		http2Server, err := core.NewServerH2(core.ServerH2Config{
			Port: port,
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
