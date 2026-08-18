package http

import "github.com/Jeudry/adventist-stack/pkg/httpx"

type SabbathSchoolRequest struct {
	Name         string  `json:"name"`
	TeacherID    *string `json:"teacherId,omitempty"`
	Location     *string `json:"location,omitempty"`
	TargetMinAge *int    `json:"targetMinAge,omitempty"`
	TargetMaxAge *int    `json:"targetMaxAge,omitempty"`
	Status       string  `json:"status,omitempty"`
}

type SabbathSchoolVM struct {
	httpx.BaseVM
	Name         string  `json:"name"`
	TeacherID    *string `json:"teacherId,omitempty"`
	Location     *string `json:"location,omitempty"`
	TargetMinAge *int    `json:"targetMinAge,omitempty"`
	TargetMaxAge *int    `json:"targetMaxAge,omitempty"`
	Status       string  `json:"status"`
}
