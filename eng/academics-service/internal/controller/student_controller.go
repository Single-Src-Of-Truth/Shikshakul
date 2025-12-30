package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentController struct {
	Service *service.StudentService
}

func NewStudentController(svc *service.StudentService) *StudentController {
	return &StudentController{Service: svc}
}

func (ctrl *StudentController) CreateSchema(c *gin.Context) {
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
	response.Created(c, "Admission schema created", schema)
}

func (ctrl *StudentController) GetActiveSchema(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	schema, err := ctrl.Service.GetActiveSchema(tenantID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "No active schema found")
		return
	}
	response.Success(c, schema)
}

func (ctrl *StudentController) DeleteSchema(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	if err := ctrl.Service.DeleteSchema(tenantID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Schema deleted", nil)
}

func (ctrl *StudentController) SubmitApplication(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.SubmitApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	student, err := ctrl.Service.SubmitApplication(tenantID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, "Application submitted", student)
}

func (ctrl *StudentController) GetAllStudents(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)

	var filter dto.StudentFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query params")
		return
	}

	students, err := ctrl.Service.GetAllStudents(tenantID, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, students)
}

func (ctrl *StudentController) GetStudentByID(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	student, err := ctrl.Service.GetStudentByID(tenantID, studentID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Student not found")
		return
	}
	response.Success(c, student)
}

func (ctrl *StudentController) UpdateStudent(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	student, err := ctrl.Service.UpdateStudent(tenantID, studentID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, student)
}

func (ctrl *StudentController) DeleteStudent(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := ctrl.Service.DeleteStudent(tenantID, studentID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Student deleted successfully", nil)
}

func (ctrl *StudentController) ApproveStudent(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.ApproveStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	student, err := ctrl.Service.ProcessApproval(tenantID, studentID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	msg := "Student approved successfully"
	if req.Action == "REJECT" {
		msg = "Student application rejected"
	}
	response.SuccessWithMsg(c, msg, student)
}
