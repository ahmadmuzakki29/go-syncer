package syncer

import (
	"errors"
	"fmt"
	"github.com/ahmadmuzakki29/go-syncer/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50051"
)

func Serve() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSyncerServer(s, &server{})

	reflection.Register(s)
	fmt.Println("serving on", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct{}

func (s *server) Lock(ctx context.Context, req *pb.LockRequest) (*pb.Reply, error) {
	id := req.Id
	if id == "" {
		return &pb.Reply{Message: "id must be specified"}, errors.New("id must be specified")
	}
	lock(id)
	return &pb.Reply{Message: "OK"}, nil
}

func (s *server) Unlock(ctx context.Context, req *pb.LockRequest) (*pb.Reply, error) {
	id := req.Id
	if id == "" {
		return &pb.Reply{Message: "id must be specified"}, errors.New("id must be specified")
	}
	unlock(id)
	return &pb.Reply{Message: "OK"}, nil
}
