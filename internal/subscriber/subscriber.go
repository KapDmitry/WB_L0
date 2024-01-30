package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/KapDmitry/WB_L0/internal/cache"
	"github.com/KapDmitry/WB_L0/internal/config"
	"github.com/KapDmitry/WB_L0/internal/logger"
	"github.com/KapDmitry/WB_L0/internal/order"
	"github.com/KapDmitry/WB_L0/internal/repo"
	"github.com/KapDmitry/WB_L0/internal/validator"
	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	Config config.NATSConfig
	CTX    context.Context
	DB     repo.Repo
	Cache  cache.Cache
	Log    logger.Logger
}

func (s *Subscriber) Listen() {
	sConn, err := stan.Connect(s.Config.ClusterID, "client1", stan.NatsURL(s.Config.URL))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer sConn.Close()

	sConn.Subscribe(s.Config.Subject, func(msg *stan.Msg) {
		var ord order.Order
		err = json.Unmarshal(msg.Data, &ord)
		if err != nil {
			s.Log.Log("Info", err.Error())
			return
		}

		if !validator.Validate(ord, s.Log) {
			msg.Ack()
			return
		}

		err = s.DB.Add(s.CTX, ord)
		if err != nil {
			s.Log.Log("Error", err.Error())
			return
		}
		s.Log.LogW("Info", "succesfully inserted in DB: ", map[string]interface{}{"orderID": ord.OrderUID})
		msg.Ack()
		err = s.Cache.Update(ord)
		if err != nil {
			s.Log.Log("Info", err.Error())
			return
		}

	}, stan.DurableName("ultradurable"), stan.SetManualAckMode(), stan.AckWait(10*time.Second), stan.MaxInflight(10))

	for {
	}
}
