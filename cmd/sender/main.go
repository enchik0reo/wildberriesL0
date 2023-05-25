package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/enchik0reo/wildberriesL0/internal/models"

	"github.com/nats-io/stan.go"
)

const (
	url       = "nats://localhost:4222"
	clusterID = "test-cluster"
	clientID  = "order-sender"
)

func main() {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	file, err := os.Open("internal/models/model.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	model := models.Basic{}

	if err = json.Unmarshal(bytes, &model); err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 100; i++ {
		if err = sc.Publish("orders", newOrder(i, &model)); err != nil {
			log.Fatal(err)
		}

		log.Printf("order №%d published\n", i)
	}
}

func newOrder(num int, model *models.Basic) []byte {
	model.OrderUid = "order" + strconv.Itoa(num)

	bytes, err := json.MarshalIndent(model, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
