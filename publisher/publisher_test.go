package publisher_test

import (
	"context"
	"testing"

    "github.com/adamluzsi/gcloud_pubsub/publisher"

    . "github.com/adamluzsi/gcloud_pubsub/testing"
)

func TestPublishing(t *testing.T) {
	SetUp(t)

	ctx := context.Background()

    p := publisher.New(ctx, TopicName)
	p.Publish([]byte(`hello world!`))
	p.Publish([]byte(`hello`), []byte(`world`))

	messages := FetchFromSubscription(ctx, t, 3)

	if len(messages) != 3 {
		t.Fatalf("expected message count is not equal to the received one (%v)\n", len(messages))
	}

}

func TestMockedPublishing(t *testing.T) {
	publisher.TurnMockModOn()
	defer publisher.TurnMockModOff()

	ctx := context.Background()

    p := publisher.New(ctx, TopicName)
	p.Publish([]byte(`hello world!`))
	p.Publish([]byte(`hello`), []byte(`world`))

	datas := make([][]byte, 0, 3)
	for i := 0; i < 3; i++ {
		data := <- publisher.MockMessageReceiver
		datas = append(datas, data)
	}

	if len(datas) != 3 {
		t.Fatalf("expected message count is not equal to the received one (%v)\n", len(datas))
	}

}