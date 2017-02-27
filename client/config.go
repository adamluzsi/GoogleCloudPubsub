package client

import (
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
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
