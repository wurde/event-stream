package utils

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/nats-io/nats.go"
)

func Random(max, min int) int {
	return rand.Intn(max-min) + min
}

func ParseTime(lastEvent *nats.Msg) time.Time {
	var lastEventData event.Event
	var lastEventPublishedAt time.Time

	if lastEvent != nil {
		_ = json.Unmarshal(lastEvent.Data, &lastEventData)
		lastEventPublishedAt = lastEventData.Time()
	} else {
		lastEventPublishedAt, _ = time.Parse("2006-01-02", "2000-01-01")
	}
	return lastEventPublishedAt
}
