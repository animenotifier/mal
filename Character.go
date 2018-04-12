package mal

// Character represents a character on MyAnimeList.
type Character struct {
	ID        string
	ImagePath string
	Name      string
	// JapaneseName  string
	// AlternateName string
	// Description   string
	// VoiceActors   []*Person
}

// ImageLink returns the URL of the image.
func (character *Character) ImageLink() string {
	return "https://myanimelist.cdn-dena.com/images/characters/" + character.ImagePath
}
