package core

import (
	"golang.org/x/net/context"
)

type Client interface {
	UploadFile(ctx context.Context, f string) (err error)
	Close () ()
}
