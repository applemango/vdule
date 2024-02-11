package youtube

import (
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
	apiKey := os.Getenv("YOUTUBE_DATA_API_KEY")
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	return &Tube{service}, nil
}
