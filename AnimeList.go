package mal

// AnimeList ...
type AnimeList struct {
	Myinfo *User    `json:"myinfo" xml:"myinfo"`
	Anime  []*Anime `json:"anime" xml:"anime>"`
}
