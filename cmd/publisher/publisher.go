package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/KapDmitry/WB_L0/internal/config"
	"github.com/KapDmitry/WB_L0/utils"
	"github.com/nats-io/stan.go"
)

func main() {
	stanCfg := config.NATSConfig{}
	stanCfg.Load("../../config/nats/")
	pConn, err := stan.Connect(stanCfg.ClusterID, "publisher", stan.NatsURL(stanCfg.URL))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer pConn.Close()

	for {
		orderData := utils.GenerateRandomOrder()
		byr, err := json.Marshal(orderData)
		if err != nil {
			fmt.Println(err.Error())
		}
		pConn.Publish(stanCfg.Subject, byr)
		fmt.Printf("order with ID: %s published!\n", orderData.OrderUID)
		time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
	}

}
