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
	sC, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		log.Fatal(err)
	}
	defer sC.Close()

	for i := 1; i <= 1000; i++ {
		if err = sC.Publish("orders", newOrder(i)); err != nil {
			log.Fatal(err)
		}

		log.Printf("order №%d published\n", i)
	}
}

func newOrder(num int) []byte {
	n := strconv.Itoa(num)
	return []byte("{" +
		"  \"order_uid\": \"order" + n + "\"," +
		"  \"track_number\": \"WBILMTESTTRACK\"," +
		"  \"entry\": \"WBIL\"," +
		"  \"delivery\": {" +
		"	 \"name\": \"Test Testov\"," +
		"    \"phone\": \"+9720000000\"," +
		"    \"zip\": \"2639809\"," +
		"    \"city\": \"Kiryat Mozkin\"," +
		"    \"address\": \"Ploshad Mira 15\"," +
		"    \"region\": \"Kraiot\"," +
		"    \"email\": \"test@gmail.com\"\n  }," +
		"  \"payment\": {" +
		"	 \"transaction\": \"b563feb7b2b84b6test\"," +
		"    \"request_id\": \"\"," +
		"    \"currency\": \"USD\"," +
		"    \"provider\": \"wbpay\"," +
		"    \"amount\": 1817," +
		"    \"payment_dt\": 1637907727," +
		"    \"bank\": \"alpha\"," +
		"    \"delivery_cost\": 1500," +
		"    \"goods_total\": 317," +
		"    \"custom_fee\": 0\n  }," +
		"  \"items\": [    {" +
		"      \"chrt_id\": 9934930," +
		"      \"track_number\": \"WBILMTESTTRACK\"," +
		"      \"price\": 453," +
		"      \"rid\": \"ab4219087a764ae0btest\"," +
		"      \"name\": \"Mascaras\"," +
		"      \"sale\": 30," +
		"      \"size\": \"0\"," +
		"      \"total_price\": 317," +
		"      \"nm_id\": 2389212," +
		"      \"brand\": \"Vivienne Sabo\"," +
		"     \"status\": 202\n    }  ]," +
		"  \"locale\": \"en\"," +
		"  \"internal_signature\": \"\"," +
		"  \"customer_id\": \"test\"," +
		"  \"delivery_service\": \"meest\"," +
		"  \"shardkey\": \"9\"," +
		"  \"sm_id\": 99," +
		"  \"date_created\": \"2021-11-26T06:22:19Z\"," +
		"  \"oof_shard\": \"1\"}")
}
