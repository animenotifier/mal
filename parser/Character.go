package malparser

import (
	"bytes"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/animenotifier/mal"
)

var (
	malCharacterIDRegEx = regexp.MustCompile(`myanimelist.net/character/(\d+)`)
)

// ParseCharacter ...
func ParseCharacter(htmlReader io.Reader) (*mal.Character, error) {
	document, err := goquery.NewDocumentFromReader(htmlReader)

	if err != nil {
		return nil, err
	}

	character := &mal.Character{}

	// Find name from og:title tag
	s := document.Find("meta[property='og:title']").First()
	name := s.AttrOr("content", "")
	character.Name = html.UnescapeString(name)

	// Find headers
	rightColumn := document.Find("#content > table > tbody > tr > td").Last()
	nameHeader := rightColumn.Find(".normal_header").First()

	// Find Japanese name
	s = nameHeader.Find("small").First()
	japaneseName := strings.TrimSpace(s.Text())
	japaneseName = strings.TrimPrefix(japaneseName, "(")
	japaneseName = strings.TrimSuffix(japaneseName, ")")
	character.JapaneseName = japaneseName

	// Description
	character.Description = getDescription(rightColumn)
	character.Spoilers = getSpoilers(rightColumn)

	// Find image from og:image tag
	s = document.Find("meta[property='og:image']").First()
	imageURL := s.AttrOr("content", "")

	if strings.HasPrefix(imageURL, "https://myanimelist.cdn-dena.com/images/characters/") {
		character.Image = strings.TrimPrefix(imageURL, "https://myanimelist.cdn-dena.com/images/characters/")
	}

	// Find ID
	document.Find("#horiznav_nav ul li a").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "Details" {
			return
		}

		character.URL = s.AttrOr("href", "")
		matches := malCharacterIDRegEx.FindStringSubmatch(character.URL)

		if len(matches) > 1 {
			character.ID = matches[1]
		}
	})

	return character, nil
}

func getSpoilers(s *goquery.Selection) []string {
	spoilers := []string{}

	s.Find(".spoiler_content").Each(func(i int, s *goquery.Selection) {
		spoiler := getFlatText(s)
		spoilers = append(spoilers, spoiler)
	})

	return spoilers
}

func getDescription(s *goquery.Selection) string {
	description := getFlatText(s)

	if description == "No biography written." {
		description = ""
	}

	return description
}

// getFlatText returns the text of a node using only top-level children (no recursion).
// It additionally handles cases for inline formatting elements.
func getFlatText(s *goquery.Selection) string {
	description := bytes.Buffer{}

	s.Contents().Each(func(i int, s *goquery.Selection) {
		switch goquery.NodeName(s) {
		case "#text":
			text := s.Text()
			description.WriteString(text)

		case "b", "i", "em", "strong":
			text := strings.TrimSpace(s.Text())
			description.WriteString(text)
		}
	})

	result := description.String()
	result = strings.Replace(result, "\n\n\n", "\n\n", -1)
	result = strings.TrimSpace(result)
	return result
}
