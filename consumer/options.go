package consumer

import (
	"time"

	"cloud.google.com/go/pubsub"
)

// Options is the public interface when working with options
type Options interface {
	Configure(config *config)
}

type config struct {
	maxExtension time.Duration
	workersCount int
	batchSize    int
}

const (
	// DefaultWorkerCount used to determine the default worker count on start
	DefaultWorkerCount = 1

	// DefaultBatchSize used to determine the default message batch and prefetch count
	DefaultBatchSize = 500
)

func newConfig() *config {
	return &config{
		workersCount: DefaultWorkerCount,
		batchSize:    DefaultBatchSize,
		maxExtension: pubsub.DefaultMaxExtension}
}

func (c *config) consumeOptions(opts []Options) {
	for _, opt := range opts {
		opt.Configure(c)
	}
}

type workersCountSetter struct {
	workersCount int
}

func (wcs *workersCountSetter) Configure(config *config) {
	config.workersCount = wcs.workersCount
}

// SetWorkersCountTo will return an options object that can configure the workers count for a Consumer
func SetWorkersCountTo(count int) Options {
	return &workersCountSetter{workersCount: count}
}

type batchSizeSetter struct {
	batchSize int
}

func (wcs *batchSizeSetter) Configure(config *config) {
	config.batchSize = wcs.batchSize
}

// SetBatchSizeTo will return an options object that can configure the prefetch and message bulk size for the Consumer
func SetBatchSizeTo(count int) Options {
	return &batchSizeSetter{batchSize: count}
}

type maxExtensionSetter struct {
	maxExtension time.Duration
}

func (mes *maxExtensionSetter) Configure(config *config) {
	config.maxExtension = mes.maxExtension
}

// SetMaxExtensionDurationTo will return an option that
// can set the maxExtension duration for google pubsub message iterator
func SetMaxExtensionDurationTo(maxExtension time.Duration) Options {
	return &maxExtensionSetter{maxExtension: maxExtension}
}
