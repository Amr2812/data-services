package main

import (
	"context"
	"flag"
	"log"
	"strconv"
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
		MessageId: channelId*1000,
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
	requests := flag.Int("requests", 10000, "Total number of requests to send")
	numChannels := flag.Int("channels", 20, "Number of unique channels to distribute requests across")
	flag.Parse()
	if (*requests <= 0) || (*numChannels <= 0) {
		log.Fatalf("Invalid requests or channels, must be greater than 0")
	}
	if *numChannels > 20 {
		log.Fatalf("Number of channels must be less than or equal to 20 because there are only 20 unique channels in cassandra")
	}

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

	var metrics []*pb.MetricsReply
	for _, client := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		metricsReply, err := client.GetAndResetMetrics(ctx, &pb.Empty{})
		if err != nil {
			log.Fatalf("Failed to get metrics: %v", err)
		}
		metrics = append(metrics, metricsReply)
	}

	totalRequests := 0
	queriesExecuted := 0
	for _, m := range metrics {
		totalRequests += int(m.TotalRequests)
		queriesExecuted += int(m.QueriesExecuted)
	}
	avgQueries := strconv.FormatFloat(float64(queriesExecuted)/float64(totalRequests), 'f', -1, 64)

	log.Printf("Total requests: %d, Total queries executed: %d", totalRequests, queriesExecuted)
	log.Printf("Average queries per request: %s", avgQueries)
	log.Printf("Saved queries by coalescing: %d", totalRequests-queriesExecuted)
	log.Printf("Total time taken: %v", elapsed)
}
