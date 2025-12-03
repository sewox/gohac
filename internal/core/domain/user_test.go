package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_HashPassword(t *testing.T) {
	user := &User{
		Password: "testpassword123",
	}

	err := user.HashPassword()
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, "testpassword123", user.Password) // Should be hashed
	assert.True(t, len(user.Password) > 50)              // bcrypt hash is long
}

func TestUser_CheckPassword(t *testing.T) {
	user := &User{
		Password: "testpassword123",
	}

	// Hash the password first
	err := user.HashPassword()
	assert.NoError(t, err)

	// Test correct password
	assert.True(t, user.CheckPassword("testpassword123"))

	// Test incorrect password
	assert.False(t, user.CheckPassword("wrongpassword"))
	assert.False(t, user.CheckPassword(""))
}

func TestUser_CheckPassword_WithDifferentHashes(t *testing.T) {
	// Test that same password produces different hashes (due to salt)
	user1 := &User{Password: "testpassword123"}
	user2 := &User{Password: "testpassword123"}

	err1 := user1.HashPassword()
	err2 := user2.HashPassword()

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// Hashes should be different (bcrypt uses random salt)
	assert.NotEqual(t, user1.Password, user2.Password)

	// But both should verify the same password
	assert.True(t, user1.CheckPassword("testpassword123"))
	assert.True(t, user2.CheckPassword("testpassword123"))
}
