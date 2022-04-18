package stream

import (
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type StreamServer struct {
	*server.Server
}

type EventStream struct {
	nats.JetStreamContext
}

func startServer() *StreamServer {
	ns, err := server.NewServer(&server.Options{
		JetStream:          true,
		JetStreamMaxStore:  1024 * 1024 * 512, // 512MB
		JetStreamMaxMemory: 1024 * 1024 * 512, // 512MB
		StoreDir:           "./store",
	})
	if err != nil {
		log.Fatal(err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		log.Fatal("not ready for connection after 4 seconds")
	}

	stream := &StreamServer{ns}

	return stream
}

func (s *StreamServer) addStream(name string, subjects []string) {
	_, js := s.Connect("Add stream " + name)

	// https://pkg.go.dev/github.com/nats-io/nats.go#StreamConfig
	js.AddStream(&nats.StreamConfig{
		// A name for the Stream. May NOT have spaces, tabs, period (.),
		// greater than (>) or asterisk (*) characters.
		Name: name,
		// The type of storage backend. Either nats.FileStorage (default)
		// or nats.MemoryStorage.
		Storage: nats.FileStorage,
		// A list of subjects to consume, supports wildcards.
		Subjects: subjects,
		// When a Stream reaches it's limits either, DiscardNew refuses
		// new messages while DiscardOld (default) deletes old messages
		Discard: nats.DiscardOld,
		// How message retention is considered, LimitsPolicy (default),
		// InterestPolicy or WorkQueuePolicy.
		Retention: nats.LimitsPolicy,
		// How many bytes the Stream may contain.
		MaxBytes: 1024 * 1024 * 128, // 128MB
	})
}

// It is recommended to create the JetStream consumer
// (using js.AddConsumer) so that the JS consumer is not
// deleted on an Unsubscribe() or Drain() when the member
// that created the consumer goes away first.
func (s *StreamServer) addConsumer(stream string, config *nats.ConsumerConfig) {
	_, js := s.Connect("Add consumer " + stream)

	js.AddConsumer(stream, config)
}

func (stream *StreamServer) Connect(name string) (*nats.Conn, *EventStream) {
	nc, err := nats.Connect(
		stream.ClientURL(),
		nats.Name(name),
	)
	if err != nil {
		log.Fatal(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	es := &EventStream{js}

	return nc, es
}

func Setup() *StreamServer {
	stream := startServer()

	stream.addStream("news", []string{"news.feed.entry.>"})
	stream.addStream("comm", []string{"comm.email.>"})
	stream.addConsumer("news", &nats.ConsumerConfig{
		Durable:       "new",
		Description:   "Deliver the news",
		DeliverPolicy: nats.DeliverAllPolicy,
		ReplayPolicy:  nats.ReplayInstantPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
	})

	return stream
}
