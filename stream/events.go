package stream

import (
	"encoding/json"
	"log"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats.go"
)

type CloudEventInput struct {
	Id          string
	Source      string
	EventType   string
	PublishedAt time.Time
	Data        interface{}
}

func CreateCloudEvent(e *CloudEventInput) cloudevents.Event {
	ce := cloudevents.NewEvent()

	// The version of the CloudEvents specification used.
	ce.SetSpecVersion("1.0")
	// Identifies the event. The `id` attribute is unique
	// across all events related to one `source`.
	ce.SetID(e.Id)
	// Identifies the context in which an event happened.
	ce.SetSource(e.Source)
	// Describes the type of event related to the originating
	// occurrence. Often this attribute is used for routing,
	// observability, policy enforcement. When a CloudEvent's
	// data changes in a backwardly-compatible way, the value
	// of the type attribute should generally stay the same.
	// Otherwise it should change. A version scheme should
	// include a a date (vYYYYMMDD) as part of the value.
	ce.SetType(e.EventType)
	// Timestamp of when the occurrence happened.
	ce.SetTime(e.PublishedAt)
	// All data not related to routing an event.
	ce.SetData(cloudevents.ApplicationJSON, e.Data)

	return ce
}

func (es *EventStream) PublishCloudEvent(input *CloudEventInput) {
	event := CreateCloudEvent(input)

	bytes, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}

	_, err = es.Publish(
		event.Type(),
		bytes,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (es *EventStream) GetLastEvent(subject string) *nats.Msg {
	sub, _ := es.SubscribeSync(
		subject,
		nats.DeliverLast(),
	)
	defer sub.Unsubscribe()

	msg, _ := sub.NextMsg(1 * time.Second)
	return msg
}

func (es *EventStream) GetAllNewsEvents(subject string) []*nats.Msg {
	var newsEvents []*nats.Msg

	sub, err := es.PullSubscribe(subject, "new")
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	msgs, _ := sub.Fetch(100)
	for len(msgs) > 0 {
		for i := 0; i < len(msgs); i++ {
			newsEvents = append(newsEvents, msgs[i])
			msgs[i].AckSync()
		}
		msgs, _ = sub.Fetch(100)
	}

	return newsEvents
}
