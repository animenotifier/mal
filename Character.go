package mal

// Character represents a character on MyAnimeList.
type Character struct {
	ID           string
	URL          string
	Image        string
	Name         string
	JapaneseName string
	Description  string
	Spoilers     []string
	// AlternateName string
	// VoiceActors   []*Person
}

// ImageLink returns the URL of the image.
func (character *Character) ImageLink() string {
	return "https://myanimelist.cdn-dena.com/images/characters/" + character.Image
}
