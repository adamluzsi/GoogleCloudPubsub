package client

import (
	"context"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// New return a google pubsub client struct using the PUBSUB_KEYFILE_JSON env configuration
func New() (*pubsub.Client, error) {

	ctx := context.Background()

	if isAnEmulatedEnvironment() {

		client, err := pubsub.NewClient(ctx, "")

		if err != nil {
			return nil, err
		}

		return client, nil

	}

	ts := getJWTConfigFromJSON().TokenSource(ctx)
	client, err := pubsub.NewClient(ctx, projectName(), option.WithTokenSource(ts))

	if err != nil {
		return nil, err
	}

	return client, nil

}
