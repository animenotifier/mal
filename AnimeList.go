package mal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aerogo/http/client"
)

// AnimeList is just a slice of anime list items.
type AnimeList []*AnimeListItem

// GetAnimeList returns the user's anime list.
func getAnimeList(userName string, page int) (AnimeList, error) {
	animeList := AnimeList{}

	offset := page * 300

	// Fetch the page
	url := fmt.Sprintf("https://myanimelist.net/animelist/%s/load.json?offset=%d&status=7", userName, offset)
	response, err := client.Get(url).End()

	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Status code: %d", response.StatusCode())
	}
	// Parse the HTML

	if err != nil {
		return nil, err
	}

	dataItems := response.String()

	// Fix is_rewatching field
	dataItems = strings.Replace(dataItems, `"is_rewatching":""`, `"is_rewatching":false`, -1)
	dataItems = strings.Replace(dataItems, `"is_rewatching":0`, `"is_rewatching":false`, -1)
	dataItems = strings.Replace(dataItems, `"is_rewatching":1`, `"is_rewatching":true`, -1)
	// Parse JSON
	err = json.Unmarshal([]byte(dataItems), &animeList)

	if err != nil {
		return nil, err
	}

	return animeList, nil
}

func GetAnimeList(userName string) (AnimeList, error) {
	animeList := AnimeList{}
	page := 0
	ticker := time.NewTicker(1100 * time.Millisecond)
	rateLimit := ticker.C
	defer ticker.Stop()

	for {
		nextAnimeList, err := getAnimeList(userName, page)

		if err != nil {
			return nil, err
		}

		if len(nextAnimeList) == 0 {
			break
		}

		animeList = append(animeList, nextAnimeList...)
		page++
		// Wait for rate limiter to allow the next request
		<-rateLimit
	}

	return animeList, nil
}
