package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

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

func getInput() (string, error) {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input = scanner.Text()
	}

	return input, nil
}

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

	/*for {
		name, err := getInput()
		if err != nil {
			log.Printf("Could not retrieve user input")
		}

	}*/

	for {
		fmt.Print("message: ")
		input, err := getInput()
		if err != nil {
			log.Fatalf("Could not retrieve user input")
		}

		if input == "exit" {
			break
		}

		client.lclock++

		connectClient.Send(&proto.MsgClient{
			Name:   client.name,
			Msg:    input,
			Lclock: client.lclock,
		})

		msg, err := connectClient.Recv()
		if err != nil {
			log.Fatalf("Could not receive the message %v", err)
		}

		if msg.Lclock > client.lclock {
			client.lclock = msg.Lclock
		}

		client.lclock++

		log.Printf("%v", msg)
		log.Printf("lclock: %v", client.lclock)

		// log.Printf("%d [%v] - Received message from %s: %s", client.lclock, msg.GetTimestamp().AsTime().Local(), msg.Name, msg.Msg)
	}

	connectClient.CloseSend()
}

func main() {
	flag.Parse() //Get the port from the command line

	client := &Client{
		name:   *name,
		port:   *port,
		lclock: 0,
	}

	startClient(client)

	// for {
	// 	time.Sleep(1 * time.Second)
	// }
}
