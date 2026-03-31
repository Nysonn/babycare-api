package parent

import (
	"database/sql"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

const maxParentUpdateUploadSize = 10 << 20 // 10 MB

// UpdateProfile allows the authenticated parent to update their profile fields.
func (h *ParentHandler) UpdateProfile(c *gin.Context) {
	currentUserRaw, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorised"})
		return
	}
	currentUser, ok := currentUserRaw.(db.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorised"})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()
	currentProfile, err := queries.GetParentProfileByUserID(ctx, currentUser.ID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "profile not found"})
		return
	}
	if err != nil {
		log.Printf("parent update_profile: get current profile: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	location := currentProfile.Location.String
	occupation := currentProfile.Occupation.String
	preferredHours := currentProfile.PreferredHours.String
	primaryLocation := currentProfile.PrimaryLocation.String
	profilePictureURL := currentProfile.ProfilePictureUrl

	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(maxParentUpdateUploadSize); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "request too large or invalid multipart form"})
			return
		}

		location = formValueOrDefault(c.Request.MultipartForm, "location", location)
		occupation = formValueOrDefault(c.Request.MultipartForm, "occupation", occupation)
		preferredHours = formValueOrDefault(c.Request.MultipartForm, "preferred_hours", preferredHours)
		primaryLocation = formValueOrDefault(c.Request.MultipartForm, "primary_location", primaryLocation)

		if _, fileHeader, err := c.Request.FormFile("profile_picture"); err == nil {
			if h.storageService == nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "storage service unavailable"})
				return
			}

			url, uploadErr := h.storageService.UploadFile(fileHeader, "babycare/parents")
			if uploadErr != nil {
				log.Printf("parent update_profile: upload profile_picture: %v", uploadErr)
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to upload profile picture"})
				return
			}

			profilePictureURL = sql.NullString{String: url, Valid: true}
		}
	} else {
		var req models.UpdateParentProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		location = req.Location
		occupation = req.Occupation
		preferredHours = req.PreferredHours
		primaryLocation = req.PrimaryLocation
	}

	updated, err := queries.UpdateParentProfile(ctx, db.UpdateParentProfileParams{
		UserID:            currentUser.ID,
		Location:          sql.NullString{String: location, Valid: location != ""},
		Occupation:        sql.NullString{String: occupation, Valid: occupation != ""},
		PreferredHours:    sql.NullString{String: preferredHours, Valid: preferredHours != ""},
		PrimaryLocation:   sql.NullString{String: primaryLocation, Valid: primaryLocation != ""},
		ProfilePictureUrl: profilePictureURL,
	})
	if err != nil {
		log.Printf("parent update_profile: db update: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, buildParentResponse(currentUser, updated))
}

func formValueOrDefault(form *multipart.Form, key, fallback string) string {
	if form == nil {
		return fallback
	}

	values, ok := form.Value[key]
	if !ok || len(values) == 0 {
		return fallback
	}

	return values[0]
}
