package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWT secret key - in production, this should be loaded from environment variables
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// CustomClaims represents the JWT claims with user ID
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a given user ID with 6-hour expiry
func GenerateToken(userID uuid.UUID) (string, error) {
	// Set token expiration time to 6 hours from now
	expirationTime := time.Now().Add(6 * time.Hour)

	// Create the claims
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "custom-rugs-app",
		},
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID if valid
func ValidateToken(tokenString string) (uuid.UUID, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return uuid.Nil, errors.New("invalid token")
}

// SetJWTSecret allows setting a custom JWT secret (useful for testing or configuration)
func SetJWTSecret(secret []byte) {
	jwtSecret = secret
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Extract the token from the header
		tokenString := extractTokenFromHeader(authHeader)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token format",
			})
			return
		}

		// Parse and validate the token
		_, err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}
		// Continue to the next handler
		c.Next()
	}
}

// extractTokenFromHeader extracts the token from the Authorization header
func extractTokenFromHeader(header string) string {
	// Check if the header has the format: "Bearer {token}"
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(header, "Bearer ")
}
