package kafkapkg

import (
	"encoding/json"
	"go-servie/utils"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func AIConsumerWork(out chan<- utils.NoteEventResponse) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "main-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Println("consumer create failed:", err)
		return
	}

	err = c.SubscribeTopics([]string{"aiwork.done"}, nil)
	if err != nil {
		log.Println("subscribe failed:", err)
		return
	}

	log.Println("✅ Go Kafka consumer listening on aiwork.done")

	for {
		msg, err := c.ReadMessage(1 * time.Second)

		if err != nil {
			// 🔇 Ignore timeout (normal Kafka behavior)
			if kafkaErr, ok := err.(kafka.Error); ok &&
				kafkaErr.Code() == kafka.ErrTimedOut {
				continue
			}

			log.Println("Kafka error:", err)
			continue
		}

		var evt utils.NoteEventResponse
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			log.Println("failed to parse json:", err)
			partitions, err := c.CommitMessage(msg)
			if err != nil {
				log.Println("commit failed:", err)
			} else {
				log.Println("committed:", partitions)
			}

			continue
		}

		log.Println("📥 AI result received:", evt.NoteId)

		// Send to service layer
		out <- evt

		// ✅ Commit AFTER successful read
		_, err = c.CommitMessage(msg)
		if err != nil {
			log.Println("commit failed:", err)
		}
	}
}
