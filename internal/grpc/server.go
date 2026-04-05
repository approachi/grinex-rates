package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/approachi/grinex-rates/gen/rates"
)

type Server struct {
	grpc *grpc.Server
	port int
}

func NewServer(handler *Handler, port int) *Server {
	srv := grpc.NewServer()
	pb.RegisterRatesServiceServer(srv, handler)
	reflection.Register(srv)

	return &Server{
		grpc: srv,
		port: port,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	return s.grpc.Serve(lis)
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
