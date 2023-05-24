package main

import (
	"log"
	"strconv"

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

	for i := 1; i <= 100; i++ {
		if err = sc.Publish("orders", newOrder(i)); err != nil {
			log.Fatal(err)
		}

		log.Printf("order №%d published\n", i)
	}
}

func newOrder(num int) []byte {
	n := strconv.Itoa(num)
	return []byte("{\n" +
		"  \"order_uid\": \"order" + n + "\",\n" +
		"  \"track_number\": \"WBILMTESTTRACK\",\n" +
		"  \"entry\": \"WBIL\",\n" +
		"  \"delivery\": {\n" +
		"    \"name\": \"Test Testov\",\n" +
		"    \"phone\": \"+9720000000\",\n" +
		"    \"zip\": \"2639809\",\n" +
		"    \"city\": \"Kiryat Mozkin\",\n" +
		"    \"address\": \"Ploshad Mira 15\",\n" +
		"    \"region\": \"Kraiot\",\n" +
		"    \"email\": \"test@gmail.com\"\n  },\n" +
		"  \"payment\": {\n" +
		"    \"transaction\": \"b563feb7b2b84b6test\",\n" +
		"    \"request_id\": \"\",\n" +
		"    \"currency\": \"USD\",\n" +
		"    \"provider\": \"wbpay\",\n" +
		"    \"amount\": 1817,\n" +
		"    \"payment_dt\": 1637907727,\n" +
		"    \"bank\": \"alpha\",\n" +
		"    \"delivery_cost\": 1500,\n" +
		"    \"goods_total\": 317,\n" +
		"    \"custom_fee\": 0\n  },\n" +
		"  \"items\": [\n    {\n" +
		"      \"chrt_id\": 9934930,\n" +
		"      \"track_number\": \"WBILMTESTTRACK\",\n" +
		"      \"price\": 453,\n" +
		"      \"rid\": \"ab4219087a764ae0btest\",\n" +
		"      \"name\": \"Mascaras\",\n" +
		"      \"sale\": 30,\n" +
		"      \"size\": \"0\",\n" +
		"      \"total_price\": 317,\n" +
		"      \"nm_id\": 2389212,\n" +
		"      \"brand\": \"Vivienne Sabo\",\n" +
		"     \"status\": 202\n    }\n  ],\n" +
		"  \"locale\": \"en\",\n" +
		"  \"internal_signature\": \"\",\n" +
		"  \"customer_id\": \"test\",\n" +
		"  \"delivery_service\": \"meest\",\n" +
		"  \"shardkey\": \"9\",\n" +
		"  \"sm_id\": 99,\n" +
		"  \"date_created\": \"2021-11-26T06:22:19Z\",\n" +
		"  \"oof_shard\": \"1\"\n}")
}
