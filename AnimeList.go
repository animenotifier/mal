package mal

import (
	"encoding/xml"

	"github.com/parnurzeal/gorequest"
)

// AnimeList ...
type AnimeList struct {
	Myinfo *User            `json:"myinfo" xml:"myinfo"`
	Items  []*AnimeListItem `json:"anime" xml:"anime"`
}

// GetAnimeList returns the user's anime list.
func GetAnimeList(userName string) (*AnimeList, error) {
	animeList := &AnimeList{}

	_, body, errs := gorequest.New().Get("https://myanimelist.net/malappinfo.php?u=" + userName + "&status=all&type=anime").EndBytes()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	err := xml.Unmarshal(body, &animeList)

	if err != nil {
		return nil, err
	}

	return animeList, nil
}
