package models

import "time"

// ---------------------------------------------------------------------------
// Auth request / response types
// ---------------------------------------------------------------------------

// RegisterParentRequest is the payload for the parent registration endpoint.
type RegisterParentRequest struct {
	FullName       string `json:"full_name"       binding:"required"`
	Email          string `json:"email"           binding:"required,email"`
	Phone          string `json:"phone"`
	Location       string `json:"location"        binding:"required"`
	Occupation     string `json:"occupation"      binding:"required"`
	PreferredHours string `json:"preferred_hours" binding:"required"`
	Password       string `json:"password"        binding:"required,min=8"`
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
}

// ---------------------------------------------------------------------------
// Parent profile types
// ---------------------------------------------------------------------------

// UpdateParentProfileRequest is the payload for updating a parent profile.
type UpdateParentProfileRequest struct {
	Location       string `json:"location"`
	Occupation     string `json:"occupation"`
	PreferredHours string `json:"preferred_hours"`
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
	ID            string    `json:"id"`
	OtherUserName string    `json:"other_user_name"`
	IsLocked      bool      `json:"is_locked"`
	CreatedAt     time.Time `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Profile view types
// ---------------------------------------------------------------------------

// ProfileViewResponse is returned when listing who viewed a babysitter's profile.
type ProfileViewResponse struct {
	ID         string    `json:"id"`
	ParentID   string    `json:"parent_id"`
	ParentName string    `json:"parent_name"`
	ViewedAt   time.Time `json:"viewed_at"`
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
