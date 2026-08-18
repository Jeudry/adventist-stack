package vo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/pkg/vo"
)

func TestNewPhone_Valid(t *testing.T) {
	cases := map[string]struct {
		raw  string
		want string
	}{
		"plain digits":        {raw: "8095551234", want: "8095551234"},
		"strips separators":   {raw: "809-555-1234", want: "8095551234"},
		"strips parens/space": {raw: "+1 (809) 555-1234", want: "+18095551234"},
		"trims whitespace":    {raw: "  8095551234  ", want: "8095551234"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			phone, err := vo.NewPhone(tc.raw)
			require.NoError(t, err)
			assert.Equal(t, tc.want, phone.String())
		})
	}
}

func TestNewPhone_Invalid(t *testing.T) {
	cases := map[string]string{
		"empty":          "",
		"whitespace":     "   ",
		"too short":      "12345",
		"too long":       "+1234567890123456789012",
		"letters":        "809ABC1234",
		"plus in middle": "809+5551234",
	}

	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := vo.NewPhone(raw)
			assert.ErrorIs(t, err, vo.ErrInvalidPhone)
		})
	}
}

func TestNewOptionalPhone(t *testing.T) {
	blank := "   "
	valid := "809-555-1234"

	for name, raw := range map[string]*string{"nil": nil, "blank": &blank} {
		t.Run(name+" is zero", func(t *testing.T) {
			p, err := vo.NewOptionalPhone(raw)
			require.NoError(t, err)
			assert.True(t, p.IsZero(), "an absent phone must be the zero value")
			assert.Nil(t, p.Ptr(), "Ptr() must stay nil so the column stores NULL")
		})
	}

	p, err := vo.NewOptionalPhone(&valid)
	require.NoError(t, err)
	require.NotNil(t, p.Ptr(), "Ptr()")
	assert.Equal(t, "8095551234", *p.Ptr(), "Ptr()")
}
