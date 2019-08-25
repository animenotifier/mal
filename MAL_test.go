package mal

import (
	"testing"

	"github.com/akyoto/assert"
)

func TestAnimeList(t *testing.T) {
	userNames := []string{
		"Aky",
	}

	for _, userName := range userNames {
		testUser(t, userName)
	}
}

func testUser(t *testing.T, userName string) {
	animeList, err := GetAnimeList(userName)
	assert.Nil(t, err)
	assert.NotEqual(t, len(animeList), 0)

	for _, item := range animeList {
		assert.NotNil(t, item)
		assert.True(t, item.AnimeID > 0)
		assert.True(t, item.NumWatchedEpisodes >= 0)
		assert.True(t, item.Score >= 0)
		assert.True(t, item.Score <= 10)
		assert.NotEqual(t, item.AnimeTitle, "")

		// Status can only be one of the given values
		switch item.Status {
		case AnimeListStatusWatching:
		case AnimeListStatusCompleted:
		case AnimeListStatusPlanned:
		case AnimeListStatusHold:
		case AnimeListStatusDropped:
		default:
			t.Fail()
		}
	}
}
