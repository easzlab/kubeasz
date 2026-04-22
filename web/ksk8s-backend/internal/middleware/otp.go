package middleware

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/http"
	"time"

	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// GenerateOTPSecret creates a new random TOTP secret.
func GenerateOTPSecret() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateOTPURL returns the provisioning URI for QR code generation.
func GenerateOTPURL(secret string, username string, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, username, secret, issuer)
}

// ValidateOTP verifies a TOTP code against a secret.
func ValidateOTP(secret string, code string) bool {
	return totp.Validate(code, secret)
}

// OTPVerifiedSessions holds recently verified OTP sessions (userID -> expiry).
var OTPVerifiedSessions = make(map[int64]time.Time)

// OTPVerifyMiddleware checks if user has OTP enabled and validates X-OTP-Code header.
// For high-risk operations, the client must provide a valid OTP code in the X-OTP-Code header.
func OTPVerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		uid := userID.(int64)

		userRepo := repository.NewUserRepository()
		user, err := userRepo.GetByID(uid)
		if err != nil {
			c.Next()
			return
		}

		// If OTP not enabled, skip
		if !user.OTPEnabled {
			c.Next()
			return
		}

		// Check recent verified session (within 5 min)
		if expiry, ok := OTPVerifiedSessions[uid]; ok && time.Now().Before(expiry) {
			c.Next()
			return
		}

		// Require X-OTP-Code header
		code := c.GetHeader("X-OTP-Code")
		if code == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "otp_required", "message": "OTP verification required for this operation"})
			return
		}

		if !ValidateOTP(user.OTPSecret, code) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid_otp", "message": "Invalid OTP code"})
			return
		}

		// Mark session as verified for 5 minutes
		OTPVerifiedSessions[uid] = time.Now().Add(5 * time.Minute)
		c.Next()
	}
}

// ClearOTPVerifiedSession removes the verified session for a user (e.g., on logout).
func ClearOTPVerifiedSession(userID int64) {
	delete(OTPVerifiedSessions, userID)
}
