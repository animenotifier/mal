package mal

import (
	"encoding/xml"

	"github.com/aerogo/http/client"
)

// AnimeList ...
type AnimeList struct {
	Myinfo *User            `json:"myinfo" xml:"myinfo"`
	Items  []*AnimeListItem `json:"anime" xml:"anime"`
}

// GetAnimeList returns the user's anime list.
func GetAnimeList(userName string) (*AnimeList, error) {
	animeList := &AnimeList{}
	response, err := client.Get("https://myanimelist.net/malappinfo.php?u=" + userName + "&status=all&type=anime").End()

	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(response.Bytes(), &animeList)

	if err != nil {
		return nil, err
	}

	return animeList, nil
}
