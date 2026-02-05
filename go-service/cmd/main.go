package main

import (
	"context"
	"fmt"
	"go-servie/internals/config"
	"go-servie/internals/db"
	"go-servie/internals/handler"
	kafkapkg "go-servie/internals/pkgs/kafka"
	"go-servie/internals/repo"
	"go-servie/internals/routes"
	"go-servie/internals/service"
	"go-servie/utils"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("err: ", err)
	}
	database, err := db.ConnectDB()
	if err != nil {
		log.Println("database connection failed")
	}
	producer, err := kafkapkg.NewKafkaProducer("kafka:9092")
	if err != nil {
		log.Println("error while producer: ", err)
	}
	firebaseapp := config.Firebase_Setup()
	db := db.DB
	rdc := config.RedisConfig()
	if err := rdc.Ping(context.Background()).Err(); err != nil {
		log.Println("redis failed to connect: ", err)
	} else {
		log.Println("redis connected")
	}

	repo := repo.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service, firebaseapp, producer, rdc)
	router := routes.V1Router(handler)
	server := &http.Server{
		Addr:    utils.PORT,
		Handler: router,
	}
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	go kafkapkg.StartAIResultConsumer(shutdown, service)
	go func() {
		log.Println("Server is running")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	<-shutdown.Done()
	log.Println("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}
	log.Println("Closing resources...")
	database.Close()
	log.Println("Server exited gracefully")
}
