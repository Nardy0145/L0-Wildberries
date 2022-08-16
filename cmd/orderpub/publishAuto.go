package main

import (
	"123/models"
	"encoding/json"
	"github.com/bxcodec/faker"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	sc, err := stan.Connect("test-cluster", "me", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
	}
	o := new(models.Order)
	err = faker.FakeData(o)
	if err != nil {
		log.Printf("can't create fake data: %s", err.Error())
	}

	b, err := json.Marshal(o)
	if err != nil {
		log.Printf("error marshaling message %s", err.Error())
	}
	err = sc.Publish("orders", b)
	if err != nil {
		log.Fatal(err)
	}
}
