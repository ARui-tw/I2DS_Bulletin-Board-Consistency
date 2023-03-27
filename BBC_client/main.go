package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"time"

	pb "github.com/ARui-tw/I2DS_Bulletin-Board-Consistency/BBC"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	config_file = flag.String("config_file", "server_config.json", "The server config file")
	addr        string
	idList      []uint32 = []uint32{}
)

type Config struct {
	Type    string `json:"type"`
	Primary int    `json:"primary"`
	Child   []int  `json:"child"`
}

func SendPost(fileName string) {
	// read file
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Error(err)
		return
	}

	content := string(file)

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: ", err)
	}
	defer conn.Close()
	c := pb.NewBulletinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Post(ctx, &pb.Content{Message: content})
	if err != nil {
		log.Error("could not post: ", err)
	}

	if r.GetSuccess() {
		fmt.Println("Post successfully!")
	} else {
		fmt.Println("Post failed!")
	}
}

func SendRead() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBulletinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Read(ctx, &pb.Empty{})
	if err != nil {
		log.Error("could not read: %v", err)
	}

	articles := r.GetMessage()
	a := exec.Command("clear")
	a.Stdout = os.Stdout
	a.Run()
	fmt.Println("\nList of Articles:")
	fmt.Print("-----------------\n\n")
	// display 10 articles at a time and wait for user input to continue and flush before displaying the next 10
	for i := 0; i < len(articles); i++ {
		fmt.Println(articles[i])
		if i%10 == 9 {
			fmt.Print("\nPress Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')

			if i != len(articles)-1 {
				c := exec.Command("clear")
				c.Stdout = os.Stdout
				c.Run()
				fmt.Println("\nList of Articles:")
				fmt.Printf("-----------------\n\n")
			}
		}
	}

	fmt.Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	a = exec.Command("clear")
	a.Stdout = os.Stdout
	a.Run()

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
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBulletinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := c.Choose(ctx, &pb.ID{NodeID: idList[i]})
	if err != nil {
		log.Error("could not choose: %v", err)
	}

	fmt.Println("\nArticle:")
	fmt.Println("-----------------")
	fmt.Printf("%s", r.GetMessage())
}

func SendReply(fileName string, i uint32) {
	// read file
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Error(err)
		return
	}

	content := string(file)

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBulletinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Reply(ctx, &pb.Node{Message: content, NodeID: idList[i]})
	if err != nil {
		log.Error("could not reply: %v", err)
	}
	if r.GetSuccess() {
		fmt.Println("Reply successfully!")
	} else {
		fmt.Println("Reply failed!")
	}
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

	jsonFile, err := os.Open(*config_file)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	json.Unmarshal(byteValue, &config)

	addr = fmt.Sprintf("localhost:%d", config.Child[rand.Intn(len(config.Child))])

	fmt.Printf("Server Address: %s\n", addr)

	buf := bufio.NewReader(os.Stdin)

	for {
		PrintMenu()
		text, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err)
			break
		}

		switch text {
		case "1\n":
			var content string
			fmt.Print("File Name: ")
			_, err := fmt.Scanf("%s", &content)
			if err != nil {
				log.Error(err)
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
				log.Error(err)
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
				log.Error(err)
				break
			}

			fmt.Print("File Name: ")
			_, err = fmt.Scanf("%s", &content)
			if err != nil {
				log.Error(err)
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
