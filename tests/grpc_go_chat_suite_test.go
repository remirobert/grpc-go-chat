package grpc_go_chat_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestGrpcGoChat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GrpcGoChat Suite")
}
