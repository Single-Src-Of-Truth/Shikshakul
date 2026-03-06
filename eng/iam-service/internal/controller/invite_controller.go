package controller

import (
	"net/http"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/dto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/response"
	"github.com/gin-gonic/gin"
)

type InviteController struct {
	onboardingService service.OnboardingService
}

func NewInviteController(os service.OnboardingService) *InviteController {
	return &InviteController{onboardingService: os}
}

func (ctrl *InviteController) GenerateInvite(c *gin.Context) {
	var req dto.GenerateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid invite payload")
		return
	}

	inviterID := c.GetString("user_id")

	requestTenantID := req.TenantID
	contextTenantID := c.GetString("tenant_id")
	if contextTenantID != "" && contextTenantID != requestTenantID {
		response.Error(c, http.StatusForbidden, "Cannot invite users to a different school")
		return
	}

	inviteID, rawToken, err := ctrl.onboardingService.GenerateInvite(c.Request.Context(), inviterID, requestTenantID, req.RoleID, req.Identifier)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate invite link")
		return
	}

	// TODO: Move it to email after it's done
	response.Created(c, "Invite generated successfully", gin.H{
		"invite_id":    inviteID,
		"invite_token": rawToken,
		"identifier":   req.Identifier,
	})
}

func (ctrl *InviteController) AcceptInvite(c *gin.Context) {
	var req dto.AcceptInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload or password too short (min 8 chars)")
		return
	}

	user, err := ctrl.onboardingService.AcceptInvite(c.Request.Context(), req.InviteToken, req.Password, req.FirstName, req.LastName)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "Account created successfully. You may now log in.", gin.H{
		"user_id":    user.ID,
		"identifier": user.PrimaryIdentifier,
	})
}

func (ctrl *InviteController) ExtendInvite(c *gin.Context) {
	inviteID := c.Param("id")

	if err := ctrl.onboardingService.ExtendInvite(c.Request.Context(), inviteID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "Invitation expiry extended successfully", nil)
}

func (ctrl *InviteController) CancelInvite(c *gin.Context) {
	inviteID := c.Param("id")

	if err := ctrl.onboardingService.CancelInvite(c.Request.Context(), inviteID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMsg(c, "Invitation cancelled successfully", nil)
}
