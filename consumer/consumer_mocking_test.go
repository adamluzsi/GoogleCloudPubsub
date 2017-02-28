package consumer_test

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	. "github.com/LxDB/testing"
	"github.com/adamluzsi/GoogleCloudPubsub/consumer"
)

type ExampleHandler struct {
	messages []consumer.Message
	wg       *sync.WaitGroup
}

func (eh *ExampleHandler) HandleMessage(msg consumer.Message) error {
	eh.messages = append(eh.messages, msg)
	msg.Done(true)
	return nil
}
func (eh *ExampleHandler) Finish() error {
	defer eh.wg.Done()

	return nil
}

func TestConsumerMockingAllPerfect(t *testing.T) {
	consumer.TurnMockModOn()
	defer consumer.TurnMockModOff()

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	eh := &ExampleHandler{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	amount := 3
	cons := consumer.New(ctx, "example", fn,
		consumer.SetBatchSizeTo(amount),
		consumer.SetMaxExtensionDurationTo(10*time.Minute),
		consumer.SetWorkersCountTo(runtime.NumCPU()))

	cons.Start()
	defer cons.Stop()

	datas := make([][]byte, amount)
	for i := 0; i < amount; i++ {
		datas[i] = []byte(`Hello world ` + strconv.Itoa(i))
		consumer.MockMessageFeeder <- datas[i]
	}

	wg.Wait()

	if len(datas) != len(eh.messages) {
		t.Log(len(eh.messages))
		t.Fatal("expected message count is not equal to the received one")
	}

	for i, data := range datas {
		if !TestEqBytes(data, eh.messages[i].Data()) {
			t.Fatal("The two data and the passed value not matching!")
		}
	}

}
