package nucleohub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"time"
)

type NucleoProducer struct {
	Brokers []string
	Writer *kafka.Writer
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
	// write
	go producer.WriteThread();

	return producer;
}

func (p * NucleoProducer) WriteThread(){
	u, _ := uuid.NewRandom()
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: p.Hub.Name+"-PRODUCER-"+u.String(),
	}
	p.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: p.Brokers,
		Dialer: dialer,
		Topic: p.Chain,
		Balancer: &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Millisecond,
		BatchTimeout: 500 * time.Millisecond,
		ReadTimeout: 500 * time.Millisecond,
		CompressionCodec: snappy.NewCompressionCodec(),
		Async: false,
	})
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
		time.Sleep(1 * time.Millisecond)
	}
}
