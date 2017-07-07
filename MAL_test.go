package mal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnimeList(t *testing.T) {
	animeList, err := GetAnimeList("Aky")

	assert.NoError(t, err)
	assert.NotNil(t, animeList)
	assert.NotNil(t, animeList.Items)

	fmt.Println(len(animeList.Items), "anime")

	for _, item := range animeList.Items {
		assert.NotEmpty(t, item.AnimeID)

		// MyStatus can only be one of the given values
		switch item.MyStatus {
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
