package publisher

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type publisher struct {
	ctx       context.Context
	topic     *pubsub.Topic
	topicName string
}

// Publish will publish a []byte message to a given topic
func (p *publisher) Publish(datas ...[]byte) ([]string, error) {
	var ids []string
	var err error

	err = p.initialSetup()
	if err != nil {
		return nil, err
	}

	_, err = p.topic.Exists(p.ctx)
	if err != nil {
		err = p.resetTopicObject()
		if err != nil {
			return ids, err
		}
	}

	ids, err = p.publishByteMessages(datas)

	return ids, err
}

func (p *publisher) publishByteMessages(datas [][]byte) ([]string, error) {
	topic := p.topic

	chunks := chunkSlice(datas)

	ids := make([]string, 0)
	for _, chunk := range chunks {
		newIds, err := p.publishChunk(topic, chunk)
		if err != nil {
			return ids, err
		}

		ids = append(ids, newIds...)
	}

	return ids, nil
}

func (p *publisher) publishChunk(top *pubsub.Topic, chunk [][]byte) ([]string, error) {

	messages := make([]*pubsub.Message, 0, len(chunk))
	for _, data := range chunk {
		messages = append(messages, &pubsub.Message{Data: data})
	}

	resps, err := top.Publish(p.ctx, messages...)

	if err != nil {
		return nil, err
	}

	return resps, nil

}
