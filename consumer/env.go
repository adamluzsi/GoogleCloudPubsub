package consumer

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
