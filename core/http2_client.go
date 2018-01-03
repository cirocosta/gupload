package core

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"time"

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

func (c *ClientH2) UploadFile(ctx context.Context, f string) (stats Stats, err error) {
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

	req, err := http.NewRequest("POST", c.address+"/upload", file)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create POST request")
		return
	}

	stats.StartedAt = time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		err = errors.Wrapf(err,
			"request failed")
		return
	}
	stats.FinishedAt = time.Now()

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
