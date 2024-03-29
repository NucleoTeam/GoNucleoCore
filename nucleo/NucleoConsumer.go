package nucleohub

import (
	"encoding/json"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"time"
)

type NucleoConsumer struct {
	Brokers []string
	GroupID string
	Hub *NucleoHub
	Topics []string
}

func NewConsumer(group string, brokers []string, hub *NucleoHub) *NucleoConsumer {
	consumer := new(NucleoConsumer);
	consumer.Brokers = brokers
	consumer.GroupID = group
	consumer.Hub = hub
	return consumer;
}

func (cHandle * NucleoConsumer) AddTopic(topic string){
	cHandle.Topics = append(cHandle.Topics, topic)
}

func (cHandle * NucleoConsumer) Start(){
	go cHandle.readThread()
}

func (cHandle * NucleoConsumer) readThread(){
	run := true
	c, _ := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cHandle.Brokers[0],
		"group.id":          cHandle.GroupID,
		"auto.offset.reset": "latest",
		"session.timeout.ms":    6000,
	})
	c.SubscribeTopics(cHandle.Topics, nil)
	for run == true {
		data, _ := c.ReadMessage(1 * time.Millisecond)
		if data != nil {
			dataTmp := NewNucleoData()
			//fmt.Println(string(data.Value))
			errX := json.Unmarshal(data.Value, &dataTmp)
			if errX != nil {
				fmt.Println(errX)
			}
			cHandle.Hub.Execute(*data.TopicPartition.Topic, dataTmp)
		}
		time.Sleep(time.Millisecond)
	}
}