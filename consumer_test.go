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

	"github.com/adamluzsi/gcloud_pubsub_pipeline/consumer"

	"cloud.google.com/go/pubsub"
)

const sourceTopicName = "test-source-topic"

var sourceTopic *pubsub.Topic

const sourceSubscriptionName = "test-source-subscription"

var sourceSubscription *pubsub.Subscription

type ExampleHandler struct {
	messages []consumer.MessageInterface
	wg       *sync.WaitGroup
}

func (eh *ExampleHandler) HandleMessage(msg consumer.MessageInterface) error {
	eh.messages = append(eh.messages, msg)
	return nil
}
func (eh *ExampleHandler) Finish() error {
	defer eh.wg.Done()
	// do something important
	return nil
}

func init() {
	os.Setenv("PUBSUB_KEYFILE_JSON", `{"project_id":"testing"}`)
}

func TestConsuming(t *testing.T) {
	setup(t)

	ctx := context.Background()

	testAmount := 10

	var wg sync.WaitGroup
	eh := &ExampleHandler{messages: []consumer.MessageInterface{}, wg: &wg}

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

// func BenchmarkConsuming(b *testing.B) {
// 	setup(b)
// 	ctx := context.Background()

// 	testAmount := 100
// 	testTimes := 100
// 	tstamp := timestamp()

// 	for index := 0; index < testTimes; index++ {
// 		messages := getMessageDatas(testAmount, tstamp+index)

// 		for _, message := range messages {
// 			_, err := sourceTopic.Publish(ctx, &pubsub.Message{Data: message})
// 			if err != nil {
// 				b.Log(err)
// 				b.Fail()
// 			}
// 		}
// 	}

// 	consumer := psub.NewConsumer(ctx, c)

// 	b.ResetTimer()
// 	consumer.Start(5, testAmount)
// 	fetchFromTargetSubscription(ctx, b)
// 	consumer.Stop()

// }

func publishExampleMessages(t tester, testAmount, timestamp int) {
	ctx := context.Background()
	messages := getMessageDatas(testAmount, timestamp)

	for _, message := range messages {
		_, err := sourceTopic.Publish(ctx, &pubsub.Message{Data: message})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}

}

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
