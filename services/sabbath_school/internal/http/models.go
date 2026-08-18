package http

import "github.com/Jeudry/adventist-stack/pkg/httpx"

// The bounds repeat the domain's constants because a struct tag only takes literals. They are not
// redundant: these are what reach the OpenAPI document, so a client learns the rule instead of
// discovering it through a 422. The domain still validates — this layer never gets to be the only
// guard, since nothing stops another caller from reaching the service.
type SabbathSchoolRequest struct {
	Name         string  `json:"name" minLength:"3" maxLength:"128" doc:"Nombre de la clase"`
	TeacherID    *string `json:"teacherId,omitempty" format:"uuid"`
	Location     *string `json:"location,omitempty" minLength:"2" maxLength:"128"`
	TargetMinAge *int    `json:"targetMinAge,omitempty" minimum:"0" maximum:"120"`
	TargetMaxAge *int    `json:"targetMaxAge,omitempty" minimum:"0" maximum:"120"`
	Status       string  `json:"status,omitempty" doc:"active o inactive; vacío toma active"`
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
