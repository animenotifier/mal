package mal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/aerogo/http/client"
)

// AnimeList is just a slice of anime list items.
type AnimeList []AnimeListItem

// GetAnimeList returns the user's anime list.
func GetAnimeList(userName string) (AnimeList, error) {
	animeList := AnimeList{}

	// Fetch the page
	response, err := client.Get("https://myanimelist.net/animelist/" + userName + "?status=7").End()

	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Status code: %d", response.StatusCode())
	}

	// Parse the HTML
	reader := bytes.NewReader(response.Bytes())
	document, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, err
	}

	dataItems, exists := document.Find(".list-table").First().Attr("data-items")

	if !exists {
		return nil, errors.New("Missing data-items attribute")
	}

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
