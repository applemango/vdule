package youtube

import (
	"errors"
	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
)

type Tube struct {
	Service *youtube.Service
}

var T, _ = CreateTube()

func CreateTube() (*Tube, error) {
	err := godotenv.Load(".env")
	if err != nil {
		panic("failed load config")
	}
	apiKey := os.Getenv("YOUTUBE_DATA_API_KEY")
	if len(apiKey) == 0 {
		return nil, errors.New("error")
	}
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
		return nil, err
	}
	return &Tube{service}, nil
}

func ResetYoutube() {
	T, _ = CreateTube()
}
