package models

import "time"

// ---------------------------------------------------------------------------
// Auth request / response types
// ---------------------------------------------------------------------------

// RegisterParentRequest is the payload for the parent registration endpoint.
type RegisterParentRequest struct {
	FullName        string `json:"full_name"        binding:"required"`
	Email           string `json:"email"            binding:"required,email"`
	Phone           string `json:"phone"`
	Location        string `json:"location"         binding:"required"`
	PrimaryLocation string `json:"primary_location"`
	Occupation      string `json:"occupation"       binding:"required"`
	PreferredHours  string `json:"preferred_hours"  binding:"required"`
	Password        string `json:"password"         binding:"required,min=8"`
}

// RegisterBabysitterRequest is the payload for the babysitter registration endpoint.
type RegisterBabysitterRequest struct {
	FullName  string   `json:"full_name" binding:"required"`
	Email     string   `json:"email"     binding:"required,email"`
	Phone     string   `json:"phone"`
	Location  string   `json:"location"  binding:"required"`
	Languages []string `json:"languages" binding:"required"`
	Password  string   `json:"password"  binding:"required,min=8"`
}

// LoginRequest is the payload for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is returned after a successful login.
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// ---------------------------------------------------------------------------
// User types
// ---------------------------------------------------------------------------

// UserResponse is the public representation of a user returned by the API.
type UserResponse struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateAdminRequest is the payload for creating an admin account.
type CreateAdminRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
}

// ---------------------------------------------------------------------------
// Babysitter profile types
// ---------------------------------------------------------------------------

// BabysitterProfileResponse is the public representation of a babysitter profile.
type BabysitterProfileResponse struct {
	UserID            string   `json:"user_id"`
	FullName          string   `json:"full_name"`
	Email             string   `json:"email"`
	Phone             string   `json:"phone"`
	Location          string   `json:"location"`
	ProfilePictureURL string   `json:"profile_picture_url"`
	Languages         []string `json:"languages"`
	DaysPerWeek       int      `json:"days_per_week"`
	HoursPerDay       int      `json:"hours_per_day"`
	RateType          string   `json:"rate_type"`
	RateAmount        float64  `json:"rate_amount"`
	PaymentMethod     string   `json:"payment_method"`
	IsApproved        bool     `json:"is_approved"`
	Gender            string   `json:"gender"`
	Availability      []string `json:"availability"`
	Currency          string   `json:"currency"`
	IsAvailable       bool     `json:"is_available"`
}

// UpdateBabysitterProfileRequest is the payload for updating a babysitter profile.
type UpdateBabysitterProfileRequest struct {
	Location      string   `json:"location"`
	Languages     []string `json:"languages"`
	DaysPerWeek   int      `json:"days_per_week"`
	HoursPerDay   int      `json:"hours_per_day"`
	RateType      string   `json:"rate_type"`
	RateAmount    float64  `json:"rate_amount"`
	PaymentMethod string   `json:"payment_method"`
	Gender        string   `json:"gender"`
	Availability  []string `json:"availability"`
	Currency      string   `json:"currency"`
}

// SetWorkStatusRequest is the payload for updating a babysitter's work status.
// IsAvailable uses a pointer so that an explicit `false` value is accepted by the binding validator.
type SetWorkStatusRequest struct {
	IsAvailable *bool `json:"is_available" binding:"required"`
}

// ---------------------------------------------------------------------------
// Parent profile types
// ---------------------------------------------------------------------------

// UpdateParentProfileRequest is the payload for updating a parent profile.
type UpdateParentProfileRequest struct {
	Location        string `json:"location"`
	Occupation      string `json:"occupation"`
	PreferredHours  string `json:"preferred_hours"`
	PrimaryLocation string `json:"primary_location"`
}

// ---------------------------------------------------------------------------
// Saved babysitters types
// ---------------------------------------------------------------------------

// SaveBabysitterRequest is the payload for saving a babysitter.
type SaveBabysitterRequest struct {
	BabysitterID string `json:"babysitter_id" binding:"required"`
}

// ---------------------------------------------------------------------------
// Messaging types
// ---------------------------------------------------------------------------

// SendMessageRequest is the payload for sending a message.
type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// MessageResponse is the public representation of a single message.
type MessageResponse struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	Content        string    `json:"content"`
	IsRead         bool      `json:"is_read"`
	SentAt         time.Time `json:"sent_at"`
}

// ConversationResponse is the public representation of a conversation thread.
type ConversationResponse struct {
	ID                         string    `json:"id"`
	OtherUserName              string    `json:"other_user_name"`
	OtherUserProfilePictureURL string    `json:"other_user_profile_picture_url"`
	IsLocked                   bool      `json:"is_locked"`
	CreatedAt                  time.Time `json:"created_at"`
}

// ConversationPreviewResponse is the public representation of a conversation
// enriched with the most recent message — used by the conversation list UI to
// display message snippets and unread indicators.
type ConversationPreviewResponse struct {
	ConversationID             string    `json:"conversation_id"`
	OtherUserName              string    `json:"other_user_name"`
	OtherUserProfilePictureURL string    `json:"other_user_profile_picture_url"`
	IsLocked                   bool      `json:"is_locked"`
	LastMessageID              string    `json:"last_message_id"`
	LastSenderID               string    `json:"last_sender_id"`
	PreviewText                string    `json:"preview_text"`
	IsRead                     bool      `json:"is_read"`
	LastMessageSent            time.Time `json:"last_message_sent"`
}

// ---------------------------------------------------------------------------
// Profile view types
// ---------------------------------------------------------------------------

// ProfileViewResponse is returned when listing who viewed a babysitter's profile.
// Restricted fields (Email, Phone, PrimaryLocation, PreferredHours) are only
// populated when HasMessaged is true — the babysitter and parent have exchanged messages.
type ProfileViewResponse struct {
	ID              string    `json:"id"`
	ParentID        string    `json:"parent_id"`
	ParentName      string    `json:"parent_name"`
	Occupation      string    `json:"occupation"`
	ViewedAt        time.Time `json:"viewed_at"`
	HasMessaged     bool      `json:"has_messaged"`
	Email           string    `json:"email,omitempty"`
	Phone           string    `json:"phone,omitempty"`
	PrimaryLocation string    `json:"primary_location,omitempty"`
	PreferredHours  string    `json:"preferred_hours,omitempty"`
}

// WeeklyViewsResponse is returned for the weekly profile views count endpoint.
type WeeklyViewsResponse struct {
	BabysitterID string `json:"babysitter_id"`
	ViewCount    int64  `json:"view_count"`
	PeriodDays   int    `json:"period_days"`
}

// ---------------------------------------------------------------------------
// Admin / reporting types
// ---------------------------------------------------------------------------

// ActivityResponse is used in admin reporting to summarise user activity.
type ActivityResponse struct {
	UserID        string `json:"user_id"`
	FullName      string `json:"full_name"`
	Role          string `json:"role"`
	ActivityLabel string `json:"activity_label"`
	MessageCount  int64  `json:"message_count"`
}

// ---------------------------------------------------------------------------
// Generic response types
// ---------------------------------------------------------------------------

// ErrorResponse is returned when a request fails.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse is returned for operations that produce no data payload.
type SuccessResponse struct {
	Message string `json:"message"`
}
