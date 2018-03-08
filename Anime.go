package mal

// Anime represents an anime on MyAnimeList.
type Anime struct {
	ID        string
	URL       string
	Synopsis  string
	Source    string
	StartDate string
	EndDate   string
	Genres    []string
	Studios   []*Producer
	Producers []*Producer
	Licensors []*Producer
}
