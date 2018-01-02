package core

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
)

// ClientH2 provides the implementation of a file
// uploader that streams data via an HTTP2-enabled
// connection.
type ClientH2 struct {
	client  *http.Client
	address string
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

	c.address = cfg.Address

	return
}

func (c *ClientH2) UploadFile(ctx context.Context, f string) (err error) {
	var (
		file *os.File
		body = new(bytes.Buffer)
	)

	file, err = os.Open(f)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to open file %s", f)
		return
	}
	defer file.Close()

	_, err = io.Copy(body, file)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to copy file %s to body buffer",
			f)
		return
	}

	req, err := http.NewRequest("POST", c.address+"/upload", body)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create POST request")
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		err = errors.Wrapf(err,
			"request failed")
		return
	}

	if resp.StatusCode != 200 {
		err = errors.Errorf("request failed - status code: %d",
			resp.StatusCode)
		return
	}

	return
}

func (c *ClientH2) Close() {
	return
}
