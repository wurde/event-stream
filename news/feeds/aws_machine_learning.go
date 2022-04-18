package feeds

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/wurde/event-stream/stream"
	"github.com/wurde/event-stream/utils"
)

var AWS_ML_URL = "https://aws.amazon.com/blogs/machine-learning/feed"

func MonitorAwsMachineLearning(s *stream.StreamServer) {
	sched := gocron.NewScheduler(time.Local)

	_, es := s.Connect("Monitoring AWS Machine Learning News Feed")

	// Repeat ever ~6 hours
	sched.Every(utils.Random(375, 360)).Minutes().Do(func() {
		log.Println("Fetching latest AWS ML publications...")

		feed := getFeed(AWS_ML_URL)

		lastEvent := es.GetLastEvent("news.feed.entry.aws.>")
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
				Source:      "com.amazon.aws",
				EventType:   "news.feed.entry.aws.published.v20220416",
				PublishedAt: publishedAt,
				Data:        item,
			})
		}
	})

	sched.StartAsync()
}
