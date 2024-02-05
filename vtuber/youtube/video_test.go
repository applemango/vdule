package youtube

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsLive(t *testing.T) {
	tube, err := CreateTube()
	live, err := tube.GetRawVideo("puS_2cpIUQA")
	if err != nil {
		t.Fatalf("tube.GetRawVideo Error: %v\n", err)
	}
	is := isLive(live)
	assert.Equal(t, is, true)

	video, err := tube.GetRawVideo("C7szSb6Cuko")
	if err != nil {
		t.Fatalf("tube.GetRawVideo Error: %v\n", err)
	}
	is = isLive(video)
	assert.Equal(t, is, false)
}
