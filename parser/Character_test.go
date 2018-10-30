package malparser_test

import (
	"bytes"
	"testing"

	"github.com/aerogo/http/client"
	"github.com/animenotifier/arn/stringutils"

	"github.com/animenotifier/mal/parser"

	"github.com/stretchr/testify/assert"
)

func TestParseCharacter(t *testing.T) {
	response, err := client.Get("https://myanimelist.net/character/117909").End()
	assert.NoError(t, err)

	character, err := malparser.ParseCharacter(bytes.NewReader(response.Bytes()))
	assert.NoError(t, err)
	stringutils.PrettyPrint(character)
}
