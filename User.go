package mal

// User ...
type User struct {
	ID                string `json:"user_id" xml:"user_id"`
	Name              string `json:"user_name" xml:"user_name"`
	Watching          string `json:"user_watching" xml:"user_watching"`
	Completed         string `json:"user_completed" xml:"user_completed"`
	Onhold            string `json:"user_onhold" xml:"user_onhold"`
	Dropped           string `json:"user_dropped" xml:"user_dropped"`
	Plantowatch       string `json:"user_plantowatch" xml:"user_plantowatch"`
	DaysSpentWatching string `json:"user_days_spent_watching" xml:"user_days_spent_watching"`
}
