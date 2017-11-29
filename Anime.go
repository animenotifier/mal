package mal

// Anime represents an anime on MyAnimeList.
type Anime struct {
	ID        string
	URL       string
	Genres    []string
	Studios   []string
	Producers []string
	Licensors []string
}
