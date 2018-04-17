package mal

// Anime represents an anime on MyAnimeList.
type Anime struct {
	ID             string
	URL            string
	Title          string
	EnglishTitle   string
	JapaneseTitle  string
	Image          string
	Synonyms       []string
	Synopsis       string
	Rating         string
	SynopsisSource string
	Status         string
	Source         string
	StartDate      string
	EndDate        string
	EpisodeCount   int
	EpisodeLength  int
	Genres         []string
	Studios        []*Producer
	Producers      []*Producer
	Licensors      []*Producer
	Characters     []*AnimeCharacter
}
