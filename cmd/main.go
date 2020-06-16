package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	messageapi "github.com/alexvassiliou/messenger/message"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	port := flag.String("port", "8080", "select the port")
	flag.Parse()

	var connections []*messageapi.Connection
	s := &messageapi.Server{connections}

	// Launch ther gRPC server
	log.Fatal(startServer(ctx, *port, s))
}

//
func startServer(ctx context.Context, port string, s messageapi.MessageServiceServer) error {
	// setup tcp connection
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("startServer: %s", err)
	}

	server := grpc.NewServer()

	messageapi.RegisterMessageServiceServer(server, s)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Printf("shutting down service on port :%s", port)
		}

		server.GracefulStop()

		<-ctx.Done()
	}()

	log.Printf("starting gRPC Messages service on port :%s", port)

	// Server returns an error or nil hence the return statement
	return server.Serve(listen)
}
