package consumer

import "log"

func validateInputString(name string) {
	if name == "" {
		log.Fatalln("Invalid parameter given")
	}
}
