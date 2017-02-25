package consumer

import (
	"context"
	"sync"
	"time"
)

// Consumer is a struct that can pull messages from a subscription,
// and pass values to a handler that can be constucted with a given handlerConstructor
//
// Handler must implement the HandlerConstructor interface
// and the return value must implement the Handler interface
type Consumer interface {
	Start()
	Stop()
}

type consumer struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	subscriptionName   string
	handlerConstructor HandlerConstructor

	workersCount int
	batchSize    int
	maxExtension time.Duration
}

// Start Begin the Google Pubsub Consuming
func (c *consumer) Start() {
	c.wg.Add(1)
	go c.subscriptionWorker()
}

// Stop will gracefully cancel the subscription consuming
// usefull when comined with exit signal handling
func (c *consumer) Stop() {
	c.cancel()
	c.wg.Wait()
}
