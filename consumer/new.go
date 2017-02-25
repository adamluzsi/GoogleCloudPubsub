package consumer

import (
	"context"
	"sync"
)

// New will create a new consumer structure that can handle work
func New(parent context.Context, subscriptionName string, handlerConstructor HandlerConstructor, opts ...Options) Consumer {

	ctx, cancel := context.WithCancel(parent)
	var wg sync.WaitGroup

	validateInputString(subscriptionName)

	conf := newConfig()
	conf.consumeOptions(opts)

	if useMock {

		return &mock{
			subscriptionName:   subscriptionName,
			handlerConstructor: handlerConstructor,
			maxExtension:       conf.maxExtension,
			workersCount:       conf.workersCount,
			batchSize:          conf.batchSize,
			ctx:                ctx,
			cancel:             cancel,
			wg:                 &wg}

	}

	return &consumer{
		subscriptionName:   subscriptionName,
		handlerConstructor: handlerConstructor,
		maxExtension:       conf.maxExtension,
		workersCount:       conf.workersCount,
		batchSize:          conf.batchSize,
		ctx:                ctx,
		cancel:             cancel,
		wg:                 &wg}

}
