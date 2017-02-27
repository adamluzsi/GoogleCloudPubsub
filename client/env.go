package client

import (
	"log"
	"os"
)

func pubsubKeyfileJSON() []byte {
	keyfileJSONStr := os.Getenv("PUBSUB_KEYFILE_JSON")

	if keyfileJSONStr == "" {
		log.Fatalln("missing env variable: PUBSUB_KEYFILE_JSON")
	}

	return []byte(keyfileJSONStr)
}

func isAnEmulatedEnvironment() bool {
	return os.Getenv("PUBSUB_EMULATOR_HOST") != ""
}
