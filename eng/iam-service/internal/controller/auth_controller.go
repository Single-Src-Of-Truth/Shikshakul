package controller

import (
	"net/http"
	"strings"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/dto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/fingerprint"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(as service.AuthService) *AuthController {
	return &AuthController{authService: as}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "success": false})
		return
	}

	deviceInfo := fingerprint.Extract(c.Request)

	rawToken, user, err := ctrl.authService.Login(c.Request.Context(), req.TenantID, req.Identifier, req.Password, deviceInfo)
	if err != nil {

		if err.Error() == "TENANT_SUSPENDED" {
			response.Error(c, http.StatusPaymentRequired, "Service suspended. Please contact administration to resume access.")
			return
		}

		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("skl_session", rawToken, int(36*3600), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
		"data": gin.H{
			"user_id":    user.ID,
			"first_name": user.FirstName,
			"status":     user.Status,
		},
	})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	token, err := c.Cookie("skl_session")
	if err == nil && token != "" {
		_ = ctrl.authService.Logout(c.Request.Context(), token)
	}

	c.SetCookie("skl_session", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "logged out successfully"})
}

func (ctrl *AuthController) Verify(c *gin.Context) {
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	permsInterface, _ := c.Get("permissions")
	permissions := permsInterface.([]string)

	c.Header("Skl-User-Id", userID)

	if tenantID != "" {
		c.Header("Skl-Tenant-Id", tenantID)
	}

	c.Header("Skl-Permissions", strings.Join(permissions, ","))

	c.Status(http.StatusOK)
}

func (ctrl *AuthController) LogoutAll(c *gin.Context) {
	userID := c.GetString("user_id")
	_ = ctrl.authService.LogoutAll(c.Request.Context(), userID)

	c.SetCookie("skl_session", "", -1, "/", "", true, true)
	response.SuccessWithMsg(c, "Logged out from all devices", nil)
}

func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := ctrl.authService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Password changed successfully", nil)
}

func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	_, _ = ctrl.authService.ForgotPassword(c.Request.Context(), req.TenantID, req.Identifier)
	response.SuccessWithMsg(c, "If the account exists, an OTP has been generated", nil)
}

func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := ctrl.authService.ResetPassword(c.Request.Context(), req.TenantID, req.Identifier, req.OTP, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Password reset successfully. You may now log in.", nil)
}
