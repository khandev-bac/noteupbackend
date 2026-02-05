package handler

import (
	kafkapkg "go-servie/internals/pkgs/kafka"
	"go-servie/internals/service"

	firebase "firebase.google.com/go/v4"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	service     *service.Service
	firebaseApp *firebase.App
	kafkaConfig *kafkapkg.Producer
	redis       *redis.Client
}

func NewHandler(service *service.Service, firebaseApp *firebase.App, kafkaConfig *kafkapkg.Producer, redis *redis.Client) *Handler {
	return &Handler{
		service:     service,
		firebaseApp: firebaseApp,
		kafkaConfig: kafkaConfig,
		redis:       redis,
	}
}
