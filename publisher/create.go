package publisher

import (
	"context"
	"errors"

	"github.com/adamluzsi/GoogleCloudPubsub/client"

	"cloud.google.com/go/pubsub"
)

// getTopic get the required topic from Google Cloud Pubsub topics
func (p *publisher) getTopic() (*pubsub.Topic, error) {

	ctx := context.Background()

	pubsub, e := client.New()

	if e != nil {
		return nil, e
	}

	t := pubsub.Topic(p.topicName)

	ok, err := t.Exists(ctx)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("topic not found")
	}

	return t, nil

}
