package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	pb "github.com/ARui-tw/I2DS_Bulletin-Board-Consistency/BBC"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 50052, "The server port")
	addr = flag.String("addr", "localhost:50051", "Primary Server Address")
)

// var addr string
// var PrimaryServerPort int

type server struct {
	pb.UnimplementedBulletinServer
}

type Node struct {
	content string
	ID      uint32
	parent  *Node
	child   []*Node
}

var root *Node

var nodes []*Node = []*Node{}

func newNode(content string, NodeID uint32, ParentID uint32) {
	var newNode *Node = &Node{content: content, ID: NodeID, parent: nodes[ParentID], child: []*Node{}}
	nodes[ParentID].child = append(nodes[ParentID].child, newNode)
	nodes = append(nodes, newNode)
	fmt.Println(len(nodes[0].child))

	if nodes[0] != root {
		fmt.Println("root")
	}
}

// func newPost(content string, ID uint32) {
// 	var newNode *Node = &Node{content: content, ID: ID, parent: root, child: []*Node{}}
// 	root.child = append(root.child, newNode)
// 	nodes = append(nodes, newNode)
// }

// func newReply(content string, parentID uint32) {
// var parentNode *Node = nodes[parentID]
// var newNode *Node = &Node{content: content, ID: ID_counter, parent: parentNode, child: []*Node{}}
// parentNode.child = append(parentNode.child, newNode)
// nodes = append(nodes, newNode)
// }

func (s *server) Post(ctx context.Context, in *pb.Content) (*pb.ACK, error) {
	var content string = in.GetMessage()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: ", err)
	}
	defer conn.Close()
	c := pb.NewPrimaryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Post(ctx, &pb.Content{Message: content})

	if err != nil {
		log.Error("could not post to primary server: ", err)
		return &pb.ACK{Success: false}, nil
	}

	// ID := r.GetNodeID()

	// newNode(content, ID, 0)

	return &pb.ACK{Success: true}, nil
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

func (s *server) Choose(ctx context.Context, in *pb.ID) (*pb.Content, error) {
	var id uint32 = in.GetNodeID()

	fmt.Println(id)
	return &pb.Content{Message: nodes[id].content}, nil
}

func (s *server) Reply(ctx context.Context, in *pb.Node) (*pb.ACK, error) {
	var parentID uint32 = in.GetNodeID()
	var content string = in.GetMessage()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: ", err)
	}
	defer conn.Close()
	c := pb.NewPrimaryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Reply(ctx, &pb.Node{ParentID: parentID, Message: content})

	if err != nil {
		log.Error("could not post to primary server: ", err)
		return &pb.ACK{Success: false}, nil
	}

	ID := r.GetNodeID()

	log.Debug("[Reply] ID: %v, Content: %v", ID, content)

	// newNode(content, ID, parentID)

	return &pb.ACK{Success: true}, nil
}

// Get the update from primary server, send back ACK when done
func (s *server) Update(ctx context.Context, in *pb.Node) (*pb.ACK, error) {
	NodeID := in.GetNodeID()
	content := in.GetMessage()
	ParentID := in.GetParentID()

	newNode(content, NodeID, ParentID)

	return &pb.ACK{Success: true}, nil
}

// Start the server
//   - @param port: the port of the server
//   - @param config: the config file of all the other servers
// func StartUPServer(port int, primaryServerPort int) {
// 	PrimaryServerPort = primaryServerPort
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	addr = fmt.Sprintf("localhost:%d", PrimaryServerPort)

// 	root = &Node{content: "", ID: 0, parent: nil, child: []*Node{}}
// 	nodes = append(nodes, root)

// 	s := grpc.NewServer()
// 	pb.RegisterBulletinServer(s, &server{})
// 	log.Info("Server listening at ", lis.Addr())
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: ", err)
// 	}
// }

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	root = &Node{content: "", ID: 0, parent: nil, child: []*Node{}}
	nodes = append(nodes, root)

	s := grpc.NewServer()
	pb.RegisterBulletinServer(s, &server{})
	log.Info("Server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: ", err)
	}
}
