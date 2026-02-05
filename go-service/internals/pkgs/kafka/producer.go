package kafkapkg

import (
	"encoding/json"
	"go-servie/utils"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	p *kafka.Producer
}

func NewKafkaProducer(broker string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"acks":              "all",
		"retries":           5,
		"linger.ms":         5,
	})
	if err != nil {
		log.Println("Kafka failed producer:", err)
		return nil, err
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Println("❌ delivery failed:", ev.TopicPartition.Error)
				} else {
					log.Println("✅ delivered:", ev.TopicPartition)
				}
			}
		}
	}()

	log.Println("✅ Kafka producer connected")
	return &Producer{p: p}, nil
}

func (p *Producer) Close() {
	log.Println("Closing Kafka producer...")
	p.p.Flush(5000)
	p.p.Close()
}

func (p *Producer) PublishNoteCreated(event, noteId, audioUrl string) error {
	evt := utils.NoteEvent{
		Event:     event,
		NoteId:    noteId,
		AudioUrl:  audioUrl,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	payload, err := json.Marshal(evt)
	if err != nil {
		log.Println("error while sending payload: ", err)
		return err
	}
	topic := "note.created"
	return p.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(noteId),
		Value: payload,
	}, nil)
}
