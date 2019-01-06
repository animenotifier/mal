package mal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnimeList(t *testing.T) {
	userNames := []string{
		"Aky",
		"soory",
		"Subpyro",
		"PaladinRaid",
	}

	for _, userName := range userNames {
		testUser(t, userName)
	}
}

func testUser(t *testing.T, userName string) {
	animeList, err := GetAnimeList(userName)

	assert.NoError(t, err)
	assert.NotEmpty(t, animeList)

	for _, item := range animeList {
		assert.NotNil(t, item)
		assert.True(t, item.AnimeID > 0)
		assert.True(t, item.NumWatchedEpisodes >= 0)
		assert.True(t, item.Score >= 0)
		assert.True(t, item.Score <= 10)
		assert.NotEmpty(t, item.AnimeTitle)

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
