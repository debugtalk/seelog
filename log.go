package seelog

import (
	"log"
)

func printInfo(msg string) {
	log.Printf("[seelog] Info: %s\n", msg)
}

func printError(err error) {
	log.Printf("[seelog] Error: %+v\n", err)
}
