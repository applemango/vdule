package youtube

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseChannelHandle(t *testing.T) {
	assert.Equal(t, ParseChannelHandle("@example"), "example")
	assert.Equal(t, ParseChannelHandle("example"), "example")
	assert.Equal(t, ParseChannelHandle("@exAmpLe"), "example")
	assert.Equal(t, ParseChannelHandle("EXampLE"), "example")
}
