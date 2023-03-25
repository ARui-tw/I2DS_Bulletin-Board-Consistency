package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/ARui-tw/I2DS_Bulletin-Board-Consistency/BBC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func SendPost(content string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBulletinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Post(ctx, &pb.Content{Message: content})
	if err != nil {
		log.Fatalf("could not post: %v", err)
	}
	fmt.Printf("%s\n", r.GetMessage())
}

func SendRead() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBulletinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Read(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not read: %v", err)
	}
	fmt.Printf("%s", r.GetMessage())
}

func PrintMenu() {
	fmt.Println("\nMenu:")
	fmt.Println("\t1. Post")
	fmt.Println("\t2. Read")
	fmt.Println("\t3. Choose")
	fmt.Println("\t4. Reply")
	fmt.Println("\tq. Exit")
	fmt.Print("> ")
}

func main() {
	flag.Parse()

	buf := bufio.NewReader(os.Stdin)

	for {
		PrintMenu()
		text, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			break
		}

		switch text {
		case "1\n":
			fmt.Print("Content: ")
			content, err := buf.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				break
			}
			SendPost(content)
		case "2\n":
			SendRead()
		case "3\n":
			fmt.Println("Choose")
		case "4\n":
			fmt.Println("Reply")
		case "q\n":
			fmt.Println("Exit")
			return
		default:
			fmt.Println("Invalid input")
		}
	}
}
