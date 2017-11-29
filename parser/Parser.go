package malparser

import (
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/animenotifier/mal"
)

var malAnimeIDRegEx = regexp.MustCompile(`myanimelist.net/anime/(\d+)`)

// ParseAnime ...
func ParseAnime(htmlReader io.Reader) (*mal.Anime, error) {
	document, err := goquery.NewDocumentFromReader(htmlReader)

	if err != nil {
		return nil, err
	}

	anime := &mal.Anime{}

	// Find ID
	document.Find("#horiznav_nav ul li a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Details" {
			anime.URL = s.AttrOr("href", "")
			matches := malAnimeIDRegEx.FindStringSubmatch(anime.URL)

			if len(matches) > 1 {
				anime.ID = matches[1]
			}
		}
	})

	// Information
	document.Find(".dark_text").Each(func(i int, s *goquery.Selection) {
		category := strings.TrimSuffix(s.Text(), ":")

		switch category {
		case "Genres":
			s.Siblings().Each(func(i int, s *goquery.Selection) {
				text := s.Text()

				if text == "add some" {
					return
				}

				anime.Genres = append(anime.Genres, text)
			})

		case "Studios":
			s.Siblings().Each(func(i int, s *goquery.Selection) {
				text := s.Text()

				if text == "add some" {
					return
				}

				anime.Studios = append(anime.Studios, text)
			})

		case "Producers":
			s.Siblings().Each(func(i int, s *goquery.Selection) {
				text := s.Text()

				if text == "add some" {
					return
				}

				anime.Producers = append(anime.Producers, text)
			})

		case "Licensors":
			s.Siblings().Each(func(i int, s *goquery.Selection) {
				text := s.Text()

				if text == "add some" {
					return
				}

				anime.Licensors = append(anime.Licensors, text)
			})
		}
	})

	return anime, nil
}
