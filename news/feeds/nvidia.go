package feeds

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/wurde/event-stream/stream"
	"github.com/wurde/event-stream/utils"
)

var NVIDIA_URL = "https://developer.nvidia.com/blog/feed"

func MonitorNvidia(s *stream.StreamServer) {
	sched := gocron.NewScheduler(time.Local)

	_, es := s.Connect("Monitoring Nvidia News Feed")

	// Repeat ever ~6 hours
	sched.Every(utils.Random(375, 360)).Minutes().Do(func() {
		log.Println("Fetching latest Nvidia publications...")

		feed := getFeed(NVIDIA_URL)

		lastEvent := es.GetLastEvent("news.feed.entry.nvidia.>")
		lastEventPublishedAt := utils.ParseTime(lastEvent)

		for i := len(feed.Items) - 1; i >= 0; i-- {
			item := feed.Items[i]

			publishedAt, _ := time.Parse(time.RFC3339, item.Published)

			if publishedAt.Before(lastEventPublishedAt) || publishedAt.Equal(lastEventPublishedAt) {
				continue
			}

			es.PublishCloudEvent(&stream.CloudEventInput{
				Id:          item.GUID,
				Source:      "com.nvidia",
				EventType:   "news.feed.entry.nvidia.published.v20220416",
				PublishedAt: publishedAt,
				Data:        item,
			})
		}
	})

	sched.StartAsync()
}
