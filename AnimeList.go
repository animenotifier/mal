package mal

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aerogo/http/client"
	jsoniter "github.com/json-iterator/go"
)

// AnimeList is just a slice of anime list items.
type AnimeList []*AnimeListItem

// GetAnimeList returns the user's anime list.
func GetAnimeList(userName string) (AnimeList, error) {
	animeList := AnimeList{}
	offset := 0
	ticker := time.NewTicker(1100 * time.Millisecond)
	defer ticker.Stop()

	for {
		nextAnimeList, err := getAnimeListWithOffset(userName, offset)

		if err != nil {
			return nil, err
		}

		if len(nextAnimeList) == 0 {
			break
		}

		animeList = append(animeList, nextAnimeList...)
		offset += len(nextAnimeList)

		// Wait for rate limiter to allow the next request
		<-ticker.C
	}

	return animeList, nil
}

// getAnimeListWithOffset returns the anime list items starting with the given offset.
func getAnimeListWithOffset(userName string, offset int) (AnimeList, error) {
	var response *client.Response
	var err error
	animeList := AnimeList{}

	// Fetch the page
	for {
		url := fmt.Sprintf("https://myanimelist.net/animelist/%s/load.json?offset=%d&status=7", userName, offset)
		response, err = client.Get(url).End()

		if err != nil {
			return nil, err
		}

		if !response.Ok() {
			if response.StatusCode() == http.StatusTooManyRequests {
				time.Sleep(time.Second)
				continue
			}

			return nil, fmt.Errorf("Status code: %d", response.StatusCode())
		}

		break
	}

	// Get JSON string
	dataItems := response.String()

	// Fix is_rewatching field
	dataItems = strings.ReplaceAll(dataItems, `"is_rewatching":""`, `"is_rewatching":false`)
	dataItems = strings.ReplaceAll(dataItems, `"is_rewatching":0`, `"is_rewatching":false`)
	dataItems = strings.ReplaceAll(dataItems, `"is_rewatching":1`, `"is_rewatching":true`)

	// Fix anime_title field sometimes including numbers instead of strings
	dataItems = strings.ReplaceAll(dataItems, `"anime_title":1`, `"anime_title":"1"`)

	// Parse JSON
	err = jsoniter.UnmarshalFromString(dataItems, &animeList)
	fmt.Println(err)

	if err != nil {
		return nil, err
	}

	return animeList, nil
}
