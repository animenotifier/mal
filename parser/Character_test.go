package malparser_test

import (
	"bytes"
	"testing"

	"github.com/aerogo/http/client"
	"github.com/akyoto/assert"
	malparser "github.com/animenotifier/mal/parser"
)

func TestParseCharacter(t *testing.T) {
	response, err := client.Get("https://myanimelist.net/character/63").End()
	assert.Nil(t, err)

	character, err := malparser.ParseCharacter(bytes.NewReader(response.Bytes()))
	assert.Nil(t, err)
	assert.Equal(t, "Winry Rockbell", character.Name)
	// stringutils.PrettyPrint(character)
}
