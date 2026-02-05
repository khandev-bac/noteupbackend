package config

import (
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func ImagekitConfig() *imagekit.Client {
	imagekitApp := imagekit.NewClient(
		option.WithBaseURL("https://upload.imagekit.io"),
		option.WithPrivateKey("private_3eARtsyol0qhbkWmgiO5P2V4zqA="),
	)
	return &imagekitApp
}
