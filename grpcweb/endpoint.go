package grpcweb

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// ToEndpoint generates an endpoint from a service descriptor and a method descriptor.
func ToEndpoint(pkg string, s *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) string {
	return fmt.Sprintf("/%s.%s/%s", pkg, s.GetName(), m.GetName())
}
