package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

func projectName() string {
	config := make(map[string]interface{})

	json.Unmarshal(pubsubKeyfileJSON(), &config)

	v, ok := config["project_id"]
	if !ok {
		log.Fatalln("missing project name from pubsubKeyfileJSON!")
	}

	return v.(string)
}

func getJWTConfigFromJSON() *jwt.Config {
	conf, err := google.JWTConfigFromJSON(pubsubKeyfileJSON(), pubsub.ScopePubSub)

	if err != nil {
		log.Fatal(err)
	}

	return conf
}

// NewPubsubClient return a google pubsub client struct using the PUBSUB_KEYFILE_JSON env configuration
func NewPubsubClient() (*pubsub.Client, error) {

	ctx := context.Background()

	if os.Getenv("PUBSUB_EMULATOR_HOST") != "" {

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
