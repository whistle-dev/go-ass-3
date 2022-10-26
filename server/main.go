package main

import (
	proto "chittychat/api/gen/go/v1"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedChatServer
	name        string
	port        string
	mu          sync.Mutex
	lclock      uint64
	clients     map[string]proto.Chat_ConnectServer
	clientNames map[string]string
}

func valueInMap(m map[string]string, value string) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}

func (s *Server) Connect(stream proto.Chat_ConnectServer) error {
	m, err := stream.Recv()
	if err != nil {
		log.Printf("Could not receive the message %v", err)
	}

	if valueInMap(s.clientNames, m.Msg) {
		log.Printf("Username already exists: %v", m.Msg)
		return errors.New("username already exists")
	}

	p, _ := peer.FromContext(stream.Context())
	s.clients[p.Addr.String()] = stream
	log.Println(m.Msg + " (" + p.Addr.String() + ") connected to the server")

	s.clientNames[p.Addr.String()] = m.Msg

	timeConnect := time.Now().Local()

	message := &proto.MsgServer{
		Name:      "Server",
		Msg:       fmt.Sprintf("Welcome to ChittyChat %v", m.Msg),
		Lclock:    s.lclock,
		Timestamp: timestamppb.New(timeConnect),
	}

	// stream.Send(message)
	for _, st := range s.clients {
		st.Send(message)
	}

	for {
		msg, err := stream.Recv()

		if err != nil {
			//if error is EOF, the client has disconnected
			if status.Code(err).String() == "Canceled" || status.Code(err).String() == "EOF" {
				log.Println(m.Msg + " (" + p.Addr.String() + ") disconnected from the server")
				break
			} else {
				log.Printf("Could not receive the message %v", err)
				break
			}
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
			Name:      s.clientNames[p.Addr.String()],
			Msg:       msg.Msg,
			Lclock:    s.lclock,
			Timestamp: timestamppb.New(time),
		}

		// stream.Send(message)
		for _, st := range s.clients {
			st.Send(message)
		}
	}

	timeLeave := time.Now().Local()

	messageLeave := &proto.MsgServer{
		Name:      "Server",
		Msg:       fmt.Sprintf("Goodbye %s", s.clientNames[p.Addr.String()]),
		Lclock:    s.lclock,
		Timestamp: timestamppb.New(timeLeave),
	}

	// stream.Send(message)
	for _, st := range s.clients {
		st.Send(messageLeave)
	}

	delete(s.clients, p.Addr.String())
	delete(s.clientNames, p.Addr.String())
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
		name:        *name,
		port:        *port,
		lclock:      0,
		clients:     make(map[string]proto.Chat_ConnectServer),
		clientNames: make(map[string]string),
	}

	go startServer(server)

	for {
		time.Sleep(1 * time.Second)
	}

}
