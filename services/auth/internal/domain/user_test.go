package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

func TestNewUser_Valid(t *testing.T) {
	user, err := domain.NewUser("  Pastor@Church.org ", "  John Doe  ", domain.RoleMember)
	require.NoError(t, err)

	assert.Equal(t, "pastor@church.org", user.Email.String(), "the email must be normalized")
	assert.Equal(t, "John Doe", user.Name, "the name must be trimmed")
}

func TestNewUser_Invalid(t *testing.T) {
	cases := map[string]struct {
		email, name string
		role        domain.Role
		wantErr     error
	}{
		"bad email":    {email: "nope", name: "John", role: domain.RoleMember, wantErr: vo.ErrInvalidEmail},
		"empty name":   {email: "a@b.co", name: "  ", role: domain.RoleMember, wantErr: domain.ErrInvalidUser},
		"invalid role": {email: "a@b.co", name: "John", role: domain.Role("bishop"), wantErr: domain.ErrInvalidUser},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := domain.NewUser(tc.email, tc.name, tc.role)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestUser_SetPassword(t *testing.T) {
	user, err := domain.NewUser("a@b.co", "John", domain.RoleMember)
	require.NoError(t, err, "build user")

	assert.ErrorIs(t, user.SetPassword("short"), domain.ErrInvalidUser, "short password")

	require.NoError(t, user.SetPassword("supersecret"), "valid password")
	assert.NotEmpty(t, user.PasswordHash, "password_hash")
}

func TestUser_Authenticate(t *testing.T) {
	user, err := domain.NewUser("a@b.co", "John", domain.RoleMember)
	require.NoError(t, err, "build user")
	require.NoError(t, user.SetPassword("supersecret"), "setup")

	assert.NoError(t, user.Authenticate("supersecret"), "the correct password must be accepted")
	assert.ErrorIs(t, user.Authenticate("wrong"), domain.ErrInvalidCredentials, "wrong password")
}

func TestUser_IsAdmin(t *testing.T) {
	admin, err := domain.NewUser("a@b.co", "Admin", domain.RoleAdmin)
	require.NoError(t, err, "build admin")
	member, err := domain.NewUser("c@d.co", "Member", domain.RoleMember)
	require.NoError(t, err, "build member")

	assert.True(t, admin.IsAdmin(), "role %q", admin.Role)
	assert.False(t, member.IsAdmin(), "role %q", member.Role)
}
