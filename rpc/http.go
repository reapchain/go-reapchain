// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/rs/cors"
)

const (
	maxHTTPRequestContentLength = 1024 * 128
)

var nullAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:0")

type httpConn struct {
	client    *http.Client
	req       *http.Request
	closeOnce sync.Once
	closed    chan struct{}
}

// httpConn is treated specially by Client.
func (hc *httpConn) LocalAddr() net.Addr              { return nullAddr }
func (hc *httpConn) RemoteAddr() net.Addr             { return nullAddr }
func (hc *httpConn) SetReadDeadline(time.Time) error  { return nil }
func (hc *httpConn) SetWriteDeadline(time.Time) error { return nil }
func (hc *httpConn) SetDeadline(time.Time) error      { return nil }
func (hc *httpConn) Write([]byte) (int, error)        { panic("Write called") }

func (hc *httpConn) Read(b []byte) (int, error) {
	<-hc.closed
	return 0, io.EOF
}

func (hc *httpConn) Close() error {
	hc.closeOnce.Do(func() { close(hc.closed) })
	return nil
}

// DialHTTP creates a new RPC clients that connection to an RPC server over HTTP.
func DialHTTP(endpoint string) (*Client, error) {
	//log.Info("DialHTTP:yichoi debug: EOF ..")
	//client := &http.Client{}
	//data, err := EncodeClientRequest(method, req) //
	//if err != nil {
	//	return err
	//}



	req, err := http.NewRequest("POST", endpoint, nil)
	//req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data)) //
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	initctx := context.Background()
	return newClient(initctx, func(context.Context) (net.Conn, error) {
		return &httpConn{client: new(http.Client), req: req, closed: make(chan struct{})}, nil
	})
}

func (c *Client) sendHTTP(ctx context.Context, op *requestOp, msg interface{}) error {
	//log.Info("sendHTTP:yichoi debug: EOF ..")
	hc := c.writeConn.(*httpConn)
	respBody, err := hc.doRequest(ctx, msg)  //client.Do(req) ?
	if err != nil {
		return err
	}
	defer respBody.Close()  //
	var respmsg jsonrpcMessage
	//var byteslice []byte
	if err := json.NewDecoder(respBody).Decode(&respmsg); err != nil {  //?
		return err
	}
	//response, err = ioutil.ReadAll(resp.Body) // <-- https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi

	//byteslice , err = ioutil.ReadAll(respBody)  //func ReadAll(r io.Reader) ([]byte, error) {

	//respmsg = byteslice
	if err != nil {
		log.Info("yichoi debug: EOF ..")
		return err
	}


	op.resp <- &respmsg
	return nil
}


func (c *Client) sendBatchHTTP(ctx context.Context, op *requestOp, msgs []*jsonrpcMessage) error {
	//log.Info("sendBatchHTTP:yichoi debug: EOF ..")
	hc := c.writeConn.(*httpConn)
	respBody, err := hc.doRequest(ctx, msgs)
	if err != nil {
		return err
	}
	defer respBody.Close()
	var respmsgs []jsonrpcMessage
	if err := json.NewDecoder(respBody).Decode(&respmsgs); err != nil {
		return err
	}
	for _, respmsg := range respmsgs {
		op.resp <- &respmsg
	}
	return nil
}

func (hc *httpConn) doRequest(ctx context.Context, msg interface{}) (io.ReadCloser, error) {
	//log.Info("doRequest:yichoi debug: EOF ..")
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req := hc.req.WithContext(ctx)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// httpReadWriteNopCloser wraps a io.Reader and io.Writer with a NOP Close method.
type httpReadWriteNopCloser struct {
	io.Reader
	io.Writer
}

// Close does nothing and returns always nil
func (t *httpReadWriteNopCloser) Close() error {
	return nil
}

// NewHTTPServer creates a new HTTP RPC server around an API provider.
//
// Deprecated: Server implements http.Handler
func NewHTTPServer(cors []string, srv *Server) *http.Server {
	return &http.Server{Handler: newCorsHandler(srv, cors)}
}

// ServeHTTP serves JSON-RPC requests over HTTP.
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxHTTPRequestContentLength {
		http.Error(w,
			fmt.Sprintf("content length too large (%d>%d)", r.ContentLength, maxHTTPRequestContentLength),
			http.StatusRequestEntityTooLarge)
		return
	}
	w.Header().Set("content-type", "application/json")

	// create a codec that reads direct from the request body until
	// EOF and writes the response to w and order the server to process
	// a single request.
	codec := NewJSONCodec(&httpReadWriteNopCloser{r.Body, w})
	defer codec.Close()
	srv.ServeSingleRequest(codec, OptionMethodInvocation)
}

func newCorsHandler(srv *Server, allowedOrigins []string) http.Handler {
	// disable CORS support if user has not specified a custom CORS configuration
	if len(allowedOrigins) == 0 {
		return srv
	}

	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"POST", "GET"},
		MaxAge:         600,
		AllowedHeaders: []string{"*"},
	})
	return c.Handler(srv)
}

//And I got error EOF
// EncodeClientRequest encodes parameters for a JSON-RPC client request.

//func EncodeClientRequest(method string, args interface{}) ([]byte, error) {
//	c := &clientRequest{
//		Version: "2.0",
//		Method: method,
//		Params: [1]interface{}{args},
//		Id:     uint64(rand.Int63()),
//	}
//
//	return json.Marshal(c)
//}
//// DecodeClientResponse decodes the response body of a client request into // the interface reply.
//
//func DecodeClientResponse(r io.Reader, reply interface{}) error {
//	var c clientResponse
//	if err := json.NewDecoder(r).Decode(&c); err != nil {
//		return err
//	}
//	if c.Error != nil {
//		return fmt.Errorf("%v", c.Error)
//	}
//	if c.Result == nil {
//		return errors.New("result is null")
//	}
//	return json.Unmarshal(*c.Result, reply)
//}

