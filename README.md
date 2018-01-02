<h1 align="center">gupload 📡  </h1>

<h5 align="center">Upload files with gRPC and/or HTTP2</h5>

<br/>

### Overview

`gupload` is an experiment to verify how uploading files via gRPC compares to a barebones HTTP2 server.

#### HTTP2 server and client

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

