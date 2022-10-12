package main

import (
	"context"
	"flag"
	"log"
	"time"

	proto "chittychat/api/gen/go/v1"

	"google.golang.org/grpc"
)

type Client struct {
	name   string
	port   string
	lclock uint64
}

var name = flag.String("name", "localhost", "The server name")
var port = flag.String("port", "8080", "The server port")

func startClient(client *Client) {

	conn, err := grpc.Dial(client.name+":"+client.port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to the server %v", err)
	}

	chatClient := proto.NewChatClient(conn)

	connectClient, err := chatClient.Connect(context.Background())
	if err != nil {
		log.Fatalf("Could not connect to the server %v", err)
	}

	connectClient.Send(&proto.MsgClient{
		Name:   client.name,
		Msg:    "Hello",
		Lclock: client.lclock + 1,
	})

	msg, err := connectClient.Recv()
	if err != nil {
		log.Fatalf("Could not receive the message %v", err)
	}

	log.Printf("%d [%v] - Received message from %s: %s", msg.Lclock, msg.GetTimestamp().AsTime().Local(), msg.Name, msg.Msg)

	connectClient.CloseSend()
}

func main() {
	flag.Parse() //Get the port from the command line

	client := &Client{
		name:   *name,
		port:   *port,
		lclock: 0,
	}

	go startClient(client)

	for {
		time.Sleep(1 * time.Second)
	}

}
