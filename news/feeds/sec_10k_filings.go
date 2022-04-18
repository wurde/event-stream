package feeds

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/wurde/event-stream/stream"
	"github.com/wurde/event-stream/utils"
)

var SEC_10K_URL = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&type=10-k&start=0&count=100&output=atom"

func MonitorSec10KFilings(s *stream.StreamServer) {
	sched := gocron.NewScheduler(time.Local)

	_, es := s.Connect("Monitoring SEC 10-K Filings News Feed")

	// Repeat ever ~6 hours
	sched.Every(utils.Random(375, 360)).Minutes().Do(func() {
		log.Println("Fetching latest SEC 10-K filings...")

		feed := getFeed(SEC_10K_URL)

		lastEvent := es.GetLastEvent("news.feed.entry.sec.>")
		lastEventPublishedAt := utils.ParseTime(lastEvent)

		for i := len(feed.Items) - 1; i >= 0; i-- {
			item := feed.Items[i]

			publishedAt, _ := time.Parse(time.RFC3339, item.Published)

			if publishedAt.Before(lastEventPublishedAt) || publishedAt.Equal(lastEventPublishedAt) {
				continue
			}

			es.PublishCloudEvent(&stream.CloudEventInput{
				Id:          item.GUID,
				Source:      "gov.sec",
				EventType:   "news.feed.entry.sec.published.v20220416",
				PublishedAt: publishedAt,
				Data:        item,
			})
		}
	})

	sched.StartAsync()
}
