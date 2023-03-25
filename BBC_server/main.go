// Package main implements a server for Bulletin service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	pb "github.com/ARui-tw/I2DS_Bulletin-Board-Consistency/BBC"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedBulletinServer
}

type Node struct {
	content string
	ID      uint32
	parent  *Node
	child   []*Node
}

var ID_counter uint32 = 0
var root *Node

func newPost(content string) *Node {
	ID_counter++
	var newNode *Node = &Node{content: content, ID: ID_counter, parent: nil, child: []*Node{}}
	root.child = append(root.child, newNode)
	newNode.parent = root
	return newNode
}

// SayHello implements helloworld.GreeterServer
func (s *server) Post(ctx context.Context, in *pb.Content) (*pb.PostResult, error) {
	var context string = in.GetMessage()
	log.Printf("[Post] Received: %v", context)

	newPost(context)

	return &pb.PostResult{Message: "[Success] Post ID: " + strconv.FormatUint(uint64(ID_counter), 10)}, nil
}

func (s *server) Read(ctx context.Context, in *pb.Empty) (*pb.ReadResult, error) {
	fmt.Println(len(root.child))
	log.Printf("[Read]")

	var result []string = []string{}
	var idList []uint32 = []uint32{}

	var queue []*Node = []*Node{root}
	for len(queue) != 0 {
		var node *Node = queue[0]
		queue = queue[1:]

		result = append(result, node.content)
		idList = append(idList, node.ID)

		for _, child := range node.child {
			queue = append(queue, child)
		}
	}
	var returnResult string = ""

	result = result[1:]

	for index, content := range result {
		returnResult += strconv.Itoa(index) + ": " + content
	}

	fmt.Print(returnResult)

	return &pb.ReadResult{Message: returnResult, Data: idList}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	root = &Node{content: "", ID: 0, parent: nil, child: []*Node{}}

	s := grpc.NewServer()
	pb.RegisterBulletinServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
