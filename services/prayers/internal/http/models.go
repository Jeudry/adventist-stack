package http

import (
	"github.com/Jeudry/adventist-stack/pkg/httpx"
)

type PrayerRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"authorName" required:"false"`
	IsAnonymous bool   `json:"isAnonymous" required:"false"`
	Status      string `json:"status" required:"false"`
}

type PrayerVM struct {
	httpx.BaseVM
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"authorName" required:"false"`
	IsAnonymous bool   `json:"isAnonymous"`
	Status      string `json:"status"`
}
