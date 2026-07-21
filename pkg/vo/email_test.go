package vo_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Jeudry/adventist-stack/pkg/vo"
)

func TestNewEmail_Valid(t *testing.T) {
	cases := map[string]struct {
		raw  string
		want string
	}{
		"already normalized": {raw: "pastor@church.org", want: "pastor@church.org"},
		"trims whitespace":   {raw: "  pastor@church.org  ", want: "pastor@church.org"},
		"lowercases":         {raw: "Pastor@Church.ORG", want: "pastor@church.org"},
		"plus addressing":    {raw: "member+news@church.org", want: "member+news@church.org"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			email, err := vo.NewEmail(tc.raw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if email.String() != tc.want {
				t.Fatalf("got %q, want %q", email.String(), tc.want)
			}
		})
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	cases := map[string]string{
		"empty":      "",
		"whitespace": "   ",
		"no at":      "pastorchurch.org",
		"no domain":  "pastor@",
		"no tld":     "pastor@church",
		"double dot": "pastor@@church.org",
		"too long":   strings.Repeat("a", 250) + "@church.org",
	}

	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := vo.NewEmail(raw)
			if !errors.Is(err, vo.ErrInvalidEmail) {
				t.Fatalf("got %v, want ErrInvalidEmail", err)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	a, _ := vo.NewEmail("Pastor@Church.org")
	b, _ := vo.NewEmail("pastor@church.org")
	c, _ := vo.NewEmail("elder@church.org")

	if !a.Equals(b) {
		t.Fatal("normalized emails should be equal")
	}
	if a.Equals(c) {
		t.Fatal("different emails should not be equal")
	}
}

func TestEmail_IsZero(t *testing.T) {
	var zero vo.Email
	if !zero.IsZero() {
		t.Fatal("zero value should report IsZero")
	}

	email, _ := vo.NewEmail("pastor@church.org")
	if email.IsZero() {
		t.Fatal("constructed email should not be zero")
	}
}
