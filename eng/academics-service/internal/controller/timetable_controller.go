package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TimetableController struct {
	Service *service.TimetableService
}

func NewTimetableController(svc *service.TimetableService) *TimetableController {
	return &TimetableController{Service: svc}
}

func (ctrl *TimetableController) CreateRoutine(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateRoutineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	routine, err := ctrl.Service.CreateRoutine(tenantID, req)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}
	response.Created(c, "Routine added successfully", routine)
}

func (ctrl *TimetableController) GetTimetable(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var filter dto.TimetableFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid filters")
		return
	}

	routines, err := ctrl.Service.GetTimetable(tenantID, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, routines)
}

func (ctrl *TimetableController) DeleteRoutine(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))

	if err := ctrl.Service.DeleteRoutine(tenantID, id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Routine deleted", nil)
}
