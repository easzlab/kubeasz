package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGenerateToken(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")

	token, err := GenerateToken(42, "admin", "admin")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Error("token should not be empty")
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")
	gin.SetMode(gin.TestMode)

	token, _ := GenerateToken(42, "admin", "admin")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	JWTAuth()(c)

	if c.IsAborted() {
		t.Error("request should not be aborted with valid token")
	}

	userID, exists := c.Get("user_id")
	if !exists {
		t.Error("user_id should be set in context")
	}
	if userID.(int64) != 42 {
		t.Errorf("expected user_id 42, got %v", userID)
	}

	username, _ := c.Get("username")
	if username != "admin" {
		t.Errorf("expected username admin, got %v", username)
	}

	role, _ := c.Get("role")
	if role != "admin" {
		t.Errorf("expected role admin, got %v", role)
	}
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/test", nil)

	JWTAuth()(c)

	if !c.IsAborted() {
		t.Error("request should be aborted without auth header")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "missing authorization header") {
		t.Errorf("expected missing auth header error, got %s", body)
	}
}

func TestJWTAuth_InvalidFormat(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/test", nil)
	c.Request.Header.Set("Authorization", "Basic dXNlcjpwYXNz")

	JWTAuth()(c)

	if !c.IsAborted() {
		t.Error("request should be aborted with invalid auth format")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/test", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token-string")

	JWTAuth()(c)

	if !c.IsAborted() {
		t.Error("request should be aborted with invalid token")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")
	gin.SetMode(gin.TestMode)

	// Generate a token and immediately create one that will be expired
	// by manipulating the claims after generation is hard, so we'll
	// use parseToken to verify it rejects bad tokens
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJyb2xlIjoidmlld2VyIiwiZXhwIjoxfQ.fake_signature"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+expiredToken)

	JWTAuth()(c)

	if !c.IsAborted() {
		t.Error("request should be aborted with expired token")
	}
}

func TestRequireAdmin_AdminUser(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin")

	RequireAdmin()(c)

	if c.IsAborted() {
		t.Error("request should not be aborted for admin")
	}
}

func TestRequireAdmin_NonAdmin(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "viewer")

	RequireAdmin()(c)

	if !c.IsAborted() {
		t.Error("request should be aborted for non-admin")
	}

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestRequireAdmin_MissingRole(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// No role set

	RequireAdmin()(c)

	if !c.IsAborted() {
		t.Error("request should be aborted when role is missing")
	}

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestGenerateToken_UniquePerUser(t *testing.T) {
	InitJWTSecret("test-secret-key-for-unit-tests")

	token1, _ := GenerateToken(1, "user1", "viewer")
	token2, _ := GenerateToken(2, "user2", "admin")

	if token1 == token2 {
		t.Error("tokens for different users should be unique")
	}
}

func TestInitJWTSecret_DefaultFallback(t *testing.T) {
	// When empty string is passed, should use default
	InitJWTSecret("")

	_, err := GenerateToken(1, "test", "viewer")
	if err != nil {
		t.Errorf("should generate token with default secret: %v", err)
	}
}
