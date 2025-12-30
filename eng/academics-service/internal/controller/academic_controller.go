package controller

import (
	"errors"
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AcademicController struct {
	Service *service.AcademicService
}

func NewAcademicController(svc *service.AcademicService) *AcademicController {
	return &AcademicController{Service: svc}
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrConflict) {
		response.Error(c, http.StatusConflict, err.Error())
	} else if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, err.Error())
	} else {
		response.Error(c, http.StatusInternalServerError, err.Error())
	}
}

func (ctrl *AcademicController) CreateYear(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	year, err := ctrl.Service.CreateAcademicYear(tenantID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, "Academic year created", year)
}

func (ctrl *AcademicController) UpdateYear(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))

	var req dto.UpdateYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	year, err := ctrl.Service.UpdateAcademicYear(tenantID, id, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, year)
}

func (ctrl *AcademicController) DeleteYear(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))

	if err := ctrl.Service.DeleteAcademicYear(tenantID, id); err != nil {
		handleError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Academic year deleted", nil)
}

func (ctrl *AcademicController) GetAllYears(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	years, err := ctrl.Service.GetAllYears(tenantID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, years)
}

func (ctrl *AcademicController) CreateClass(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	class, err := ctrl.Service.CreateClass(tenantID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, "Class created", class)
}

func (ctrl *AcademicController) UpdateClass(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("class_id"))

	var req dto.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	class, err := ctrl.Service.UpdateClass(tenantID, id, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, class)
}

func (ctrl *AcademicController) DeleteClass(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("class_id"))

	if err := ctrl.Service.DeleteClass(tenantID, id); err != nil {
		handleError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Class deleted", nil)
}

func (ctrl *AcademicController) GetAllClasses(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	classes, err := ctrl.Service.GetAllClasses(tenantID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, classes)
}

func (ctrl *AcademicController) CreateSection(c *gin.Context) {
	classID := c.Param("class_id")
	var req dto.CreateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	section, err := ctrl.Service.CreateSection(classID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, "Section created", section)
}

func (ctrl *AcademicController) UpdateSection(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	var req dto.UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	section, err := ctrl.Service.UpdateSection(id, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, section)
}

func (ctrl *AcademicController) DeleteSection(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	if err := ctrl.Service.DeleteSection(id); err != nil {
		handleError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Section deleted", nil)
}

func (ctrl *AcademicController) GetSectionsByClass(c *gin.Context) {
	classID := uuid.MustParse(c.Param("class_id"))
	sections, err := ctrl.Service.GetSectionsByClass(classID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, sections)
}

func (ctrl *AcademicController) CreateSubject(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	subject, err := ctrl.Service.CreateSubject(tenantID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, "Subject created", subject)
}

func (ctrl *AcademicController) UpdateSubject(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))
	var req dto.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	subject, err := ctrl.Service.UpdateSubject(tenantID, id, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, subject)
}

func (ctrl *AcademicController) DeleteSubject(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))
	if err := ctrl.Service.DeleteSubject(tenantID, id); err != nil {
		handleError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Subject deleted", nil)
}

func (ctrl *AcademicController) GetAllSubjects(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	subjects, err := ctrl.Service.GetAllSubjects(tenantID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, subjects)
}

func (ctrl *AcademicController) AssignSubjectsToClass(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	classID := c.Param("class_id")
	var req dto.AssignSubjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := ctrl.Service.AssignSubjectsToClass(tenantID, classID, req); err != nil {
		handleError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Subjects assigned", nil)
}

func (ctrl *AcademicController) UnassignSubject(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	classID := uuid.MustParse(c.Param("class_id"))
	subjectID := uuid.MustParse(c.Param("subject_id"))

	if err := ctrl.Service.UnassignSubject(tenantID, classID, subjectID); err != nil {
		handleError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Subject unassigned from class", nil)
}

func (ctrl *AcademicController) GetCurrentYear(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	year, err := ctrl.Service.GetCurrentYear(tenantID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, year)
}
