package core

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
)

type ClientH2 struct {
	client *http.Client
}

type ClientH2Config struct {
	RootCertificate string
	Address         string
}

func NewClientH2(cfg ClientH2Config) (c ClientH2, err error) {
	if cfg.Address == "" {
		err = errors.Errorf("Address must be non-empty")
		return
	}

	if cfg.RootCertificate == "" {
		err = errors.Errorf("RootCertificate must be specified")
		return
	}

	cert, err := ioutil.ReadFile(cfg.RootCertificate)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to read root certificate")
		return
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(cert)
	if !ok {
		err = errors.Errorf(
			"failed to root certificate %s to cert pool",
			cfg.RootCertificate)
		return
	}

	c.client = &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
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
