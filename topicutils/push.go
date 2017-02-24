package topicutils

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
)

// PublishToTopic will publish to a given topic with 10 retry
func PublishToTopic(ctx context.Context, top *pubsub.Topic, messages []*pubsub.Message) error {

	for i := 0; i < 10; i++ {

		_, err := top.Publish(ctx, messages...)

		if err != nil {
			log.Println(err.Error())
			time.Sleep(1 * time.Second)
			continue

		} else {
			log.Printf("%s amount of message published", strconv.Itoa(len(messages)))
			for _, message := range messages {
				message.Done(true)
			}
			return nil
		}

	}

	return errors.New("publication failed even with 10 try")

}

// AckMessages simple iterate over the given message array and Ack all with the defined value
func AckMessages(actType bool, messages []*pubsub.Message) {
	for _, m := range messages {
		m.Done(actType)
	}
}
