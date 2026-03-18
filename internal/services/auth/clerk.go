package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// VerifyToken validates a session token with Clerk and returns the Clerk user ID.
func (s *ClerkService) VerifyToken(sessionToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, clerkBaseURL+"/clients/verify", nil)
	if err != nil {
		return "", fmt.Errorf("clerk: build verify request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.secretKey)

	q := req.URL.Query()
	q.Set("token", sessionToken)
	req.URL.RawQuery = q.Encode()

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("clerk: verify token: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("clerk: verify token read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("clerk: verify token failed with status %d", resp.StatusCode)
	}

	// The client verify response contains a list of sessions, each with a user_id.
	var result struct {
		Sessions []struct {
			UserID string `json:"user_id"`
		} `json:"sessions"`
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("clerk: verify token unmarshal: %w", err)
	}

	if len(result.Sessions) == 0 || result.Sessions[0].UserID == "" {
		return "", fmt.Errorf("clerk: verify token: no active session found")
	}

	return result.Sessions[0].UserID, nil
}

// CreateSession creates a Clerk session for the given user and returns a session token.
func (s *ClerkService) CreateSession(clerkUserID string) (string, error) {
	payload := map[string]any{
		"user_id": clerkUserID,
	}

	var result struct {
		ID             string `json:"id"`
		LastActiveToken struct {
			JWT string `json:"jwt"`
		} `json:"last_active_token"`
	}

	if _, err := s.clerkDo(http.MethodPost, "/sessions", payload, &result); err != nil {
		return "", fmt.Errorf("clerk: create session: %w", err)
	}

	// Prefer the JWT token; fall back to the session ID.
	if result.LastActiveToken.JWT != "" {
		return result.LastActiveToken.JWT, nil
	}
	if result.ID != "" {
		return result.ID, nil
	}

	return "", fmt.Errorf("clerk: create session: no token in response")
}

// RevokeSession revokes a Clerk session by token/session ID.
func (s *ClerkService) RevokeSession(sessionToken string) error {
	path := fmt.Sprintf("/sessions/%s/revoke", sessionToken)
	if _, err := s.clerkDo(http.MethodDelete, path, nil, nil); err != nil {
		// Log-worthy but non-fatal — caller decides whether to surface this.
		return fmt.Errorf("clerk: revoke session: %w", err)
	}
	return nil
}
