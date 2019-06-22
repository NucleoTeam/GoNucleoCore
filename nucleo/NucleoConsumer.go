package nucleohub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type NucleoConsumer struct {
	Brokers []string
	Chain string
	GroupID string
	Reader *kafka.Reader
	Hub *NucleoHub
}

func NewConsumer(chain string, group string, brokers []string, hub *NucleoHub) *NucleoConsumer {
	fmt.Println("Registered chain: " + chain)
	consumer := new(NucleoConsumer);
	consumer.Brokers = brokers
	consumer.GroupID = group
	consumer.Hub = hub
	consumer.Chain = chain
	c, _ := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	go consumer.readThread()

	return consumer;
}
func (c * NucleoConsumer) readThread(){
		for {
			m, err := c.Reader.FetchMessage(context.Background())
			if err != nil {
				fmt.Println(err.Error())
			}
			data := NewNucleoData()
			err = json.Unmarshal([]byte(m.Value), &data)
			if err != nil {
				fmt.Println(err)
			}
			if data.Objects == nil {
				data.Objects = map[string]interface{}{}
			}
			c.Hub.Execute(m.Topic, data)
		}
}