package main

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/Amr2812/data-services/messages"
	"github.com/gocql/gocql"
	"google.golang.org/grpc"
)

var session *gocql.Session
var requestsMap *RequestsMap

type RequestId struct {
	channelId int64 // not go channel, but a channel for chat messages
	messageId int64
}

// map of request ids to array of channels
type RequestsMap struct {
	mu sync.Mutex
	// channel returns pointer to avoid a lot of copying and it won't be modified by the caller
	requests map[RequestId][]chan *pb.MessageReply
}

func NewRequestsMap() *RequestsMap {
	return &RequestsMap{
		requests: make(map[RequestId][]chan *pb.MessageReply),
	}
}

func (rm *RequestsMap) HandleRequest(requestId RequestId) *pb.MessageReply {
	resultChan := make(chan *pb.MessageReply)
	rm.mu.Lock()

	channels, ok := rm.requests[requestId]
	if !ok {
		channels = make([]chan *pb.MessageReply, 0)
		channels = append(channels, resultChan)
		rm.requests[requestId] = channels
		rm.mu.Unlock()

		go rm.ExecuteQuery(requestId)
	} else {
		channels = append(channels, resultChan)
		rm.requests[requestId] = channels
		rm.mu.Unlock()
	}

	return <-resultChan
}

// execute the cassandra query
// send the result to all the channels in the map
// remove the request from the map
func (rm *RequestsMap) ExecuteQuery(requestId RequestId) {
	var MessageReply pb.MessageReply

	err := session.Query("SELECT * FROM messages WHERE channel_id = ? AND message_id = ?", requestId.channelId, requestId.messageId).
		Scan(&MessageReply.ChannelId, &MessageReply.MessageId, &MessageReply.AuthorId, &MessageReply.Content)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
	}

	rm.mu.Lock()
	channels := rm.requests[requestId]
	delete(rm.requests, requestId)
	rm.mu.Unlock()

	log.Printf("Sending result to %d channels", len(channels))
	for _, ch := range channels {
		ch <- &MessageReply
	}
}

type MessageService struct {
	pb.UnimplementedMessagesServiceServer
}

func (s *MessageService) GetMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageReply, error) {
	requestId := RequestId{channelId: in.ChannelId, messageId: in.MessageId}
	return requestsMap.HandleRequest(requestId), nil
}

func main() {
	var err error
	cluster := gocql.NewCluster("cassandra-node1")
	cluster.Keyspace = "dataservices"
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("failed to connect to cassandra: %v", err)
	}
	log.Printf("Connected to cassandra")
	defer session.Close()

	requestsMap = NewRequestsMap()

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMessagesServiceServer(s, &MessageService{})
	log.Printf("Starting server on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		panic(err)
	}
}
