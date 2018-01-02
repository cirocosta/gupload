<h1 align="center">gupload 📡  </h1>

<h5 align="center">Upload files with gRPC and/or HTTP2</h5>

<br/>

### Overview

`gupload` is an experiment to verify how uploading files via gRPC compares to a barebones HTTP2 server.

```
NAME:
   gupload - upload files as fast as possible

USAGE:
   gupload [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     serve    initiates a gRPC server
     upload   uploads a file
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        enables debug logging (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

Use `serve` to initiate a server (either `gRPC` or `http2` based) and `upload` to upload a file to a given address (either via `gRPC` or `http2`).


#### HTTP2

The `http2` version of both server and client require certificates / private keys. This is needed to have a well formed TLS connection between them.

The server takes both of them (certificate and private key) while the client just takes the certificate.

I've already created some certificates at `./certs` so you can just reference them as you wish (you can regenerate at any time with `make certs`).

```sh
# Start the HTTP2 server making use of HTTP2
# with the TLS configuration using the key
# and certificate from ./certs
gupload serve \
        --port 1313 \
        --http2 \
        --key ./certs/localhost.key \
        --certificate ./certs/localhost.cert


# Perform the upload of the file `./main.go` to the 
# server at `localhost:1313` via HTTP2 appending
# the self-signed certificate from `certs` to the
# root CAs.
gupload upload \
        --http2 \
        --address localhost:1313 \
        --root-certificate ./certs/localhost.cert \
        --file ./main.go
```

*note.: the certificates have `CN` (common name) set to `localhost`. As the client is not skipping insecure certificates, it'll check the address you're trying to connect to and see if it matches the certificate's CN. If you want to customize that (e.g, connect to `example.com`, make sure you issue a certificate with `CN=example.com` - or make use of SAN).*

#### GRPC

`grpc` is the default mechanism used (i.e., to make use of it you should **not** specify `--http2`) for both clients and servers.

There are two forms of running it:

- via "plain-text" TCP
- via TLS-based http2

To use the first, don't specify certificates, private keys or root certificates. To use the second, do the opposite.

For instance, to use plain tcp:

```
# Create a server
gupload serve

# Upload a file
gupload --file ./main.go
```

