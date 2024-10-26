package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"sync"
	"time"

	pb "github.com/Amr2812/data-services/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sendRequest(client pb.MessagesServiceClient, channelId int64, wg *sync.WaitGroup) {
	defer wg.Done()

	message := &pb.MessageRequest{
		ChannelId: channelId,
		MessageId: channelId*1000 + int64(rand.Intn(4)+1),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.GetMessage(ctx, message)
	if err != nil {
		log.Printf("Failed to get message for channel %d: %v", channelId, err)
		return
	}
}

func main() {
	requests := flag.Int("requests", 1000, "Total number of requests to send")
	numChannels := flag.Int("channels", 3, "Number of unique channels to distribute requests across")
	flag.Parse()

	dataServices := []string{"localhost:50051"}
	var clients []pb.MessagesServiceClient

	for _, ds := range dataServices {
		conn, err := grpc.NewClient(ds, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("could not connect to %s: %v", ds, err)
		}
		defer conn.Close()
		clients = append(clients, pb.NewMessagesServiceClient(conn))
	}

	var wg sync.WaitGroup
	wg.Add(*requests)

	start := time.Now()

	for i := 0; i < *requests; i++ {
		channelId := int64((i % *numChannels) + 1)         // Add 1 to start channels from 1
		client := clients[int(channelId-1)%len(clients)]   // Distribute requests across data services (-1 for 0-based index)
		go sendRequest(client, channelId, &wg)
	}

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Completed %d requests for %d channels in %v", *requests, *numChannels, elapsed)
}
