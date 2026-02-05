package kafkapkg

import (
	"context"
	"encoding/json"
	"go-servie/internals/service"
	"go-servie/utils"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

func StartAIResultConsumer(ctx context.Context, svc *service.Service) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092"),
		"group.id":           "ai-result-worker",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	_ = c.SubscribeTopics([]string{"aiwork.done"}, nil)
	log.Println("✅ AI result worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Kafka consumer shutting down")
			return
		default:
			msg, err := c.ReadMessage(500 * time.Millisecond)
			if err != nil {
				continue
			}

			var evt utils.NoteEventResponse
			if err := json.Unmarshal(msg.Value, &evt); err != nil {
				_, _ = c.CommitMessage(msg)
				continue
			}

			log.Println("📥 AI result received:", evt.NoteId)

			noteID, err := uuid.Parse(evt.NoteId)
			if err != nil {
				_, _ = c.CommitMessage(msg)
				continue
			}

			_, err = svc.CompleteProcessing(
				context.Background(),
				noteID,
				evt.Title,
				evt.Transcript,
				int32(len(evt.Transcript)),
			)
			if err != nil {
				log.Println("DB update failed:", err)
				continue // retry
			}

			_, _ = c.CommitMessage(msg)
			log.Println("✅ Note updated:", evt.NoteId)
		}
	}
}
