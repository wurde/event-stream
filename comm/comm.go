package comm

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/go-co-op/gocron"
	"github.com/mmcdole/gofeed"
	"github.com/wurde/event-stream/stream"
)

func ScheduleDailyEmailDigest(s *stream.StreamServer) {
	sched := gocron.NewScheduler(time.Local)

	_, es := s.Connect("Scheduling Daily Digest Email")

	sched.Every(1).Day().At("08:00").Do(func() {
		FROM_EMAIL := os.Getenv("FROM_EMAIL")
		TO_EMAIL := os.Getenv("TO_EMAIL")
		if FROM_EMAIL == "" || TO_EMAIL == "" {
			log.Fatal("AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and TO_EMAIL must be set")
		}

		entries := es.GetAllNewsEvents("news.feed.entry.>")
		if len(entries) == 0 {
			log.Println("No new entries found")
			return
		}
		log.Printf("Sending daily digest (%v)...", len(entries))

		var cloudEvent event.Event
		var eventData gofeed.Item
		var bodyText string

		for _, entry := range entries {
			json.Unmarshal(entry.Data, &cloudEvent)
			json.Unmarshal(cloudEvent.Data(), &eventData)

			bodyText += eventData.Title + "\n" + eventData.Link + "\n\n"
		}

		log.Println("")
		log.Println("")
		log.Println(bodyText)

		e := newEmail()
		e.From = FROM_EMAIL
		e.To = []string{TO_EMAIL}
		e.Subject = "My Daily Digest"
		e.Text = bodyText

		email, err := e.send()
		if err != nil {
			log.Fatal(err)
		}

		es.PublishCloudEvent(&stream.CloudEventInput{
			Id:          *email.MessageId,
			Source:      "aws.ses." + AWS_REGION,
			EventType:   "comm.email.sent.v20220416",
			PublishedAt: time.Now(),
			Data:        email,
		})
	})

	sched.StartAsync()
}
