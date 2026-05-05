package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const clerkBaseURL = "https://api.clerk.com/v1"

// ClerkService handles all communication with the Clerk backend API.
type ClerkService struct {
	secretKey  string
	httpClient *http.Client
}

// NewClerkService constructs a ClerkService with the provided secret key.
func NewClerkService(secretKey string) *ClerkService {
	return &ClerkService{
		secretKey:  secretKey,
		httpClient: &http.Client{},
	}
}

// clerkDo is a helper that executes an HTTP request with Clerk auth headers and
// decodes the JSON response body into dest. Pass nil for dest to discard the body.
func (s *ClerkService) clerkDo(method, path string, body any, dest any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("clerk: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, clerkBaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("clerk: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.secretKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("clerk: do request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("clerk: read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("clerk: request failed with status %d: %s", resp.StatusCode, string(respBytes))
	}

	if dest != nil {
		if err := json.Unmarshal(respBytes, dest); err != nil {
			return resp, fmt.Errorf("clerk: unmarshal response: %w", err)
		}
	}

	return resp, nil
}

// CreateUser creates a new user in Clerk and returns their Clerk user ID.
func (s *ClerkService) CreateUser(email, password, fullName string) (string, error) {
	payload := map[string]any{
		"email_address": []string{email},
		"password":      password,
		"first_name":    fullName,
	}

	var result struct {
		ID string `json:"id"`
	}

	if _, err := s.clerkDo(http.MethodPost, "/users", payload, &result); err != nil {
		return "", fmt.Errorf("clerk: create user: %w", err)
	}

	if result.ID == "" {
		return "", fmt.Errorf("clerk: create user: empty id in response")
	}

	return result.ID, nil
}

// DeleteUser removes a user from Clerk. Used for rollback on failed registration.
func (s *ClerkService) DeleteUser(clerkUserID string) error {
	if _, err := s.clerkDo(http.MethodDelete, "/users/"+clerkUserID, nil, nil); err != nil {
		return fmt.Errorf("clerk: delete user: %w", err)
	}
	return nil
}

// VerifyPassword checks whether password matches the Clerk-stored credential for
// the given user. Returns (true, nil) on success, (false, nil) on a wrong
// password, and (false, err) on any transport or API error.
func (s *ClerkService) VerifyPassword(clerkUserID, password string) (bool, error) {
	var result struct {
		Verified bool `json:"verified"`
	}
	_, err := s.clerkDo(
		http.MethodPost,
		"/users/"+clerkUserID+"/verify_password",
		map[string]any{"password": password},
		&result,
	)
	if err != nil {
		// Clerk returns 422 for a wrong password — treat that as (false, nil).
		if isClerkVerificationFailure(err) {
			return false, nil
		}
		return false, fmt.Errorf("clerk: verify password: %w", err)
	}
	return result.Verified, nil
}

// isClerkVerificationFailure returns true when a Clerk API error indicates that
// the supplied password was simply wrong (HTTP 422 / 400) rather than a
// transport or server error that the caller should surface as a 500.
func isClerkVerificationFailure(err error) bool {
	msg := err.Error()
	// Clerk returns 422 Unprocessable Entity for wrong passwords.
	return strings.Contains(msg, "status 422") || strings.Contains(msg, "status 400")
}

// GenerateToken creates a signed JWT for the given Clerk user ID, expiring at expiry.
func (s *ClerkService) GenerateToken(clerkUserID string, expiry time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub": clerkUserID,
		"exp": expiry.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("clerk: sign token: %w", err)
	}
	return signed, nil
}

// VerifyToken parses a signed JWT and returns the Clerk user ID from the "sub" claim.
func (s *ClerkService) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("clerk: unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("clerk: verify token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("clerk: verify token: invalid token")
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return "", fmt.Errorf("clerk: verify token: missing sub claim")
	}

	return sub, nil
}
