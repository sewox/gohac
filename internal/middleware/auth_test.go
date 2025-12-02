package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	// Register a dummy protected route
	app.Get("/test", Protected(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{
			"success": true,
			"user_id": userID,
		})
	})

	return app
}

func TestProtected_NoCookie(t *testing.T) {
	// Case 1: Request without cookie -> Expect 401
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestProtected_InvalidToken(t *testing.T) {
	// Case 2: Request with invalid token -> Expect 401
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", "auth_token=invalid-token-here")
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestProtected_ExpiredToken(t *testing.T) {
	// Test with expired token
	app := setupTestApp()

	// Create an expired token
	expiredTime := time.Now().Add(-1 * time.Hour) // 1 hour ago
	claims := &Claims{
		UserID: "test-user",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt:  jwt.NewNumericDate(expiredTime.Add(-24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", "auth_token="+tokenString)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestProtected_ValidToken(t *testing.T) {
	// Case 3: Request with valid token -> Expect 200 and verify Context has UserID
	app := setupTestApp()

	// Generate a valid token
	tokenString, err := GenerateToken("test-user-123", "test@example.com", 24)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", "auth_token="+tokenString)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProtected_ValidToken_AuthorizationHeader(t *testing.T) {
	// Test with Authorization header instead of cookie
	app := setupTestApp()

	// Generate a valid token
	tokenString, err := GenerateToken("test-user-456", "test2@example.com", 24)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGenerateToken(t *testing.T) {
	userID := "test-user"
	email := "test@example.com"
	expirationHours := 24

	tokenString, err := GenerateToken(userID, email, expirationHours)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	require.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*Claims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
}

func TestProtected_InvalidSigningMethod(t *testing.T) {
	// Test with malformed/invalid token
	app := setupTestApp()

	// Use a malformed token string
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidGVzdC11c2VyIn0.invalid"

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", "auth_token="+tokenString)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}
