package grpcweb

import (
	"github.com/golang/protobuf/proto"
)

type Request struct {
	endpoint string
	in, out  interface{}
}

// NewRequest instantiates new API request from passed endpoint and I/O types.
// endpoint must be formed like:
//
//   "/{package name}.{service name}/{method name}"
//
func NewRequest(
	endpoint string,
	in proto.Message,
	out proto.Message,
) *Request {
	return &Request{
		endpoint: endpoint,
		in:       in,
		out:      out,
	}
}
