// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package consumer_test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"google.golang.org/api/iterator"

	. "github.com/LxDB/testing"
	"github.com/adamluzsi/gcloud_pubsub/consumer"

	"cloud.google.com/go/pubsub"
)

const sourceTopicName = "test-source-topic"

var sourceTopic *pubsub.Topic

const sourceSubscriptionName = "test-source-subscription"

var sourceSubscription *pubsub.Subscription

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

func TestConsuming(t *testing.T) {
	setup(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandlerAcker{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	publishExampleMessages(t, testAmount, timestamp())

	wg.Add(1)
	c := consumer.New(ctx, sourceSubscriptionName, fn, consumer.SetBatchSizeTo(testAmount))
	c.Start()
	wg.Wait()
	c.Stop()

	if len(eh.messages) != testAmount {
		t.Log(len(eh.messages))
		t.Fatal("expected message count is not equal to the received one")
	}

}

func TestConsumingDataPassTheMessageValue(t *testing.T) {
	setup(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandlerAcker{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	messages := publishExampleMessages(t, testAmount, timestamp())

	wg.Add(1)
	c := consumer.New(ctx, sourceSubscriptionName, fn, consumer.SetBatchSizeTo(testAmount))
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
	setup(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandlerNotAcker{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	publishExampleMessages(t, testAmount, timestamp())

	wg.Add(1)
	c := consumer.New(ctx, sourceSubscriptionName, fn, consumer.SetBatchSizeTo(testAmount))
	c.Start()
	wg.Wait()
	time.Sleep(1 * time.Second)
	c.Stop()

	messages := fetchFromSubscription(ctx, t, testAmount)

	if len(messages) != testAmount {
		t.Fatalf("expected message count is not equal to the received one: %v\n", len(messages))
	}

}

func TestConsumingDependOnHandlerForAckTheMessages(t *testing.T) {
	setup(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandlerAcker{messages: []consumer.Message{}, wg: &wg}
	fn := func() consumer.Handler { return eh }

	publishExampleMessages(t, testAmount, timestamp())

	wg.Add(1)
	c := consumer.New(ctx, sourceSubscriptionName, fn, consumer.SetBatchSizeTo(testAmount))
	c.Start()
	wg.Wait()
	time.Sleep(1 * time.Second)
	c.Stop()

	messages := fetchFromSubscription(ctx, t, 1)

	if len(messages) != 0 {
		t.Fatalf("expected message count is not equal to the received one: %v\n", len(messages))
	}

}

// -------------------------------------------HELPERS-------------------------------------------

func timestamp() int {
	return int(time.Now().Unix())
}

type tester interface {
	Log(...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fail()
}

func setup(t tester) {

	ctx := context.Background()
	t.Log("begin to setup test environment")

	pubsubEmulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")
	if pubsubEmulatorHost == "" {
		log.Fatalln("No PUBSUB_EMULATOR_HOST variable set!")
	}

	t.Log("Create pubsub client for testing")
	client, err := consumer.NewPubsubClient()

	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t.Log("create source topic")
	createSourceTopic(ctx, t, client)

	t.Log("create source subscription")
	createSourceSubscription(ctx, t, client)

}

func createSourceTopic(ctx context.Context, t tester, client *pubsub.Client) {
	top := client.Topic(sourceTopicName)

	ok, err := top.Exists(ctx)

	if err != nil {
		t.Fatalf("failed to check if topic exists: %v", err)
	}

	if !ok {
		if top, err = client.CreateTopic(ctx, sourceTopicName); err != nil {
			t.Fatalf("failed to create the topic: %v", err)
		}
	}

	sourceTopic = top

}

func createSourceSubscription(ctx context.Context, t tester, client *pubsub.Client) {

	sub := client.Subscription(sourceSubscriptionName)

	ok, err := sub.Exists(ctx)

	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}

	if ok {
		// cleanup
		if err := sub.Delete(ctx); err != nil {
			t.Fatalf("failed to cleanup the Subscription (%q): %v", sourceSubscriptionName, err)
		}
	}

	t.Log(sourceTopic)
	sub, err = client.CreateSubscription(context.Background(), sourceSubscriptionName, sourceTopic, 0, nil)

	ok, err = sub.Exists(ctx)

	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}

	if !ok {
		t.Fatalf("failed to create the Subscription (%q)", sourceSubscriptionName)
	}

	sourceSubscription = sub

}

func getMessageDatas(quantity int, seed int) [][]byte {
	messageDatas := make([][]byte, 0, 1000)

	for index := 0; index < quantity; index++ {
		m := map[string]int{"seed": seed}
		data, _ := json.Marshal(m)
		messageDatas = append(messageDatas, data)
	}

	return messageDatas
}

func publishExampleMessages(t tester, testAmount, timestamp int) [][]byte {
	ctx := context.Background()
	messages := getMessageDatas(testAmount, timestamp)

	for _, message := range messages {
		_, err := sourceTopic.Publish(ctx, &pubsub.Message{Data: message})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}

	return messages
}

func fetchFromSubscription(parent context.Context, t tester, amount int) []*pubsub.Message {
	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	it, err := sourceSubscription.Pull(ctx)

	if err != nil {
		t.Fatal("Message iterator failed to create!")
	}

	messages := make([]*pubsub.Message, 0, amount)

	for index := 0; index < amount; index++ {

		msg, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			t.Log(err)
			break
		}

		messages = append(messages, msg)
		msg.Done(true)

	}

	it.Stop()
	return messages
}
