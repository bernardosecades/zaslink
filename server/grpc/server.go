package grpc

import (
	"fmt"
	"net"
	"os"

	sharesecretgrpc "github.com/bernardosecades/sharesecret/build"
	"github.com/bernardosecades/sharesecret/server"
	"github.com/bernardosecades/sharesecret/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type grpcServer struct {
	config          server.Config
	secretService   service.SecretService
}

func NewServer(config server.Config, ss service.SecretService) server.Server {
	return &grpcServer{config: config, secretService: ss}
}

func (s *grpcServer) Serve() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	listener, err := net.Listen(s.config.Protocol, addr)
	if err != nil {
		return err
	}

	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)

	srv := grpc.NewServer()

	serviceServer := NewShareSecretServer(s.secretService)
	sharesecretgrpc.RegisterSecretAppServer(srv, serviceServer)

	if err := srv.Serve(listener); err != nil {
		return err
	}

	return nil
}

