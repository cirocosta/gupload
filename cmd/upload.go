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
		&cli.StringFlag{
			Name: "file",
		},
	},
}

func uploadAction(c *cli.Context) (err error) {
	var (
		address = c.String("address")
		file    = c.String("file")
	)

	if address == "" {
		must(errors.New("address"))
	}

	if file == "" {
		must(errors.New("file must be set"))
	}

	client, err := core.NewClient(core.ClientConfig{
		Address: address,
	})
	must(err)

	err = client.UploadFile(context.Background(), file)
	must(err)
	defer client.Close()

	return
}
