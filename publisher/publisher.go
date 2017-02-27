package publisher

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type publisher struct {
	ctx       context.Context
	topicName string
}

// Publish will publish a []byte message to a given topic
func (p *publisher) Publish(datas ...[]byte) ([]string, error) {

	topic, err := p.getTopic()

	if err != nil {
		return nil, err
	}

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

func chunkSlice(slices [][]byte) [][][]byte {
	chunks := make([][][]byte, 0)
	chunk := make([][]byte, 0, pubsub.MaxPublishBatchSize)

	for _, slice := range slices {
		if len(chunk) != pubsub.MaxPublishBatchSize {
			chunk = append(chunk, slice)
		} else {
			chunks = append(chunks, chunk)
			chunk = make([][]byte, 0, pubsub.MaxPublishBatchSize)
		}
	}

	if len(chunk) != 0 {
		chunks = append(chunks, chunk)
	}

	return chunks
}

func (p *publisher) publishChunk(top *pubsub.Topic, chunk [][]byte) ([]string, error) {

	messages := make([]*pubsub.Message, 0, pubsub.MaxPublishBatchSize)
	for _, data := range chunk {
		messages = append(messages, &pubsub.Message{Data: data})
	}

	resps, err := top.Publish(p.ctx, messages...)

	if err != nil {
		return nil, err
	}

	return resps, nil

}
