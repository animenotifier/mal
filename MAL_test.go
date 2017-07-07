package mal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnimeList(t *testing.T) {
	animeList, err := GetAnimeList("Aky")
	assert.NoError(t, err)
	fmt.Println(len(animeList.Anime), "anime")
}
