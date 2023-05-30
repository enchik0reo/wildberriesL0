package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/enchik0reo/wildberriesL0/internal/models"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/stan.go"
)

const (
	url       = "nats://localhost:4222"
	clusterID = "test-cluster"
	clientID  = "order-producer"
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
	gf := gofakeit.New(0)
	model := models.Basic{}

	model.OrderUid = "order" + strconv.Itoa(num)
	model.TrackNumber = "WBILMTESTTRACK"
	model.Entry = "WBIL"
	model.Delivery.Name = gf.Name()
	model.Delivery.Phone = gf.Phone()
	model.Delivery.Zip = strconv.Itoa(gf.Number(0, 9999999))

	adr := gf.Address()

	model.Delivery.City = adr.City
	model.Delivery.Address = adr.Street
	model.Delivery.Region = adr.State
	model.Delivery.Email = gf.Email()

	for i := 0; i < gf.Number(1, 10); i++ {
		item := models.Items{}

		item.ChrtId = gf.Number(0, 9999999)
		item.TrackNumber = "WBILMTESTTRACK"
		item.Price = gf.Number(1, 10000)
		item.Rid = gf.UUID()
		item.Name = "toy car"
		item.Sale = gf.Number(0, 99)
		item.Size = "0"
		item.TotalPrice = item.Price - ((item.Price * item.Sale) / 100)
		item.NmId = gf.Number(0, 9999999)
		item.Brand = gf.Car().Brand
		item.Status = 202

		model.Items = append(model.Items, item)
	}

	model.Payment.Transaction = gf.UUID()
	model.Payment.RequestId = ""
	model.Payment.Currency = gf.CurrencyShort()
	model.Payment.Provider = "wbpay"
	model.Payment.PaymentDt = gf.Number(0, 9999999999)
	model.Payment.Bank = "alpha"
	model.Payment.DeliveryCost = gf.Number(0, 5000)

	model.Locale = adr.State
	model.InternalSignature = ""
	model.CustomerId = "test"
	model.DeliveryService = "meets"
	model.Shardkey = "9"
	model.SmId = gf.Number(0, 1000)
	model.DateCreated = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.UTC)
	model.OofShard = "1"

	sumPrice := func(c []models.Items) int {
		s := 0
		for _, i := range c {
			s += i.TotalPrice
		}
		return s
	}

	model.Payment.Amount = sumPrice(model.Items) + model.Payment.DeliveryCost
	model.Payment.GoodsTotal = gf.Number(0, 1000)
	model.Payment.CustomFee = 0

	bytes, err := json.MarshalIndent(model, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
