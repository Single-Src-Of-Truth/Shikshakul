package controller

import (
	"net/http"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/dto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/response"
	"github.com/gin-gonic/gin"
)

type TenantController struct {
	tenantService service.TenantService
}

func NewTenantController(ts service.TenantService) *TenantController {
	return &TenantController{tenantService: ts}
}

func (ctrl *TenantController) CreateTenant(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload data")
		return
	}

	tenantResp, err := ctrl.tenantService.CreateTenant(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create tenant")
		return
	}

	response.Created(c, "Tenant created successfully", tenantResp)
}

func (ctrl *TenantController) ListTenants(c *gin.Context) {
	tenants, err := ctrl.tenantService.ListTenants(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch tenants")
		return
	}

	response.Success(c, tenants)
}

func (ctrl *TenantController) GetTenant(c *gin.Context) {
	tenant, err := ctrl.tenantService.GetTenant(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Tenant not found")
		return
	}
	response.Success(c, tenant)
}

func (ctrl *TenantController) UpdateTenant(c *gin.Context) {
	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload")
		return
	}
	if err := ctrl.tenantService.UpdateTenant(c.Request.Context(), c.Param("id"), req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update tenant")
		return
	}
	response.SuccessWithMsg(c, "Tenant updated successfully", nil)
}

func (ctrl *TenantController) UpdateUserStatus(c *gin.Context) {
	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid status payload. Must be 'ACTIVE' or 'SUSPENDED'")
		return
	}

	targetUserID := c.Param("userId")

	if err := ctrl.tenantService.UpdateUserStatus(c.Request.Context(), targetUserID, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMsg(c, "User status successfully updated to "+req.Status, nil)
}

func (ctrl *TenantController) DeleteUser(c *gin.Context) {
	targetUserID := c.Param("userId")

	if err := ctrl.tenantService.DeleteUser(c.Request.Context(), targetUserID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMsg(c, "User account successfully deleted", nil)
}

func (ctrl *TenantController) ChangeUserRole(c *gin.Context) {
	targetUserID := c.Param("userId")

	var req dto.ChangeUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload. 'role_id' is required and must be a valid UUID")
		return
	}

	if err := ctrl.tenantService.ChangeUserRole(c.Request.Context(), targetUserID, req.RoleID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMsg(c, "User role successfully updated. The user has been logged out of all active sessions to apply the new permissions.", nil)
}

func (ctrl *TenantController) GetPublicTenants(c *gin.Context) {
	tenants, err := ctrl.tenantService.GetPublicTenants(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, tenants)
}
