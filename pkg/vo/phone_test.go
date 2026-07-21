package vo_test

import (
	"errors"
	"testing"

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
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if phone.String() != tc.want {
				t.Fatalf("got %q, want %q", phone.String(), tc.want)
			}
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
			if !errors.Is(err, vo.ErrInvalidPhone) {
				t.Fatalf("got %v, want ErrInvalidPhone", err)
			}
		})
	}
}

func TestNewOptionalPhone(t *testing.T) {
	blank := "   "
	valid := "809-555-1234"

	for name, raw := range map[string]*string{"nil": nil, "blank": &blank} {
		t.Run(name+" is zero", func(t *testing.T) {
			p, err := vo.NewOptionalPhone(raw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !p.IsZero() || p.Ptr() != nil {
				t.Fatal("optional phone should be the zero value")
			}
		})
	}

	p, err := vo.NewOptionalPhone(&valid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := p.Ptr(); got == nil || *got != "8095551234" {
		t.Fatalf("Ptr() = %v, want 8095551234", got)
	}
}
