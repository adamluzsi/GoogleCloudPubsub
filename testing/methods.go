package helpers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/adamluzsi/gcloud_pubsub/client"
	"google.golang.org/api/iterator"
)

const TopicName = "test-source-topic"

var Topic *pubsub.Topic

const SubscriptionName = "test-source-subscription"

var Subscription *pubsub.Subscription

func GetTimestamp() int {
	return int(time.Now().Unix())
}

type tester interface {
	Log(...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fail()
}

func SetUp(t tester) {

	ctx := context.Background()
	t.Log("begin to setup test environment")

	pubsubEmulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")
	if pubsubEmulatorHost == "" {
		log.Fatalln("No PUBSUB_EMULATOR_HOST variable set!")
	}

	t.Log("Create pubsub client for testing")
	client, err := client.New()

	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t.Log("create source topic")
	createTopic(ctx, t, client)

	t.Log("create source subscription")
	createSubscription(ctx, t, client)

}

func createTopic(ctx context.Context, t tester, client *pubsub.Client) {
	top := client.Topic(TopicName)

	ok, err := top.Exists(ctx)

	if err != nil {
		t.Fatalf("failed to check if topic exists: %v", err)
	}

	if !ok {
		if top, err = client.CreateTopic(ctx, TopicName); err != nil {
			t.Fatalf("failed to create the topic: %v", err)
		}
	}

	Topic = top

}

func createSubscription(ctx context.Context, t tester, client *pubsub.Client) {

	sub := client.Subscription(SubscriptionName)

	ok, err := sub.Exists(ctx)

	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}

	if ok {
		// cleanup
		if err := sub.Delete(ctx); err != nil {
			t.Fatalf("failed to cleanup the Subscription (%q): %v", SubscriptionName, err)
		}
	}

	t.Log(Topic)
	sub, err = client.CreateSubscription(context.Background(), SubscriptionName, Topic, 0, nil)

	ok, err = sub.Exists(ctx)

	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}

	if !ok {
		t.Fatalf("failed to create the Subscription (%q)", SubscriptionName)
	}

	Subscription = sub

}

func GetMessageDatas(quantity int, seed int) [][]byte {
	messageDatas := make([][]byte, 0, 1000)

	for index := 0; index < quantity; index++ {
		m := map[string]int{"seed": seed}
		data, _ := json.Marshal(m)
		messageDatas = append(messageDatas, data)
	}

	return messageDatas
}

func PublishExampleMessages(t tester, testAmount, GetTimestamp int) [][]byte {
	ctx := context.Background()
	messages := GetMessageDatas(testAmount, GetTimestamp)

	for _, message := range messages {
		_, err := Topic.Publish(ctx, &pubsub.Message{Data: message})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}

	return messages
}

func FetchFromSubscription(parent context.Context, t tester, amount int) []*pubsub.Message {
	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	it, err := Subscription.Pull(ctx)

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

