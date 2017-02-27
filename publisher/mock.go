package publisher

import (
	"context"
	"math/rand"
	"strconv"
)

var useMock bool

// MockMessageReceiver is a channel that can take []byte messages that will be sent from the mock publisher
var MockMessageReceiver chan []byte

// TurnMockModOn will enable the consumer New method to return with Mock struct instead of the real one.
// This is only for testing purpose!
func TurnMockModOn() {
	MockMessageReceiver = make(chan []byte)
	useMock = true
}

// TurnMockModOff will disable the consumer New method to return with Mock struct instead of the real one.
// This is only for testing purpose!
func TurnMockModOff() {
	MockMessageReceiver = nil
	useMock = false
}

func init() {
	TurnMockModOff()
}

type mock struct {
	ctx       context.Context
	topicName string
}

func (m *mock) Publish(datas ...[]byte) ([]string, error) {

	go func() {
		for _, data := range datas {
			MockMessageReceiver <- data
		}
	}()

	resps := make([]string, 0, len(datas))
	for i := 0; i < len(datas); i++ {
		resps = append(resps, strconv.Itoa(rand.Int()))
	}

	return resps, nil
}
