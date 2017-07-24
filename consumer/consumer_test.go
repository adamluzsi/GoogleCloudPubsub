package consumer_test

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	. "github.com/LxDB/testing"
	"github.com/adamluzsi/GoogleCloudPubsub/consumer"
	. "github.com/adamluzsi/GoogleCloudPubsub/testing"
)

type ExampleHandlerAcker struct {
	messages []consumer.Message
	wg       *sync.WaitGroup
}

func (eh *ExampleHandlerAcker) HandleMessage(msg consumer.Message) error {
	eh.messages = append(eh.messages, msg)
	msg.Done(true)
	return nil
}
func (eh *ExampleHandlerAcker) Finish() error {
	defer eh.wg.Done()

	return nil
}

type ExampleHandlerNotAcker struct {
	messages []consumer.Message
	wg       *sync.WaitGroup
}

func (eh *ExampleHandlerNotAcker) HandleMessage(msg consumer.Message) error {
	eh.messages = append(eh.messages, msg)

	return nil
}
func (eh *ExampleHandlerNotAcker) Finish() error {
	defer eh.wg.Done()

	return nil
}

func ExampleConsumer(ctx context.Context, SubscriptionName string, handlerConstructor consumer.HandlerConstructor) consumer.Consumer {
	return consumer.New(ctx, SubscriptionName, handlerConstructor)
}

func ExampleConsumer_withMaxExtensionConfigured(ctx context.Context, SubscriptionName string, handlerConstructor consumer.HandlerConstructor) consumer.Consumer {
	return consumer.New(ctx, SubscriptionName, handlerConstructor, consumer.SetMaxExtensionDurationTo(10*time.Minute))
}

func ExampleConsumer_withWorkersCountConfigured(ctx context.Context, SubscriptionName string, handlerConstructor consumer.HandlerConstructor) consumer.Consumer {
	return consumer.New(ctx, SubscriptionName, handlerConstructor, consumer.SetWorkersCountTo(runtime.NumCPU()))
}

func ExampleConsumer_withBatchAmount(ctx context.Context, SubscriptionName string, handlerConstructor consumer.HandlerConstructor, testAmount int) consumer.Consumer {
	return consumer.New(ctx, SubscriptionName, handlerConstructor, consumer.SetBatchSizeTo(testAmount))
}

func TestConsuming(t *testing.T) {
	SetUp(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandlerAcker{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	PublishExampleMessages(t, testAmount, GetTimestamp())

	wg.Add(1)
	c := ExampleConsumer_withBatchAmount(ctx, SubscriptionName, fn, testAmount)
	c.Start()
	wg.Wait()
	c.Stop()

	if len(eh.messages) != testAmount {
		t.Log(len(eh.messages))
		t.Fatal("expected message count is not equal to the received one")
	}

}

func TestConsumingDataPassTheMessageValue(t *testing.T) {
	SetUp(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandlerAcker{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	messages := PublishExampleMessages(t, testAmount, GetTimestamp())

	wg.Add(1)
	c := ExampleConsumer_withBatchAmount(ctx, SubscriptionName, fn, testAmount)
	c.Start()
	wg.Wait()
	c.Stop()

	// bad pattern to depend on array order create testing method for that
	for i, originalMessage := range messages {
		if !TestEqBytes(originalMessage, eh.messages[i].Data()) {
			t.Log(len(eh.messages))
			t.Fatal("expected message count is not equal to the received one")
		}
	}

}

func TestConsumingNotAckedMessagesWillReturnToSubscription(t *testing.T) {
	SetUp(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandlerNotAcker{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	PublishExampleMessages(t, testAmount, GetTimestamp())

	wg.Add(1)
	c := ExampleConsumer_withBatchAmount(ctx, SubscriptionName, fn, testAmount)
	c.Start()
	wg.Wait()
	time.Sleep(1 * time.Second)
	c.Stop()

	messages := FetchFromSubscription(ctx, t, testAmount)

	if len(messages) != testAmount {
		t.Fatalf("expected message count is not equal to the received one: %v\n", len(messages))
	}

}

func TestConsumingDependOnHandlerForAckTheMessages(t *testing.T) {
	SetUp(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandlerAcker{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	PublishExampleMessages(t, testAmount, GetTimestamp())

	wg.Add(1)
	c := ExampleConsumer_withBatchAmount(ctx, SubscriptionName, fn, testAmount)
	c.Start()
	wg.Wait()
	time.Sleep(1 * time.Second)
	c.Stop()

	messages := FetchFromSubscription(ctx, t, 1)

	if len(messages) != 0 {
		t.Fatalf("expected message count is not equal to the received one: %v\n", len(messages))
	}

}
