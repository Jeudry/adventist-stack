package domain_test

import (
	"errors"
	"testing"

	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

func TestNewUser_Valid(t *testing.T) {
	user, err := domain.NewUser("  Pastor@Church.org ", "  John Doe  ", domain.RoleMember)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email.String() != "pastor@church.org" {
		t.Fatalf("email not normalized: %q", user.Email.String())
	}
	if user.Name != "John Doe" {
		t.Fatalf("name not trimmed: %q", user.Name)
	}
}

func TestNewUser_Invalid(t *testing.T) {
	cases := map[string]struct {
		email, name string
		role        domain.Role
		wantErr     error
	}{
		"bad email":    {email: "nope", name: "John", role: domain.RoleMember, wantErr: vo.ErrInvalidEmail},
		"empty name":   {email: "a@b.co", name: "  ", role: domain.RoleMember, wantErr: domain.ErrInvalidUser},
		"invalid role": {email: "a@b.co", name: "John", role: domain.Role(99), wantErr: domain.ErrInvalidUser},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := domain.NewUser(tc.email, tc.name, tc.role)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestUser_SetPassword(t *testing.T) {
	user, _ := domain.NewUser("a@b.co", "John", domain.RoleMember)

	if err := user.SetPassword("short"); !errors.Is(err, domain.ErrInvalidUser) {
		t.Fatalf("short password: got %v, want ErrInvalidUser", err)
	}

	if err := user.SetPassword("supersecret"); err != nil {
		t.Fatalf("valid password: unexpected error: %v", err)
	}
	if user.PasswordHash == "" {
		t.Fatal("password hash was not set")
	}
}

func TestUser_Authenticate(t *testing.T) {
	user, _ := domain.NewUser("a@b.co", "John", domain.RoleMember)
	if err := user.SetPassword("supersecret"); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := user.Authenticate("supersecret"); err != nil {
		t.Fatalf("correct password rejected: %v", err)
	}
	if err := user.Authenticate("wrong"); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("wrong password: got %v, want ErrInvalidCredentials", err)
	}
}

func TestUser_IsAdmin(t *testing.T) {
	admin, _ := domain.NewUser("a@b.co", "Admin", domain.RoleAdmin)
	member, _ := domain.NewUser("c@d.co", "Member", domain.RoleMember)

	if !admin.IsAdmin() {
		t.Fatal("admin should be admin")
	}
	if member.IsAdmin() {
		t.Fatal("member should not be admin")
	}
}
