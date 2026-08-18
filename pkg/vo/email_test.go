package vo_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			require.NoError(t, err)
			assert.Equal(t, tc.want, email.String())
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
			assert.ErrorIs(t, err, vo.ErrInvalidEmail)
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	a, _ := vo.NewEmail("Pastor@Church.org")
	b, _ := vo.NewEmail("pastor@church.org")
	c, _ := vo.NewEmail("elder@church.org")

	assert.True(t, a.Equals(b), "%q and %q normalize to the same address", a.String(), b.String())
	assert.False(t, a.Equals(c), "%q and %q are different addresses", a.String(), c.String())
}

func TestEmail_IsZero(t *testing.T) {
	var zero vo.Email
	assert.True(t, zero.IsZero(), "the zero value must report IsZero")

	email, _ := vo.NewEmail("pastor@church.org")
	assert.False(t, email.IsZero(), "a constructed email must not report IsZero")
}
