package schedule

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTwitterUrl(t *testing.T) {
	assert.Equal(t, ParseTwitterUrl("https://twitter.com/MahiroYukishiro/status/1759187824618991840"), "MahiroYukishiro")
	assert.Equal(t, ParseTwitterUrl("https://twitter.com/ars_almal/status/1759232329720189154"), "ars_almal")
	assert.Equal(t, ParseTwitterUrl("https://twitter.com/NaSera2434/status/1759225097687277683"), "NaSera2434")
}

func TestParseYoutubeUrl(t *testing.T) {
	assert.Equal(t, ParseYoutubeUrl("https://www.youtube.com/watch?v=4-f4hWEzkCo"), "4-f4hWEzkCo")
	assert.Equal(t, ParseYoutubeUrl("https://www.youtube.com/watch?v=IZwfdJKxEhQ"), "IZwfdJKxEhQ")
	assert.Equal(t, ParseYoutubeUrl("https://www.youtube.com/watch?v=IZwfdJKxEhQ&source=www.google.com"), "IZwfdJKxEhQ")
}
