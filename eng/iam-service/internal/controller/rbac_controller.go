package controller

import (
	"net/http"
	"strings"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/dto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/response"
	"github.com/gin-gonic/gin"
)

type RBACController struct {
	rbacService service.RBACService
}

func NewRBACController(rs service.RBACService) *RBACController {
	return &RBACController{rbacService: rs}
}

func (ctrl *RBACController) ListPermissions(c *gin.Context) {
	permissions, err := ctrl.rbacService.GetAllPermissions(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch permissions")
		return
	}

	var responseData []map[string]string
	for _, p := range permissions {
		responseData = append(responseData, map[string]string{
			"id":          p.ID,
			"description": p.Description,
		})
	}

	response.Success(c, responseData)
}

func (ctrl *RBACController) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role payload")
		return
	}

	tenantID := c.GetString("tenant_id")
	var tID *string
	if tenantID != "" {
		tID = &tenantID
	}

	role, err := ctrl.rbacService.CreateCustomRole(c.Request.Context(), tID, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "unique constraint") {
			response.Error(c, http.StatusConflict, "A role with this name already exists in your school")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to create custom role")
		return
	}

	response.Created(c, "Role created successfully", dto.RoleResponse{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		Permissions: req.PermissionIDs,
	})
}

func (ctrl *RBACController) BulkCreatePermissions(c *gin.Context) {
	var req dto.BulkCreatePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload")
		return
	}
	if err := ctrl.rbacService.BulkSeedPermissions(c.Request.Context(), req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to seed permissions")
		return
	}
	response.SuccessWithMsg(c, "Permissions seeded successfully", nil)
}

func (ctrl *RBACController) UpdateRole(c *gin.Context) {
	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload")
		return
	}
	if err := ctrl.rbacService.UpdateCustomRole(c.Request.Context(), c.Param("id"), req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Role updated successfully", nil)
}

func (ctrl *RBACController) DeleteRole(c *gin.Context) {
	if err := ctrl.rbacService.DeleteCustomRole(c.Request.Context(), c.Param("id")); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "conflict: cannot delete role while users are assigned to it" {
			status = http.StatusConflict
		}
		response.Error(c, status, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Role deleted successfully", nil)
}

func (ctrl *RBACController) ListRoles(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var tID *string
	if tenantID != "" {
		tID = &tenantID
	}

	roles, err := ctrl.rbacService.ListRoles(c.Request.Context(), tID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch roles")
		return
	}
	response.Success(c, roles)
}
