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
type Consumer struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	subscriptionName   string
	handlerConstructor HandlerConstructor

	workersCount int
	batchSize    int
	maxExtension time.Duration
}

// New will create a new consumer structure that can handle work
func New(parent context.Context, subscriptionName string, handlerConstructor HandlerConstructor, opts ...Options) *Consumer {

	ctx, cancel := context.WithCancel(parent)
	var wg sync.WaitGroup

	validateInputString(subscriptionName)

	conf := newConfig()
	conf.consumeOptions(opts)

	return &Consumer{
		subscriptionName:   subscriptionName,
		handlerConstructor: handlerConstructor,
		maxExtension:       conf.maxExtension,
		workersCount:       conf.workersCount,
		batchSize:          conf.batchSize,
		ctx:                ctx,
		cancel:             cancel,
		wg:                 &wg}
}

// Stop will gracefully cancel the subscription consuming
// usefull when comined with exit signal handling
func (c *Consumer) Stop() {
	c.cancel()
	c.wg.Wait()
}
