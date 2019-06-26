package nucleohub

import (
	"encoding/json"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"time"
)

type NucleoProducer struct {
	Brokers []string
	Producer *kafka.Producer
	Hub *NucleoHub
	Queue *NucleoList
}
func NewProducer(brokers []string, hub *NucleoHub) *NucleoProducer {
	producer := new(NucleoProducer)
	producer.Brokers = brokers
	producer.Hub = hub
	producer.Queue = newList()

	p, _ := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers[0]})
	producer.Producer = p

	go producer.WriteThread();

	return producer;
}

func (p * NucleoProducer) WriteThread(){
	for{
		if p.Queue.Size>0 {
			item := p.Queue.Pop()
			if item != nil {
				dataJson, err := json.Marshal(item.Data)
				if err != nil {
					fmt.Println(err)
				}
				p.Producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &item.Chain, Partition: kafka.PartitionAny},
					Value: dataJson,
				}, nil)
			}
			p.Producer.Flush(1)
		}
		time.Sleep(time.Millisecond)
	}
}
