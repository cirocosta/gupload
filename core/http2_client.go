package core

import (
	"os"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type ClientH2 struct {
}

type ClientH2Config struct {
	Address string
}

func NewClientH2(cfg ClientH2Config) (c ClientH2, err error) {
	if cfg.Address == "" {
		err = errors.Errorf("Address must be non-empty")
		return
	}

	return
}

func (c *ClientH2) UploadFile(ctx context.Context, f string) (err error) {
	var (
		file *os.File
	)

	file, err = os.Open(f)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to open file %s", f)
		return
	}
	defer file.Close()

	return
}

func (c *ClientH2) Close() {
	return
}
