package nucleohub

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type NucleoConsumer struct {
	Brokers []string
	Chain string
	GroupID string
	Reader *kafka.Reader
	Hub *NucleoHub
}

func NewConsumer(chain string, group string, brokers []string, hub *NucleoHub) *NucleoConsumer {
	consumer := new(NucleoConsumer);
	consumer.Brokers = brokers
	consumer.GroupID = group
	consumer.Hub = hub
	consumer.Chain = chain
	consumer.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   consumer.Brokers,
		GroupID:   consumer.GroupID,
		Topic:     consumer.Chain,
		MinBytes:  1, // 1B
		MaxBytes:  2e6, // 2MB
		CommitInterval: 1000,
	})
	go consumer.readThread()

	return consumer;
}
func (c * NucleoConsumer) readThread(){
	for {
		m, err := c.Reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}