package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
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
}

type server struct {
	pb.UnimplementedPrimaryServer
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

var ServerList []string

func (s *server) Post(ctx context.Context, in *pb.Content) (*pb.ID, error) {
	var content string = in.GetMessage()
	log.Debug("[Post] Received: ", content)

	ID_counter++

	// Update the data in servers
	for _, addr := range ServerList {
		{

			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error("did not connect: ", err)
			}
			defer conn.Close()
			c := pb.NewBulletinClient(conn)

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

	return &pb.ID{NodeID: ID_counter}, nil
}

func (s *server) Reply(ctx context.Context, in *pb.Node) (*pb.ID, error) {
	var parentID uint32 = in.GetParentID()
	var content string = in.GetMessage()

	log.Debug("[Reply] ID: %v,\nContent: %v", parentID, content)
	ID_counter++

	// Update the data in servers
	for _, addr := range ServerList {
		{
			log.Info("Update to ", addr)
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error("did not connect: ", err)
			}
			defer conn.Close()
			c := pb.NewBulletinClient(conn)

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

	return &pb.ID{NodeID: ID_counter}, nil
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPrimaryServer(s, &server{})
	log.Info("Primary server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: ", err)
	}
}
