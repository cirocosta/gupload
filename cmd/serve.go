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
		port = c.Int("port")
	)

	server, err := core.NewServer(core.ServerConfig{
		Port: port,
	})
	must(err)

	err = server.Listen()
	must(err)
	defer server.Close()

	return
}
