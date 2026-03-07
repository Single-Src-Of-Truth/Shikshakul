package controller

import (
	"net/http"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/dto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/response"
	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	profileService service.ProfileService
}

func NewProfileController(ps service.ProfileService) *ProfileController {
	return &ProfileController{profileService: ps}
}

func (ctrl *ProfileController) GetMyProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	profile, err := ctrl.profileService.GetMyProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch profile")
		return
	}

	response.Success(c, profile)
}

func (ctrl *ProfileController) GetMySessions(c *gin.Context) {
	userID := c.GetString("user_id")

	rawToken, err := c.Cookie("skl_session")
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Session cookie missing")
		return
	}

	sessions, err := ctrl.profileService.GetMySessions(c.Request.Context(), userID, rawToken)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch active devices")
		return
	}

	response.SuccessWithMsg(c, "Active sessions retrieved", sessions)
}

func (ctrl *ProfileController) UpdateMyProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := ctrl.profileService.UpdateMyProfile(c.Request.Context(), userID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	response.SuccessWithMsg(c, "Profile updated successfully", nil)
}
