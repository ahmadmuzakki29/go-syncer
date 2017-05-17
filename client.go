package syncer

import (
	"context"
	"errors"
	"github.com/ahmadmuzakki29/go-syncer/pb"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

var client pb.SyncerClient

func Init(address string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client = pb.NewSyncerClient(conn)

	go func(conn *grpc.ClientConn) {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		conn.Close()
		os.Exit(1)
	}(conn)

	return nil
}

func Lock(id string) error {
	reply, err := client.Lock(context.Background(), &pb.LockRequest{Id: id})
	if err != nil {
		return err
	}
	if reply.Message != "OK" {
		return errors.New(reply.Message)
	}
	return nil
}

func Unlock(id string) error {
	reply, err := client.Unlock(context.Background(), &pb.LockRequest{Id: id})
	if err != nil {
		return err
	}
	if reply.Message != "OK" {
		return errors.New(reply.Message)
	}
	return nil
}
