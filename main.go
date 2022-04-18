package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/wurde/event-stream/comm"
	"github.com/wurde/event-stream/news"
	"github.com/wurde/event-stream/stream"
)

func main() {
	log.Println("Running event stream example...")

	godotenv.Load()
	s := stream.Setup()
	news.StartMonitoringFeeds(s)
	comm.ScheduleDailyEmailDigest(s)
	s.WaitForShutdown()
}
