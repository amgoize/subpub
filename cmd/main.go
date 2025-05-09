package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"vk/config"
	"vk/internal/subpub"

	"vk/amgoize/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type subpubServer struct {
	proto.UnimplementedSubPubServer
	subPub subpub.SubPub
}

func newSubpubServer() *subpubServer {
	return &subpubServer{
		subPub: subpub.NewSubPub(),
	}
}

func (s *subpubServer) Subscribe(req *proto.SubscribeRequest, stream proto.SubPub_SubscribeServer) error {
	log.Printf("Client subscribed with key: %s", req.GetKey())

	sub, err := s.subPub.Subscripe(req.GetKey(), func(msg interface{}) {
		event, ok := msg.(*proto.Event)
		if ok {
			if err := stream.Send(event); err != nil {
				log.Printf("Error sending event: %v", err)
			} else {
				log.Printf("Event sent to client: %v", event)
			}
		}
	})
	if err != nil {
		log.Printf("Failed to subscribe: %v", err)
		return status.Errorf(codes.Internal, "failed to subscribe: %v", err)
	}

	defer func() {
		sub.Unsubscribe()
		log.Printf("Unsubscribed from key: %s", req.GetKey())
	}()

	<-stream.Context().Done()
	log.Printf("Client disconnected: %s", req.GetKey())
	return nil
}

func (s *subpubServer) Publish(ctx context.Context, req *proto.PublishRequest) (*emptypb.Empty, error) {
	log.Printf("Received publish request with key: %s and data: %s", req.GetKey(), req.GetData())

	err := s.subPub.Publish(req.GetKey(), &proto.Event{Data: req.GetData()})
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to publish event: %v", err)
	}

	log.Printf("Event published successfully with key: %s", req.GetKey())
	return &emptypb.Empty{}, nil
}

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Log level from config (ignored): %s", cfg.Log.Level)

	server := grpc.NewServer()
	service := newSubpubServer()
	proto.RegisterSubPubServer(server, service)

	addr := fmt.Sprintf(":%d", cfg.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Printf("Server started on port %d", cfg.GRPC.Port)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh
	log.Println("Shutting down gracefully...")
	server.GracefulStop()
	log.Println("Server stopped")
}
