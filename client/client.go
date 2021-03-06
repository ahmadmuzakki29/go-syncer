package client

import (
	"context"
	"errors"
	"github.com/ahmadmuzakki29/go-syncer/pb"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var pbClient pb.SyncerClient

type Config struct {
	EndPoint    string
	LockTimeout time.Duration
}

type Client struct {
	config *Config
}

func NewClient(cfg Config) (*Client, error) {
	conn, err := grpc.Dial(cfg.EndPoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	if cfg.LockTimeout.Nanoseconds() == 0 {
		// default lock timeout is 30s
		cfg.LockTimeout = time.Duration(30) * time.Second
	}
	cli := Client{
		&cfg,
	}

	pbClient = pb.NewSyncerClient(conn)

	go func(conn *grpc.ClientConn) {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		conn.Close()
		os.Exit(1)
	}(conn)

	return &cli, nil
}

func (cli *Client) Lock(id string) error {
	config := cli.config
	if config == nil {
		return errors.New("Syncer Not initialized")
	}

	locktimeout := config.LockTimeout.String()
	reply, err := pbClient.Lock(context.Background(), &pb.LockRequest{Id: id, Locktimeout: locktimeout})
	if err != nil {
		return err
	}
	if reply.Message != "OK" {
		return errors.New(reply.Message)
	}
	return nil
}

func (cli *Client) Unlock(id string) error {
	config := cli.config
	if config == nil {
		return errors.New("Syncer Not initialized")
	}

	reply, err := pbClient.Unlock(context.Background(), &pb.LockRequest{Id: id})
	if err != nil {
		return err
	}
	if reply.Message != "OK" {
		return errors.New(reply.Message)
	}
	return nil
}
