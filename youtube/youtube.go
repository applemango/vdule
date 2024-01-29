package youtube

import (
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
)

type Tube struct {
	Service *youtube.Service
}

func CreateTube() (*Tube, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyA2KkxjSd-s5ydoVPCym9yOH9lsyInxoKE"},
	}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	return &Tube{service}, nil
}
