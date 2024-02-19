package schedule

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFixNijisanjiWikiRawDateNumString(t *testing.T) {
	assert.Equal(t, FixNijisanjiWikiRawDateNumString("01"), "01")
	assert.Equal(t, FixNijisanjiWikiRawDateNumString("1"), "01")
	assert.Equal(t, FixNijisanjiWikiRawDateNumString(""), "00")
	assert.Equal(t, FixNijisanjiWikiRawDateNumString("10or11"), "10")
	assert.Equal(t, FixNijisanjiWikiRawDateNumString("朝の10"), "10")
}

func TestParseNijisanjiWikiRawDate(t *testing.T) {
	a, _ := ParseNijisanjiWikiRawDate(2023, 12, 1, "2時00分")
	b, _ := time.Parse("2006/01/02 15:04:05", "2023/12/01 02:00:00")
	assert.Equal(t, a, b)

	a, _ = ParseNijisanjiWikiRawDate(2026, 2, 1, "09時30分")
	b, _ = time.Parse("2006/01/02 15:04:05", "2026/02/01 09:30:00")
	assert.Equal(t, a, b)

	a, _ = ParseNijisanjiWikiRawDate(2023, 12, 1, "(？時40分")
	b, _ = time.Parse("2006/01/02 15:04:05", "2023/12/01 00:40:00")
	assert.Equal(t, a, b)

	a, _ = ParseNijisanjiWikiRawDate(2023, 12, 1, "(20時？？分")
	b, _ = time.Parse("2006/01/02 15:04:05", "2023/12/01 20:00:00")
	assert.Equal(t, a, b)

	a, _ = ParseNijisanjiWikiRawDate(2023, 12, 1, "18or20時40分")
	b, _ = time.Parse("2006/01/02 15:04:05", "2023/12/01 18:40:00")
	assert.Equal(t, a, b)
}
