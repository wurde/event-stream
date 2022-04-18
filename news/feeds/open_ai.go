package feeds

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/wurde/event-stream/stream"
	"github.com/wurde/event-stream/utils"
)

var NOT_SO_OPEN_AI_URL = "https://openai.com/blog/rss"

func MonitorNotSoOpenAI(s *stream.StreamServer) {
	sched := gocron.NewScheduler(time.Local)

	_, es := s.Connect("Monitoring NotSoOpenAI News Feed")

	// Repeat ever ~6 hours
	sched.Every(utils.Random(375, 360)).Minutes().Do(func() {
		log.Println("Fetching latest NotSoOpenAI publications...")

		feed := getFeed(NOT_SO_OPEN_AI_URL)

		lastEvent := es.GetLastEvent("news.feed.entry.openai.>")
		lastEventPublishedAt := utils.ParseTime(lastEvent)

		for i := len(feed.Items) - 1; i >= 0; i-- {
			item := feed.Items[i]

			timeLayout := "Mon, 02 Jan 2006 15:04:05 MST"
			publishedAt, _ := time.Parse(timeLayout, item.Published)

			if publishedAt.Before(lastEventPublishedAt) || publishedAt.Equal(lastEventPublishedAt) {
				continue
			}

			es.PublishCloudEvent(&stream.CloudEventInput{
				Id:          item.GUID,
				Source:      "com.openai",
				EventType:   "news.feed.entry.openai.published.v20220416",
				PublishedAt: publishedAt,
				Data:        item,
			})
		}
	})

	sched.StartAsync()
}
