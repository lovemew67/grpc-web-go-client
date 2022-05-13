package grpcweb

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/grpc/encoding"
	pb "google.golang.org/grpc/encoding/proto"
)

type ClientOption func(*Client)

// Client starts each API session.
type Client struct {
	host string

	tb    TransportBuilder
	stb   StreamTransportBuilder
	codec encoding.Codec
}

// NewClient instantiates new API client for a gRPC Web API server.
// Client accepts some options to configure transports, codec, and so on.
// The default codec is Protocol Buffers.
func NewClient(host string, opts ...ClientOption) *Client {
	c := &Client{
		host: host,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.tb == nil {
		c.tb = DefaultTransportBuilder
	}

	if c.stb == nil {
		c.stb = DefaultStreamTransportBuilder
	}

	if c.codec == nil {
		// use Protocol Buffers as a default codec.
		c.codec = encoding.GetCodec(pb.Name)
	}

	return c
}

// Unary sends an unary request. (also known as simple request)
func (c *Client) Unary(ctx context.Context, req *Request, insecure bool) (*Response, error) {
	r, err := parseRequestBody(c.codec, req.in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build the request body")
	}

	rawBody, err := c.tb(c.host, req, insecure).Send(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send the request")
	}
	defer rawBody.Close()

	resBody, err := parseResponseBody(rawBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build the response body")
	}

	if err := c.codec.Unmarshal(resBody, req.out); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal response body by codec %s", c.codec.Name())
	}

	return &Response{
		ContentType: c.codec.Name(),
		Content:     req.out,
	}, nil
}

// copied from rpc_util.go#msgHeader
const headerLen = 5

func header(body []byte) []byte {
	h := make([]byte, 5)
	h[0] = byte(0)
	binary.BigEndian.PutUint32(h[1:], uint32(len(body)))
	return h
}

// header (compressed-flag(1) + message-length(4)) + body
// TODO: compressed message
func parseRequestBody(codec encoding.Codec, in interface{}) (io.Reader, error) {
	body, err := codec.Marshal(in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal the request body")
	}
	buf := bytes.NewBuffer(make([]byte, 0, headerLen+len(body)))
	buf.Write(header(body))
	buf.Write(body)
	return buf, nil
}

// copied from rpc_util#parser.recvMsg
// TODO: compressed message
func parseResponseBody(resBody io.Reader) ([]byte, error) {
	var h [5]byte
	if _, err := resBody.Read(h[:]); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(h[1:])
	if length == 0 {
		return nil, nil
	}

	// TODO: check message size

	content := make([]byte, int(length))
	if n, err := resBody.Read(content); err != nil {
		if err == io.EOF && int(n) != int(length) {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}

	return content, nil
}
