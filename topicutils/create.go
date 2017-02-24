package topicutils

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// CreateTopicIfNotExists Create Google Cloud Pubsub topic if not Exists
func CreateTopicIfNotExists(pub *pubsub.Client, topicName string) (*pubsub.Topic, error) {

	ctx := context.Background()

	t := pub.Topic(topicName)
	ok, err := t.Exists(ctx)

	if err != nil {
		return nil, err
	}

	if ok {
		return t, nil
	}

	t, err = pub.CreateTopic(ctx, topicName)

	if err != nil {
		return nil, err
	}

	return t, nil

}
