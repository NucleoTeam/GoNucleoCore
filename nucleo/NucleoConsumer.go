package nucleohub

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type NucleoConsumer struct {
	Brokers []string
	Chain string
	GroupID string
	Consumer *kafka.Consumer
	Hub *NucleoHub
}

func NewConsumer(chain string, group string, brokers []string, hub *NucleoHub) *NucleoConsumer {
	consumer := new(NucleoConsumer);
	consumer.Brokers = brokers
	consumer.GroupID = group
	consumer.Hub = hub
	consumer.Chain = chain
	c, _ := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
		"group.id":          consumer.GroupID,
		"auto.offset.reset": "earliest",
	})
	consumer.Consumer = c
	consumer.Consumer.Subscribe(chain, nil)

	go consumer.readThread()

	return consumer;
}
func (c * NucleoConsumer) readThread(){
	for {
		m, err := c.Consumer.ReadMessage(-1)
		if err != nil {
			fmt.Println(err)
			break
		}
		data := NewNucleoData()
		err = json.Unmarshal(m.Value, &data)
		if err != nil {
			fmt.Println(err)
		}
		c.Hub.Execute(m.TopicPartition.String(), data)
	}
}