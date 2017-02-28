package publisher

import (
	"context"
	"errors"

	"github.com/adamluzsi/GoogleCloudPubsub/client"

	"cloud.google.com/go/pubsub"
)

func (p *publisher) initialSetup() error {
	var err error
	if p.topic == nil {
		err = p.resetTopicObject()
	}
	return err
}

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

func (p *publisher) resetTopicObject() error {
	topic, err := p.getTopic()

	if err != nil {
		return err
	}

	p.topic = topic
	return nil
}
