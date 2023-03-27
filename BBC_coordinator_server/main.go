package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/ARui-tw/I2DS_Bulletin-Board-Consistency/BBC"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Type    string `json:"type"`
	Primary int    `json:"primary"`
	Child   []int  `json:"child"`
	NR      int    `json:"NR"`
	NW      int    `json:"NW"`
}

type server struct {
	pb.UnimplementedQuorumCoordinatorServer
}

type Node struct {
	content string
	ID      uint32
	parent  *Node
	child   []*Node
}

var (
	port        = flag.Int("port", 50051, "The server port")
	config_file = flag.String("config_file", "server_config.json", "The server config file")
)

var ID_counter uint32 = 0
var (
	NR         int
	NW         int
	N          int
	root       *Node
	nodes      []*Node = []*Node{}
	ServerList []string
)

func newPost(content string) {
	var newNode *Node = &Node{content: content, ID: ID_counter, parent: root, child: []*Node{}}
	root.child = append(root.child, newNode)
	nodes = append(nodes, newNode)
}

func newReply(content string, parentID uint32) {
	var parentNode *Node = nodes[parentID]
	var newNode *Node = &Node{content: content, ID: ID_counter, parent: parentNode, child: []*Node{}}
	parentNode.child = append(parentNode.child, newNode)
	nodes = append(nodes, newNode)
}

func (s *server) Post(ctx context.Context, in *pb.Content) (*pb.ACK, error) {
	var content string = in.GetMessage()
	log.Debug("[Post] Received: ", content)

	ID_counter++

	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(N)

	// Update data in NW random servers
	for _, idx := range p[:NW] {
		{
			conn, err := grpc.Dial(ServerList[idx], grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error("did not connect: ", err)
			}
			defer conn.Close()
			c := pb.NewQuorumClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := c.Update(ctx, &pb.Node{Message: content, NodeID: ID_counter, ParentID: 0})

			if err != nil || !r.GetSuccess() {
				log.Error("could not post to child server: ", err)
				ID_counter--
				return nil, err
			}
		}
	}

	newPost(content)

	return &pb.ACK{Success: true}, nil
}

func (s *server) Read(ctx context.Context, in *pb.Empty) (*pb.ReadResult, error) {
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

	result = result[1:]
	idList = idList[1:]

	for index, content := range result {
		result[index] = strconv.Itoa(index) + ": " + content
	}

	return &pb.ReadResult{Message: result, Data: idList}, nil
}

func (s *server) Choose(ctx context.Context, in *pb.ID) (*pb.Content, error) {
	var id uint32 = in.GetNodeID()

	return &pb.Content{Message: nodes[id].content}, nil
}

func (s *server) Reply(ctx context.Context, in *pb.Node) (*pb.ACK, error) {
	var parentID uint32 = in.GetParentID()
	var content string = in.GetMessage()

	log.Debug("[Reply] ID: %v,\nContent: %v", parentID, content)
	ID_counter++

	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(N)

	// Update data in NW random servers
	for _, idx := range p[:NW] {
		{
			conn, err := grpc.Dial(ServerList[idx], grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error("did not connect: ", err)
			}
			defer conn.Close()
			c := pb.NewQuorumClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := c.Update(ctx, &pb.Node{Message: content, NodeID: ID_counter, ParentID: parentID})

			if err != nil || !r.GetSuccess() {
				log.Error("could not post to child server: ", err)
				ID_counter--
				return nil, err
			}
		}

	}

	newReply(content, parentID)

	return &pb.ACK{Success: true}, nil
}

func (s *server) Synch(ctx context.Context, in *pb.IDs) (*pb.Nodes, error) {

	var node pb.Node = pb.Node{Message: "test", NodeID: 1, ParentID: 0}
	var nodes []*pb.Node = []*pb.Node{&node}

	return &pb.Nodes{Nodes: nodes}, nil
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

	for _, port := range config.Child {
		ServerList = append(ServerList, fmt.Sprintf("localhost:%d", port))
	}

	NR = config.NR
	NW = config.NW
	N = len(ServerList)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	root = &Node{content: "", ID: 0, parent: nil, child: []*Node{}}
	nodes = append(nodes, root)

	s := grpc.NewServer()
	pb.RegisterQuorumCoordinatorServer(s, &server{})
	log.Info("Primary server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: ", err)
	}
}
