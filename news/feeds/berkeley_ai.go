package feeds

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/wurde/event-stream/stream"
	"github.com/wurde/event-stream/utils"
)

var BERKELEY_AI_URL = "https://bair.berkeley.edu/blog/feed.xml"

func MonitorBerkeleyAI(s *stream.StreamServer) {
	sched := gocron.NewScheduler(time.Local)

	_, es := s.Connect("Monitoring Berkeley AI News Feed")

	// Repeat ever ~6 hours
	sched.Every(utils.Random(375, 360)).Minutes().Do(func() {
		log.Println("Fetching latest Berkeley AI publications...")

		feed := getFeed(BERKELEY_AI_URL)

		lastEvent := es.GetLastEvent("news.feed.entry.berkeley.>")
		lastEventPublishedAt := utils.ParseTime(lastEvent)

		for i := len(feed.Items) - 1; i >= 0; i-- {
			item := feed.Items[i]

			timeLayout := "Mon, 02 Jan 2006 15:04:05 -0700"
			publishedAt, _ := time.Parse(timeLayout, item.Published)

			if publishedAt.Before(lastEventPublishedAt) || publishedAt.Equal(lastEventPublishedAt) {
				continue
			}

			es.PublishCloudEvent(&stream.CloudEventInput{
				Id:          item.GUID,
				Source:      "edu.berkeley.bair",
				EventType:   "news.feed.entry.berkeley.published.v20220416",
				PublishedAt: publishedAt,
				Data:        item,
			})
		}
	})

	sched.StartAsync()
}
