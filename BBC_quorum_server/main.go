package main

import (
	"context"
	"flag"
	"fmt"
	"net"
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

type server struct {
	pb.UnimplementedQuorumServer
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
	if (ParentID >= uint32(len(nodes))) || (ParentID < 0) {
		return
	}

	var newNode *Node = &Node{content: content, ID: NodeID, parent: nodes[ParentID], child: []*Node{}}
	nodes[ParentID].child = append(nodes[ParentID].child, newNode)
	nodes = append(nodes, newNode)
	fmt.Println(len(nodes[0].child))
}

func (s *server) Post(ctx context.Context, in *pb.Content) (*pb.ACK, error) {
	var content string = in.GetMessage()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: ", err)
	}
	defer conn.Close()
	c := pb.NewQuorumCoordinatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Post(ctx, &pb.Content{Message: content})

	if err != nil {
		log.Error("could not post to primary server: ", err)
		return &pb.ACK{Success: false}, nil
	}

	return &pb.ACK{Success: true}, nil
}

func (s *server) Read(ctx context.Context, in *pb.Empty) (*pb.ReadResult, error) {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: ", err)
	}
	defer conn.Close()
	c := pb.NewQuorumCoordinatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Read(ctx, &pb.Empty{})

	if err != nil {
		return nil, err
	}

	return &pb.ReadResult{Message: r.GetMessage(), Data: r.GetData()}, nil
}

func (s *server) Choose(ctx context.Context, in *pb.ID) (*pb.Content, error) {
	var id uint32 = in.GetNodeID()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: ", err)
	}
	defer conn.Close()
	c := pb.NewQuorumCoordinatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Choose(ctx, &pb.ID{NodeID: id})

	if err != nil {
		return nil, err
	}

	return &pb.Content{Message: r.GetMessage()}, nil
}

func (s *server) Reply(ctx context.Context, in *pb.Node) (*pb.ACK, error) {
	var parentID uint32 = in.GetNodeID()
	var content string = in.GetMessage()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: ", err)
	}
	defer conn.Close()
	c := pb.NewQuorumCoordinatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Reply(ctx, &pb.Node{ParentID: parentID, Message: content})

	if err != nil {
		log.Error("could not post to primary server: ", err)
		return &pb.ACK{Success: false}, err
	}

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

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	root = &Node{content: "", ID: 0, parent: nil, child: []*Node{}}
	nodes = append(nodes, root)

	s := grpc.NewServer()
	pb.RegisterQuorumServer(s, &server{})
	log.Info("Server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: ", err)
	}
}
