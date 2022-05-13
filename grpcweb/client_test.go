package grpcweb

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/ktr0731/grpc-test/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// var defaultAddr = "localhost:8880"
var defaultAddr4 = "grpcweb.howardzhou.dev:443"

type protoHelper struct {
	*desc.FileDescriptor

	s map[string]*descriptor.ServiceDescriptorProto
	m map[string]*dynamic.Message
}

func (h *protoHelper) getServiceByName(t *testing.T, n string) *descriptor.ServiceDescriptorProto {
	if h.s == nil {
		h.s = map[string]*descriptor.ServiceDescriptorProto{}
		for _, svc := range h.GetServices() {
			h.s[svc.GetName()] = svc.AsServiceDescriptorProto()
		}
	}
	svc, ok := h.s[n]
	if !ok {
		require.FailNowf(t, "ServiceDescriptor not found", "no such *desc.ServiceDescriptor: %s", n)
	}
	return svc
}

func (h *protoHelper) getMessageTypeByName(t *testing.T, n string) *dynamic.Message {
	if h.m == nil {
		h.m = map[string]*dynamic.Message{}
		for _, msg := range h.GetMessageTypes() {
			h.m[msg.GetName()] = dynamic.NewMessage(msg)
		}
	}
	msg, ok := h.m[n]
	if !ok {
		require.FailNowf(t, "MessageDescriptor not found", "no such *desc.MessageDescriptor: %s", n)
	}
	return msg
}

func getAPIProto(t *testing.T) *protoHelper {
	t.Helper()

	pkgs := parseProto(t, "echo.proto")
	require.Len(t, pkgs, 1)

	return &protoHelper{FileDescriptor: pkgs[0]}
}

func Test_ClientE2E(t *testing.T) {
	pkg := getAPIProto(t)
	service := pkg.getServiceByName(t, "EchoService")

	t.Run("Unary", func(t *testing.T) {
		defer server.New().Serve().Stop()

		client := NewClient(defaultAddr4)
		endpoint := ToEndpoint("echo", service, service.GetMethod()[0])

		in := pkg.getMessageTypeByName(t, "HiRequest")

		cases := []string{
			time.Now().UTC().String(),
			time.Now().UTC().String(),
		}

		for _, c := range cases {
			in.SetFieldByName("message", c)

			out := pkg.getMessageTypeByName(t, "HiResponse")

			req := NewRequest(endpoint, in, out)
			res, err := client.Unary(context.Background(), req)
			assert.NoError(t, err)

			assert.Equal(t, c, extractMessage(t, res))
			// expected := fmt.Sprintf("hello, %s", c)
			// assert.Equal(t, expected, extractMessage(t, res))
		}
	})
}

func extractMessage(t *testing.T, res *Response) string {
	require.NotNil(t, res.Content)

	m, ok := res.Content.(*dynamic.Message)
	require.True(t, ok)

	msg := m.GetFieldByName("message")
	s, ok := msg.(string)
	require.True(t, ok)

	return s
}
