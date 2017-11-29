package mal

// Anime represents an anime on MyAnimeList.
type Anime struct {
	ID        string
	Genres    []string
	Studios   []string
	Producers []string
	Licensors []string
}
