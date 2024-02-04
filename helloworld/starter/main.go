package main

import (
	"github.com/temporalio/samples-go/helloworld"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	result, err := helloworld.Run(c, "hello-world")
	if err != nil {
		log.Fatalln("Workflow run failed with err:", err)
	}
	log.Println("Workflow result:", result)
}
