package main

import (
	"flag"
	"github.com/nats-io/stan.go"
	"log"
)

func parseArgs() (sc stan.Conn, ch string, msg string) {
	crid := flag.String("crid", "test-cluster", "NS Cluster ID")
	clid := flag.String("clid", "me", "NS Client ID")
	chz := flag.String("ch", "orders", "NS Channel name")
	message := flag.String("msg", "Default message", "Message to send")
	flag.Parse()
	sc, _ = stan.Connect(*crid, *clid, stan.NatsURL(stan.DefaultNatsURL))
	return sc, *chz, *message
}

func publish(sc stan.Conn, ch string, msg string) {
	if err := sc.Publish(ch, []byte(msg)); err != nil {
		log.Fatal(err)
	}
}

func main() {
	sc, ch, msg := parseArgs()
	publish(sc, ch, msg)
}
