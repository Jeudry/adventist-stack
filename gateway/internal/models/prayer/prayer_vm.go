package prayer

import "github.com/Jeudry/adventist-stack/gateway/internal/models/base"

type PrayerVM struct {
	base.BaseVM
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"authorName"`
	IsAnonymous bool   `json:"isAnonymous"`
	Status      string `json:"status"`
}
