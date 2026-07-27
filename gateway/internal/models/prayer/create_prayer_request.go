package prayer

type CreatePrayerRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"authorName"`
	IsAnonymous bool   `json:"isAnonymous"`
	Status      string `json:"status"`
}
