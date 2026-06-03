package rpc

import (
	"bytes"
	"io"
	"net/rpc/jsonrpc"
)

// Wrapper for jsonrpc.ServeConn that simplifies its usage
type rpcRequest struct {
	r    io.Reader
	rw   io.ReadWriter
	done chan bool
}

// Takes a reader that can be an *http.Request or anthing that implements
// io.ReadWriter interface.
// Call() will perform the jsonRPC call and write or read from the ReadWriter
func newRequest(r io.Reader) *rpcRequest {
	var buf bytes.Buffer
	done := make(chan bool)
	return &rpcRequest{r, &buf, done}
}

func (r *rpcRequest) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *rpcRequest) Write(p []byte) (n int, err error) {
	return r.rw.Write(p)
}

func (r *rpcRequest) Close() error {
	r.done <- true
	return nil
}

func (r *rpcRequest) Call() io.Reader {
	go jsonrpc.ServeConn(r)
	<-r.done
	return r.rw
}
