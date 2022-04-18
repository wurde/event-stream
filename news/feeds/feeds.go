package feeds

import (
	"context"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

func getFeed(url string) *gofeed.Feed {
	parser := gofeed.NewParser()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	feed, err := parser.ParseURLWithContext(url, ctx)
	if err != nil {
		log.Fatal(err)
	}
	return feed
}
