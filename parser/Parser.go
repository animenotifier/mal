package malparser

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/animenotifier/mal"
)

const malDateFormat = "Jan _2, 2006"

var (
	malAnimeIDRegEx    = regexp.MustCompile(`myanimelist.net/anime/(\d+)`)
	malProducerIDRegEx = regexp.MustCompile(`/anime/producer/(\d+)`)
	sourceRegex        = regexp.MustCompile(`\(Source: (.*?)\)`)
	writtenByRegex     = regexp.MustCompile(`\[Written by (.*?)\]`)
)

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

	// Find description from og:description tag
	document.Find("meta[property='og:description']").Each(func(i int, s *goquery.Selection) {
		synopsis := s.AttrOr("content", "")
		matches := writtenByRegex.FindStringSubmatch(synopsis)

		if len(matches) >= 2 {
			anime.SynopsisSource = matches[1]
		}

		synopsis = writtenByRegex.ReplaceAllString(synopsis, "")
		synopsis = strings.TrimSpace(synopsis)
		anime.Synopsis = synopsis
	})

	// Title
	title := document.Find("h1.h1 span[itemprop='name']").Text()
	title = strings.TrimSpace(title)
	anime.Title = title

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

		case "English":
			anime.EnglishTitle = darkTextValue(s)

		case "Japanese":
			anime.JapaneseTitle = darkTextValue(s)

		case "Synonyms":
			anime.Synonyms = strings.Split(darkTextValue(s), ", ")

		case "Episodes":
			text := darkTextValue(s)
			number, err := strconv.Atoi(text)

			if err == nil {
				anime.EpisodeCount = number
			}

		case "Status":
			anime.Status = darkTextValue(s)

		case "Source":
			anime.Source = darkTextValue(s)

		case "Rating":
			anime.Rating = darkTextValue(s)

		case "Aired":
			aired := s.Parent().Contents().Not(".dark_text").Text()
			aired = strings.TrimSpace(aired)
			parts := strings.Split(aired, " to ")
			startDate := parts[0]
			endDate := parts[1]

			if startDate == "?" {
				startDate = ""
			}

			if endDate == "?" {
				endDate = ""
			}

			startTime, err := time.Parse(malDateFormat, startDate)

			if err == nil {
				anime.StartDate = startTime.Format("2006-01-02")
			}

			endTime, err := time.Parse(malDateFormat, endDate)

			if err == nil {
				anime.EndDate = endTime.Format("2006-01-02")
			}

		case "Studios":
			s.Siblings().Each(producerHandler(&anime.Studios))

		case "Producers":
			s.Siblings().Each(producerHandler(&anime.Producers))

		case "Licensors":
			s.Siblings().Each(producerHandler(&anime.Licensors))
		}
	})

	return anime, nil
}

func darkTextValue(s *goquery.Selection) string {
	text := s.Parent().Contents().Not(".dark_text").Text()
	text = strings.TrimSpace(text)
	return text
}

func producerHandler(slice *[]*mal.Producer) func(int, *goquery.Selection) {
	return func(i int, s *goquery.Selection) {
		text := s.Text()

		if text == "add some" {
			return
		}

		id := ""
		url := s.AttrOr("href", "")
		matches := malProducerIDRegEx.FindStringSubmatch(url)

		if len(matches) > 1 {
			id = matches[1]
		}

		*slice = append(*slice, &mal.Producer{
			ID:   id,
			Name: text,
		})
	}
}
