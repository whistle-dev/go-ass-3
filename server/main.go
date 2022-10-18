package main

import (
	proto "chittychat/api/gen/go/v1"
	"flag"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedChatServer
	name    string
	port    string
	mu      sync.Mutex
	lclock  uint64
	clients map[string]proto.Chat_ConnectServer
}

func (s *Server) Connect(stream proto.Chat_ConnectServer) error {
	p, _ := peer.FromContext(stream.Context())
	s.clients[p.Addr.String()] = stream
	log.Println(p.Addr.String())

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Could not receive the message %v", err)
			break
		}

		s.mu.Lock()
		if msg.Lclock > s.lclock {
			s.lclock = msg.Lclock
		}
		s.lclock++
		s.mu.Unlock()

		log.Printf("%v", msg)
		log.Printf("lclock: %v", s.lclock)
		// log.Printf("%d - Received message from %s: %s", s.lclock, msg.Name, msg.Msg)

		time := time.Now().Local()

		s.mu.Lock()
		s.lclock++
		s.mu.Unlock()

		message := &proto.MsgServer{
			Name:      s.name,
			Msg:       msg.Msg,
			Lclock:    s.lclock,
			Timestamp: timestamppb.New(time),
		}

		// stream.Send(message)
		for _, st := range s.clients {
			st.Send(message)
		}
	}

	delete(s.clients, p.Addr.String())
	return nil
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
		name:    *name,
		port:    *port,
		lclock:  0,
		clients: make(map[string]proto.Chat_ConnectServer),
	}

	go startServer(server)

	for {
		time.Sleep(1 * time.Second)
	}

}
