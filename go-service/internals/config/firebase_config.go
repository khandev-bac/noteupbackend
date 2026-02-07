package config

import (
	"context"
	"go-servie/utils"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func Firebase_Setup() *firebase.App {
	firebase_file := utils.FIREBASE_SERVICE
	opt := option.WithCredentialsFile(firebase_file)
	config := &firebase.Config{ProjectID: "notev2-34951"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalln("Firebase setup failed", err)
	}
	return app
}
