package mal

import "fmt"

// AnimeListItem represents a single list item.
type AnimeListItem struct {
	AnimeID            int  `json:"anime_id"`
	Status             int  `json:"status"`
	Score              int  `json:"score"`
	IsRewatching       bool `json:"is_rewatching"`
	NumWatchedEpisodes int  `json:"num_watched_episodes"`

	// Tags                  string      `json:"tags"`
	// AnimeTitle            string      `json:"anime_title"`
	// AnimeNumEpisodes      int         `json:"anime_num_episodes"`
	// AnimeAiringStatus     int         `json:"anime_airing_status"`
	// AnimeStudios          interface{} `json:"anime_studios"`
	// AnimeLicensors        interface{} `json:"anime_licensors"`
	// AnimeSeason           interface{} `json:"anime_season"`
	// HasEpisodeVideo       bool        `json:"has_episode_video"`
	// HasPromotionVideo     bool        `json:"has_promotion_video"`
	// HasVideo              bool        `json:"has_video"`
	// VideoURL              string      `json:"video_url"`
	// AnimeURL              string      `json:"anime_url"`
	// AnimeImagePath        string      `json:"anime_image_path"`
	// IsAddedToList         bool        `json:"is_added_to_list"`
	// AnimeMediaTypeString  string      `json:"anime_media_type_string"`
	// AnimeMpaaRatingString string      `json:"anime_mpaa_rating_string"`
	// StartDateString       interface{} `json:"start_date_string"`
	// FinishDateString      interface{} `json:"finish_date_string"`
	// AnimeStartDateString  string      `json:"anime_start_date_string"`
	// AnimeEndDateString    string      `json:"anime_end_date_string"`
	// DaysString            interface{} `json:"days_string"`
	// StorageString         string      `json:"storage_string"`
	// PriorityString        string      `json:"priority_string"`
}

// AnimeLink ...
func (item AnimeListItem) AnimeLink() string {
	return fmt.Sprintf("https://myanimelist.net/anime/%d", item.AnimeID)
}
