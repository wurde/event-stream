package news

import (
	"github.com/wurde/event-stream/news/feeds"
	"github.com/wurde/event-stream/stream"
)

func StartMonitoringFeeds(s *stream.StreamServer) {
	feeds.MonitorAwsMachineLearning(s)
	feeds.MonitorBerkeleyAI(s)
	feeds.MonitorFAIR(s)
	feeds.MonitorGoogleAI(s)
	feeds.MonitorMitMachineLearning(s)
	feeds.MonitorNvidia(s)
	feeds.MonitorNotSoOpenAI(s)
	feeds.MonitorSec10KFilings(s)
}
