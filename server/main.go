package main

import (
	proto "chittychat/api/gen/go/v1"
	"flag"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedChatServer
	name   string
	port   string
	mu     sync.Mutex
	lclock uint64
}

func (s *Server) Connect(stream proto.Chat_ConnectServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Could not receive the message %v", err)
		}

		log.Printf("%d - Received message from %s: %s", msg.Lclock, msg.Name, msg.Msg)

		time := time.Now().Local()

		message := &proto.MsgServer{
			Name:      s.name,
			Msg:       "Hello " + msg.Name,
			Lclock:    s.lclock + 1,
			Timestamp: timestamppb.New(time),
		}
		stream.Send(message)
	}
}

var name = flag.String("name", "localhost", "The server name")
var port = flag.String("port", "8080", "The server port")

func startServer(server *Server) {

	listen, err := net.Listen("tcp", server.name+":"+server.port)
	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}

	grpcServer := grpc.NewServer()

	log.Printf("Starting server on port %v", server.port)

	proto.RegisterChatServer(grpcServer, server)
	serveError := grpcServer.Serve(listen)

	if serveError != nil {
		log.Fatalf("Could not start the server %v", serveError)
	}
}

func main() {
	flag.Parse() //Get the port from the command line

	server := &Server{
		name:   *name,
		port:   *port,
		lclock: 0,
	}

	go startServer(server)

	for {
		time.Sleep(1 * time.Second)
	}

}
