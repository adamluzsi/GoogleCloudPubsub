package publisher

import "context"

// Publisher is the public interface when you work with the publisher package
// This is used to force the developer writing code that can be tested easily
type Publisher interface {
	Publish(...[]byte) ([]string, error)
}

// New will return a new publisher object that can publish []byte messages to a given google cloud pubsub topic
func New(ctx context.Context, topicName string, opts ...Options) Publisher {
	conf := newConfig()
	conf.consumeOptions(opts)

	if useMock {
		m := &mock{ctx: ctx, topicName: topicName}
		return m
	}

	return &publisher{ctx: ctx, topicName: topicName}
}
