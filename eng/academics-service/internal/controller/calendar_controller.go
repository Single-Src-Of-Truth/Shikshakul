package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CalendarController struct {
	Service *service.CalendarService
}

func NewCalendarController(svc *service.CalendarService) *CalendarController {
	return &CalendarController{Service: svc}
}

func (ctrl *CalendarController) CreateEvent(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	event, err := ctrl.Service.CreateEvent(tenantID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, "Event created successfully", event)
}

func (ctrl *CalendarController) UpdateEvent(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	eventID := uuid.MustParse(c.Param("id"))

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	event, err := ctrl.Service.UpdateEvent(tenantID, eventID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, event)
}

func (ctrl *CalendarController) GetEvents(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var filter dto.CalendarFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid filters")
		return
	}

	events, err := ctrl.Service.GetMonthEvents(tenantID, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, events)
}

func (ctrl *CalendarController) DeleteEvent(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))

	if err := ctrl.Service.DeleteEvent(tenantID, id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Event deleted", nil)
}
