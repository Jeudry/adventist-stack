package prayer

type UpdatePrayerRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"authorName"`
	IsAnonymous bool   `json:"isAnonymous"`
}
