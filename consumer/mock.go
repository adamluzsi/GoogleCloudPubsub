package consumer

import (
	"context"
	"log"
	"sync"
	"time"
)

var useMock bool

// MockMessageFeeder is a channel that can take []byte messages that will be feeded to the mock consumers
var MockMessageFeeder chan []byte

// TurnMockModOn will enable the consumer New method to return with Mock struct instead of the real one.
// This is only for testing purpose!
func TurnMockModOn() {
	log.Println("NOTICE: Mock mod enabled for gcloud pubsub consumer")
	MockMessageFeeder = make(chan []byte)
	useMock = true
}

// TurnMockModOff will disable the consumer New method to return with Mock struct instead of the real one.
// This is only for testing purpose!
func TurnMockModOff() {
	MockMessageFeeder = nil
	useMock = false
}

func init() {
	TurnMockModOff()
}

type mockMessage struct {
	data         []byte
	wg           *sync.WaitGroup
	ackedAlready bool
}

func newMockMessage(data []byte, wg *sync.WaitGroup) Message {
	wg.Add(1)
	return &mockMessage{data: data, wg: wg}
}

func (mm *mockMessage) Data() []byte {
	return mm.data
}

func (mm *mockMessage) Done(_ bool) {
	if !mm.ackedAlready {
		mm.ackedAlready = true
		mm.wg.Done()
	}
}

type mock struct {
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
func (m *mock) Start() {
	m.wg.Add(1)
	go m.work()
}

// Stop will stop the mock consumer processing
func (m *mock) Stop() {
	m.cancel()
	m.wg.Wait()
}

func (m *mock) work() {
	defer m.wg.Done()

working:
	for {
		select {
		case <-m.ctx.Done():
			break working

		default:

			var wg sync.WaitGroup
			handler := m.handlerConstructor()
			for i := 0; i < m.batchSize; i++ {
				select {
				case <-m.ctx.Done():
					break working

				case data, ok := <-MockMessageFeeder:

					if !ok {
						break working
					}

					msg := newMockMessage(data, &wg)
					err := handler.HandleMessage(msg)

					if err != nil {
						log.Println(err)
						msg.Done(false) // like it matter or anything...
					}

				}
			}

			err := handler.Finish()
			if err != nil {
				// nack all msg
				continue working
			}

			wg.Wait()

		}
	}

}
