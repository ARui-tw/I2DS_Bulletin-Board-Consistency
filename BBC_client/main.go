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

var idList []uint32 = []uint32{}

var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

func SendPost(fileName string) {
	// read file
	file, err := os.ReadFile(fileName)
	if err != nil {
		Error.Println(err)
		return
	}

	content := string(file)

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		Error.Println("did not connect:", err)
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

	fmt.Println("\nList of Articles:")
	fmt.Println("-----------------")
	fmt.Printf("%s", r.GetMessage())

	idList = append(r.GetData())

	/*
		// print idList
		fmt.Println("ID List:")
		for i := 0; i < len(idList); i++ {
			fmt.Println(idList[i])
		}
	*/
}

func SendChoose(i uint32) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBulletinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Choose(ctx, &pb.ChooseMessage{NodeID: idList[i]})
	if err != nil {
		log.Fatalf("could not choose: %v", err)
	}

	fmt.Println("\nArticle:")
	fmt.Println("-----------------")
	fmt.Printf("%s", r.GetMessage())
}

func SendReply(fileName string, i uint32) {
	// read file
	file, err := os.ReadFile(fileName)
	if err != nil {
		Error.Println(err)
		return
	}

	content := string(file)
	// Set up a connection to the server.

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBulletinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Reply(ctx, &pb.ReplyMessage{Message: content, NodeID: idList[i]})
	if err != nil {
		log.Fatalf("could not reply: %v", err)
	}
	fmt.Printf("%s\n", r.GetMessage())
}

func PrintMenu() {
	fmt.Println("-----------------")
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
			var content string
			fmt.Print("File Name: ")
			_, err := fmt.Scanf("%s", &content)
			if err != nil {
				fmt.Println(err)
				break
			}
			SendPost(content)
		case "2\n":
			SendRead()
		case "3\n":
			if len(idList) == 0 {
				fmt.Println("No article to choose, please read first")
				break
			}

			var i uint32
			fmt.Print("ID: ")
			_, err := fmt.Scanf("%d", &i)
			if err != nil {
				fmt.Println(err)
				break
			}

			SendChoose(i)

		case "4\n":
			var i uint32
			var content string

			if len(idList) == 0 {
				fmt.Println("No article to reply, please read first")
				break
			}

			fmt.Print("ID: ")
			_, err := fmt.Scanf("%d", &i)
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Print("File Name: ")
			_, err = fmt.Scanf("%s", &content)
			if err != nil {
				fmt.Println(err)
				break
			}

			SendReply(content, i)

		case "q\n":
			fmt.Println("Exit")
			return
		default:
			fmt.Println("Invalid input")
		}
	}
}
