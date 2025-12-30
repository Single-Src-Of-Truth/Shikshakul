package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeacherController struct {
	Service *service.TeacherService
}

func NewTeacherController(svc *service.TeacherService) *TeacherController {
	return &TeacherController{Service: svc}
}

func (ctrl *TeacherController) CreateSchema(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid JSON structure")
		return
	}
	schema, err := ctrl.Service.CreateSchema(tenantID, input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(c, "Teacher schema created", schema)
}

func (ctrl *TeacherController) GetActiveSchema(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	schema, err := ctrl.Service.GetActiveSchema(tenantID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "No active schema found")
		return
	}
	response.Success(c, schema)
}

func (ctrl *TeacherController) DeleteSchema(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	if err := ctrl.Service.DeleteSchema(tenantID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Schema deleted", nil)
}

func (ctrl *TeacherController) OnboardTeacher(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.OnboardTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	teacher, err := ctrl.Service.OnboardTeacher(tenantID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, "Teacher onboarding submitted", teacher)
}

func (ctrl *TeacherController) UpdateTeacher(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))

	var req dto.UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	teacher, err := ctrl.Service.UpdateTeacher(tenantID, id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, teacher)
}

func (ctrl *TeacherController) DeleteTeacher(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))

	if err := ctrl.Service.DeleteTeacher(tenantID, id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Teacher deleted", nil)
}

func (ctrl *TeacherController) GetTeacherByID(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))

	teacher, err := ctrl.Service.GetTeacherByID(tenantID, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Teacher not found")
		return
	}
	response.Success(c, teacher)
}

func (ctrl *TeacherController) GetAllTeachers(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var filter dto.TeacherFilter
	c.ShouldBindQuery(&filter)

	teachers, err := ctrl.Service.GetAllTeachers(tenantID, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, teachers)
}

func (ctrl *TeacherController) ApproveTeacher(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	teacherID := uuid.MustParse(c.Param("id"))
	var req dto.ApproveTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	teacher, err := ctrl.Service.ApproveTeacher(tenantID, teacherID, req.Action)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, teacher)
}

func (ctrl *TeacherController) AssignClassTeacher(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	sectionID := uuid.MustParse(c.Param("section_id"))
	var req dto.AssignClassTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	teacherID := uuid.MustParse(req.TeacherID)

	err := ctrl.Service.AssignClassTeacher(tenantID, sectionID, teacherID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Class teacher assigned", nil)
}

func (ctrl *TeacherController) AssignSubjectTeacher(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	sectionID := uuid.MustParse(c.Param("section_id"))

	var req dto.AssignSubjectTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := ctrl.Service.AssignSubjectTeacher(tenantID, sectionID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Subject teacher allocated", nil)
}

func (ctrl *TeacherController) GetSectionAllocations(c *gin.Context) {
	sectionID := uuid.MustParse(c.Param("section_id"))
	yearID := uuid.MustParse(c.Query("academic_year_id"))

	allocs, err := ctrl.Service.GetAllocations(sectionID, yearID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, allocs)
}
