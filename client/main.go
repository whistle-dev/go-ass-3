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

// ANSI codes for clearing lines
const ClearLine = "\033[2K\r"
const ClearPreviousLine = "\033[1A\033[K"

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
	fmt.Printf("Enter your username: ")

	var connectClient proto.Chat_ConnectClient

	for {
		conn, err := grpc.Dial(client.name+":"+client.port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect to the server %v", err)
		}

		chatClient := proto.NewChatClient(conn)

		connectClient, err = chatClient.Connect(context.Background())
		if err != nil {
			log.Fatalf("Could not connect to the server %v", err)
		}

		username, err := getInput()

		if err != nil {
			log.Fatalf("Could not get username: %v", err)
		}

		connectClient.Send(&proto.MsgClient{
			Name:   client.name,
			Msg:    username,
			Lclock: client.lclock,
		})

		msg, err1 := connectClient.Recv()

		if err1 == nil {
			fmt.Printf("%s [%s] (%d): %s\n", msg.Timestamp.AsTime().Local().Format("15:04:05"), msg.Name, client.lclock, msg.Msg)
			break
		}
		fmt.Println("Username already in use, please enter a new username")
	}

	/*for {
		name, err := getInput()
		if err != nil {
			log.Printf("Could not retrieve user input")
		}

	}*/

	go listenforMsg(connectClient, client)

	for {
		input, err := getInput()
		if err != nil {
			log.Fatalf("Could not retrieve user input")
		}

		if input == "exit" {
			break
		}

		fmt.Print(ClearPreviousLine)

		client.lclock++

		connectClient.Send(&proto.MsgClient{
			Name:   client.name,
			Msg:    input,
			Lclock: client.lclock,
		})

		// log.Printf("%d [%v] - Received message from %s: %s", client.lclock, msg.GetTimestamp().AsTime().Local(), msg.Name, msg.Msg)
	}

	connectClient.CloseSend()
}

func listenforMsg(connectClient proto.Chat_ConnectClient, client *Client) {
	for {
		msg, err := connectClient.Recv()
		if err != nil {
			log.Fatalf("Could not receive the message %v", err)
		}

		if msg.Lclock > client.lclock {
			client.lclock = msg.Lclock
		}

		client.lclock++

		// fmt.Print(ClearLine) to remove current line

		fmt.Printf("%s [%s] (%d): %s\n", msg.Timestamp.AsTime().Local().Format("15:04:05"), msg.Name, client.lclock, msg.Msg)

		//log.Printf("%v", msg)
		//log.Printf("lclock: %v", client.lclock)
	}
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
