package cmd

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

var Upload = cli.Command{
	Name:   "upload",
	Usage:  "uploads a file",
	Action: serveAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "address",
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

	fmt.Println(address)
	fmt.Println(file)

	return
}
