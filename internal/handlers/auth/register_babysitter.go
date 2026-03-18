package auth

import (
	"database/sql"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"babycare-api/internal/db"
	"babycare-api/internal/models"
	services_auth "babycare-api/internal/services/auth"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 20 << 20 // 20 MB

// RegisterBabysitter creates a new babysitter account with profile and document uploads.
// Expects multipart/form-data with text fields and four file fields.
func (h *AuthHandler) RegisterBabysitter(c *gin.Context) {
	// Parse the multipart form with a 20 MB limit.
	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "request too large or invalid multipart form"})
		return
	}

	// --- Text fields ---
	fullName := strings.TrimSpace(c.Request.FormValue("full_name"))
	email := strings.TrimSpace(c.Request.FormValue("email"))
	phone := strings.TrimSpace(c.Request.FormValue("phone"))
	location := strings.TrimSpace(c.Request.FormValue("location"))
	password := c.Request.FormValue("password")
	languagesRaw := strings.TrimSpace(c.Request.FormValue("languages"))

	// Validate required text fields.
	switch "" {
	case fullName:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "full_name is required"})
		return
	case email:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "email is required"})
		return
	case phone:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "phone is required"})
		return
	case location:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "location is required"})
		return
	case password:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "password is required"})
		return
	case languagesRaw:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "languages is required"})
		return
	}

	if len(password) < 8 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "password must be at least 8 characters"})
		return
	}

	// Parse comma-separated languages into a slice.
	var languages []string
	for _, lang := range strings.Split(languagesRaw, ",") {
		if trimmed := strings.TrimSpace(lang); trimmed != "" {
			languages = append(languages, trimmed)
		}
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	// Check for duplicate email.
	_, err := queries.GetUserByEmail(ctx, email)
	if err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: "email already registered"})
		return
	}
	if err != sql.ErrNoRows {
		log.Printf("register_babysitter: check email: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// --- File fields ---
	nationalIDHeader, err := fileFromForm(c, "national_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "national_id file is required"})
		return
	}
	lciLetterHeader, err := fileFromForm(c, "lci_letter")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "lci_letter file is required"})
		return
	}
	cvHeader, err := fileFromForm(c, "cv")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "cv file is required"})
		return
	}
	profilePictureHeader, err := fileFromForm(c, "profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "profile_picture file is required"})
		return
	}

	// Upload files to Cloudinary.
	folder := "babycare/babysitters"

	nationalIDURL, err := h.storageService.UploadFile(nationalIDHeader, folder)
	if err != nil {
		log.Printf("register_babysitter: upload national_id: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to upload national_id"})
		return
	}
	lciLetterURL, err := h.storageService.UploadFile(lciLetterHeader, folder)
	if err != nil {
		log.Printf("register_babysitter: upload lci_letter: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to upload lci_letter"})
		return
	}
	cvURL, err := h.storageService.UploadFile(cvHeader, folder)
	if err != nil {
		log.Printf("register_babysitter: upload cv: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to upload cv"})
		return
	}
	profilePictureURL, err := h.storageService.UploadFile(profilePictureHeader, folder)
	if err != nil {
		log.Printf("register_babysitter: upload profile_picture: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to upload profile_picture"})
		return
	}

	// Create user in Clerk.
	clerkUserID, err := h.clerkService.CreateUser(email, password, fullName)
	if err != nil {
		log.Printf("register_babysitter: clerk create user: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create account"})
		return
	}

	// Hash the password for local storage.
	hash, err := services_auth.HashPassword(password)
	if err != nil {
		log.Printf("register_babysitter: hash password: %v", err)
		_ = h.clerkService.DeleteUser(clerkUserID)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Persist the user row.
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		FullName:     fullName,
		Email:        email,
		Phone:        sql.NullString{String: phone, Valid: phone != ""},
		Role:         db.UserRoleBabysitter,
		PasswordHash: sql.NullString{String: hash, Valid: true},
		ClerkUserID:  sql.NullString{String: clerkUserID, Valid: true},
	})
	if err != nil {
		log.Printf("register_babysitter: db create user: %v", err)
		_ = h.clerkService.DeleteUser(clerkUserID)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to save user"})
		return
	}

	// Create the babysitter profile.
	_, err = queries.CreateBabysitterProfile(ctx, db.CreateBabysitterProfileParams{
		UserID:            user.ID,
		Location:          sql.NullString{String: location, Valid: location != ""},
		NationalIDUrl:     sql.NullString{String: nationalIDURL, Valid: nationalIDURL != ""},
		LciLetterUrl:      sql.NullString{String: lciLetterURL, Valid: lciLetterURL != ""},
		CvUrl:             sql.NullString{String: cvURL, Valid: cvURL != ""},
		ProfilePictureUrl: sql.NullString{String: profilePictureURL, Valid: profilePictureURL != ""},
		Languages:         languages,
	})
	if err != nil {
		log.Printf("register_babysitter: create babysitter profile: %v", err)
		// Profile creation failure is logged; user row exists and can be retried.
	}

	c.JSON(http.StatusCreated, models.UserResponse{
		ID:        user.ID.String(),
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     user.Phone.String,
		Role:      string(user.Role),
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
	})
}

// fileFromForm extracts a single file header from the multipart form.
func fileFromForm(c *gin.Context, field string) (*multipart.FileHeader, error) {
	_, fh, err := c.Request.FormFile(field)
	return fh, err
}
