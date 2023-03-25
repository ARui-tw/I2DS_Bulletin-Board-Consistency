// Package main implements a server for Bulletin service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

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

var Nodes []*Node = []*Node{}

func newPost(content string) {
	ID_counter++
	var newNode *Node = &Node{content: content, ID: ID_counter, parent: root, child: []*Node{}}
	root.child = append(root.child, newNode)
	Nodes = append(Nodes, newNode)
}

func newReply(content string, parentID uint32) {
	var parentNode *Node = Nodes[parentID]
	ID_counter++
	var newNode *Node = &Node{content: content, ID: ID_counter, parent: parentNode, child: []*Node{}}
	parentNode.child = append(parentNode.child, newNode)
	Nodes = append(Nodes, newNode)
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

	// DFS to get all the content and append tap for each level
	var DFS func(node *Node, level int)
	DFS = func(node *Node, level int) {
		var tap string = ""

		var i int
		for i = 0; i < level-2; i++ {
			tap += "│ "
		}
		if i < level-1 {
			tap += "└─"
		}

		result = append(result, tap+strings.Split(node.content, "\n")[0])
		idList = append(idList, node.ID)
		for _, child := range node.child {
			DFS(child, level+1)
		}
	}

	DFS(root, 0)

	var returnResult string = ""

	result = result[1:]
	idList = idList[1:]

	for index, content := range result {
		returnResult += strconv.Itoa(index) + ": " + content + "\n"
	}

	return &pb.ReadResult{Message: returnResult, Data: idList}, nil
}

func (s *server) Choose(ctx context.Context, in *pb.ChooseMessage) (*pb.PostResult, error) {
	var id uint32 = in.GetNodeID()

	fmt.Println(id)
	return &pb.PostResult{Message: Nodes[id].content}, nil
}

func (s *server) Reply(ctx context.Context, in *pb.ReplyMessage) (*pb.PostResult, error) {
	var id uint32 = in.GetNodeID()
	var content string = in.GetMessage()

	log.Printf("[Reply] ID: %v, Content: %v", id, content)

	newReply(content, id)

	return &pb.PostResult{Message: "[Success] Reply ID: " + strconv.FormatUint(uint64(ID_counter), 10)}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	root = &Node{content: "", ID: 0, parent: nil, child: []*Node{}}
	Nodes = append(Nodes, root)

	s := grpc.NewServer()
	pb.RegisterBulletinServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
