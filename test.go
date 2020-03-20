package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Sectorbob/mlab-ns2/gae/ns/digest"
	"github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func newClient(publicKey, privateKey string) (*mongodbatlas.Client, error) {
	transport := digest.NewTransport(publicKey, privateKey)

	client, err := transport.Client()
	if err != nil {
		return nil, err
	}

	return mongodbatlas.NewClient(client), nil
}

func main() {
	publicKey := os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
	privateKey := os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	if publicKey == "" || privateKey == "" || projectID == "" {
		log.Fatalln("MONGODB_ATLAS_PROJECT_ID, MONGODB_ATLAS_PUBLIC_KEY and MONGODB_ATLAS_PRIVATE_KEY must be set to run this example")
	}

	client, err := newClient(publicKey, privateKey)
	if err != nil {
		log.Fatalf(err.Error())
	}

	processes, _, err := client.Processes.List(context.Background(), projectID, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, process := range processes {
		fmt.Printf("%s:%d\n", process.Hostname, process.Port)

		measurements, _, err := client.Measurements.List(context.Background(), projectID, process.Hostname, process.Port, nil)
		if err != nil {
			log.Fatalf(err.Error())
		}

		for _, measurement := range measurements {
			fmt.Printf("  %s\n", measurement.Name)

			for _, datapoint := range measurement.Datapoints {
				fmt.Printf("    %s %f\n", datapoint.Timestamp, datapoint.Value)
			}
		}
	}
}
