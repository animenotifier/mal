package malparser

import (
	"fmt"
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
func ParseAnime(htmlReader io.Reader) (*mal.Anime, []*mal.Character, error) {
	document, err := goquery.NewDocumentFromReader(htmlReader)

	if err != nil {
		return nil, nil, err
	}

	anime := &mal.Anime{}
	characters := []*mal.Character{}

	// Title
	title := document.Find("h1.h1 span[itemprop='name']").Text()
	title = strings.TrimSpace(title)
	anime.Title = title

	// Characters
	document.Find("div.detail-characters-list > div > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		children := s.Children()

		// Disregard staff fields
		if children.Length() != 3 {
			return
		}

		image := children.Eq(0).Find("img")
		info := children.Eq(1)
		// voiceActor := children.Eq(2)

		// ID, link and name
		link := info.Find("a")
		name := strings.TrimSpace(link.Text())
		url := link.AttrOr("href", "")
		id := strings.TrimPrefix(url, "https://myanimelist.net/character/")
		slashPos := strings.Index(id, "/")

		if slashPos != -1 {
			id = id[:slashPos]
		}

		// If ID is empty, something went wrong
		if id == "" {
			fmt.Println("No ID found for the character:", name)
			return
		}

		// Image
		imgSrc := image.AttrOr("data-src", "")
		imgSrc = strings.Replace(imgSrc, "/r/23x32", "", 1)
		imgSrc = strings.TrimPrefix(imgSrc, "https://myanimelist.cdn-dena.com/images/characters/")
		queryPos := strings.Index(imgSrc, "?s")

		if queryPos != -1 {
			imgSrc = imgSrc[:queryPos]
		}

		if strings.Contains(imgSrc, "questionmark") {
			imgSrc = ""
		}

		// Role
		role := link.Next().Text()
		role = strings.TrimSpace(role)
		role = strings.ToLower(role)

		// Create character
		character := &mal.Character{
			ID:        id,
			Name:      name,
			ImagePath: imgSrc,
		}

		characters = append(characters, character)

		// Create anime relation
		animeCharacter := &mal.AnimeCharacter{
			ID:   id,
			Role: role,
		}

		anime.Characters = append(anime.Characters, animeCharacter)
	})

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

	// Find image from og:image tag
	document.Find("meta[property='og:image']").Each(func(i int, s *goquery.Selection) {
		imageURL := s.AttrOr("content", "")

		if !strings.Contains(imageURL, "/images/anime/") || strings.Contains(imageURL, "/icon/") {
			return
		}

		imageURL = strings.Replace(imageURL, ".jpg", "l.jpg", 1)
		anime.Image = imageURL
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

		// Remove fake synopsis
		if strings.HasPrefix(anime.Synopsis, "Looking for information on ") {
			anime.Synopsis = ""
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

		case "Duration":
			duration := darkTextValue(s)
			duration = strings.TrimSuffix(duration, " min. per ep.")
			episodeLength, err := strconv.Atoi(duration)

			if err == nil {
				anime.EpisodeLength = episodeLength
			}

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
			status := darkTextValue(s)

			switch status {
			case "Finished Airing":
				anime.Status = "finished"

			case "Currently Airing":
				anime.Status = "current"

			case "Not yet aired":
				anime.Status = "tba"

			default:
				anime.Status = status
			}

		case "Source":
			anime.Source = darkTextValue(s)
			anime.Source = strings.ToLower(anime.Source)

			if anime.Source == "unknown" {
				anime.Source = ""
			}

		case "Rating":
			anime.Rating = darkTextValue(s)

		case "Aired":
			aired := s.Parent().Contents().Not(".dark_text").Text()
			aired = strings.TrimSpace(aired)
			parts := strings.Split(aired, " to ")
			startDate := parts[0]
			endDate := ""

			if len(parts) > 1 {
				endDate = parts[1]
			}

			if startDate == "?" {
				startDate = ""
			}

			if endDate == "?" {
				endDate = ""
			}

			startTime, err := time.Parse(malDateFormat, startDate)

			if err == nil {
				anime.StartDate = startTime.Format("2006-01-02")

				if anime.Status == "tba" {
					anime.Status = "upcoming"
				}
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

	return anime, characters, nil
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
