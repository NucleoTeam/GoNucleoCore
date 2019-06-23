package nucleohub

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type NucleoProducer struct {
	Brokers []string
	Producer *kafka.Producer
	Chain string
	Hub *NucleoHub
	Queue *NucleoList
}
func NewProducer(chain string, brokers []string, hub *NucleoHub) *NucleoProducer {
	producer := new(NucleoProducer)
	producer.Brokers = brokers
	producer.Hub = hub
	producer.Chain = chain
	producer.Queue = newList()

	p, _ := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers[0]})
	producer.Producer = p

	go producer.WriteThread();

	return producer;
}

func (p * NucleoProducer) WriteThread(){
	go func() {
		for e := range p.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	for{
		if p.Queue.Size>0 {
			item := p.Queue.Pop()
			if item != nil {
				dataJson, err := json.Marshal(item.Data)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(item.Chain)
				fmt.Println(string(dataJson))
				p.Producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &item.Chain, Partition: kafka.PartitionAny},
					Value: dataJson,
				}, nil)
			}
			p.Producer.Flush(1)
		}
		time.Sleep(1 * time.Microsecond)
	}
}
