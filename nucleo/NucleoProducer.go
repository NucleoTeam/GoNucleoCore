package nucleohub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type NucleoProducer struct {
	Brokers []string
	Writer *kafka.Writer
	Chain string
	Hub *NucleoHub
	Queue *NucleoList
}
func newProducer(chain string, brokers []string, hub *NucleoHub) *NucleoProducer {
	producer := new(NucleoProducer)
	producer.Brokers = brokers
	producer.Hub = hub
	producer.Chain = chain
	producer.Queue = newList()
	producer.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: producer.Brokers,
		Topic: producer.Chain,
		Balancer: &kafka.LeastBytes{},
	})
	// write
	go producer.writeThread();

	return producer;
}

func (p * NucleoProducer) writeThread(){
	for{
		if p.Queue.Size>0 {
			item := p.Queue.Pop()
			if item != nil {
				dataJson, err := json.Marshal(item.Data)
				if err != nil {
					fmt.Println(err)
				}
				key, err := uuid.NewRandom()
				if err != nil {
					fmt.Println(err)
				}
				erro := p.Writer.WriteMessages(context.Background(), kafka.Message{
					Key: []byte(key.String()),
					Value: dataJson,
				})
				if erro != nil {
					fmt.Println(erro)
				}
			}
		}
	}
}
