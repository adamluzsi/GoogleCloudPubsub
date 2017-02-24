package consumer

import (
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

// Start Begin the Google Pubsub Consuming
func (c *Consumer) Start() {
	c.wg.Add(1)
	go c.subscriptionWorker()
}

func (c *Consumer) subscriptionWorker() {
	defer c.wg.Done()

	var subscriptionWaitGroup sync.WaitGroup

initLoop:
	for {
		select {
		case <-c.ctx.Done():
			break initLoop

		default:

			pub, err := NewPubsubClient()

			if err != nil {
				log.Println("pullWorker: pubsub client creation failed!")
				time.Sleep(5 * time.Second)
				continue initLoop
			}

			sub := pub.Subscription(c.subscriptionName)

			for index := 0; index < c.workersCount; index++ {
				subscriptionWaitGroup.Add(1)
				go c.subscribe(sub, &subscriptionWaitGroup)
			}

			subscriptionWaitGroup.Wait()

		}
	}

	subscriptionWaitGroup.Wait()
	log.Println("subscriptionWorker done")
}

func (c *Consumer) subscribe(sub *pubsub.Subscription, w *sync.WaitGroup) {
	defer w.Done()

messageProcessingLoop:
	for {
		select {
		case <-c.ctx.Done():
			break messageProcessingLoop

		default:
			err := c.processMessages(sub)

			if err != nil {
				log.Println(err)
				break messageProcessingLoop
			}
		}
	}
}

func (c *Consumer) processMessages(sub *pubsub.Subscription) error {
	handler := c.handlerConstructor()

	it, err := sub.Pull(c.ctx, pubsub.MaxExtension(c.maxExtension), pubsub.MaxPrefetch(c.batchSize))

	if err != nil {
		return err
	}

	defer it.Stop()

	messages := make([]*pubsub.Message, 0, c.batchSize)

	for index := 0; index < c.batchSize; index++ {

		msg, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return err
		}

		messages = append(messages, msg)
		err = handler.HandleMessage(newMessageWrapper(msg))

		if err != nil {
			msg.Done(false)
		}

	}

	err = handler.Finish()

	if err != nil {
		for _, m := range messages {
			m.Done(false)
		}
	}

	return nil
}
