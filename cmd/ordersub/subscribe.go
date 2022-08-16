package main

import (
	"123/cache"
	"123/database"
	"123/models"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
)

func parseArgs() (sc stan.Conn, ch string) {
	crid := flag.String("crid", "test-cluster", "NS Cluster ID")
	clid := flag.String("clid", "subscriber", "NS Client ID")
	chz := flag.String("ch", "orders", "NS Channel name")
	flag.Parse()
	sc, err := stan.Connect(*crid, *clid, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
		return nil, ""
	}
	return sc, *chz
}

func subscribe(sc stan.Conn, ch string, db *database.DbConnect) {
	_, _ = sc.Subscribe(ch, func(m *stan.Msg) {
		fmt.Printf(string(m.Data))
		fmt.Printf("New message. Channel: %s\n", m.Subject)
		err := cache.AppendCache(m.Data, db)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
	})
	fmt.Println("Subscribed to channel: ", ch)
}

func renewCache(db *database.DbConnect) {
	rows, _ := db.Db.Query("select data from wborders")
	for rows.Next() {
		var x []byte
		err := rows.Scan(&x)
		if err != nil {
			log.Fatal(err)
		}
		var i models.Order
		json.Unmarshal(x, &i)
		data, _ := json.Marshal(i)
		cache.UpdateCache(i, data)
	}
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	sc, ch := parseArgs()
	fmt.Println("Loading cache...")
	renewCache(db)
	fmt.Println("Cache has been loaded")
	subscribe(sc, ch, db)
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
