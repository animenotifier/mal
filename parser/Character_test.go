package malparser_test

import (
	"bytes"
	"testing"

	"github.com/aerogo/http/client"

	"github.com/animenotifier/mal/parser"

	"github.com/stretchr/testify/assert"
)

func TestParseCharacter(t *testing.T) {
	response, err := client.Get("https://myanimelist.net/character/63").End()
	assert.NoError(t, err)

	character, err := malparser.ParseCharacter(bytes.NewReader(response.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, "Winry Rockbell", character.Name)
	// stringutils.PrettyPrint(character)
}
