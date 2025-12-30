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

type FeeController struct {
	Service *service.FeeService
}

func NewFeeController(svc *service.FeeService) *FeeController {
	return &FeeController{Service: svc}
}

func handleFeeError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrConflict) {
		response.Error(c, http.StatusConflict, err.Error())
	} else if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, err.Error())
	} else {
		response.Error(c, http.StatusInternalServerError, err.Error())
	}
}

func (ctrl *FeeController) CreateFeeHead(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateFeeHeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	head, err := ctrl.Service.CreateFeeHead(tenantID, req)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.Created(c, "Fee Head Created", head)
}

func (ctrl *FeeController) GetAllFeeHeads(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	heads, err := ctrl.Service.GetAllFeeHeads(tenantID)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.Success(c, heads)
}

func (ctrl *FeeController) UpdateFeeHead(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))
	var req dto.UpdateFeeHeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	head, err := ctrl.Service.UpdateFeeHead(tenantID, id, req)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.Success(c, head)
}

func (ctrl *FeeController) DeleteFeeHead(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))
	if err := ctrl.Service.DeleteFeeHead(tenantID, id); err != nil {
		handleFeeError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Fee Head Deleted", nil)
}

func (ctrl *FeeController) CreateFeeStructure(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.CreateFeeStructureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	structure, err := ctrl.Service.CreateFeeStructure(tenantID, req)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.Created(c, "Fee Structure Created", structure)
}

func (ctrl *FeeController) GetStructuresByClass(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	classID := uuid.MustParse(c.Query("class_id"))
	yearID := uuid.MustParse(c.Query("academic_year_id"))

	structures, err := ctrl.Service.GetStructuresByClass(tenantID, classID, yearID)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.Success(c, structures)
}

func (ctrl *FeeController) UpdateFeeStructure(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))
	var req dto.UpdateFeeStructureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	st, err := ctrl.Service.UpdateFeeStructure(tenantID, id, req)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.Success(c, st)
}

func (ctrl *FeeController) DeleteFeeStructure(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	id := uuid.MustParse(c.Param("id"))
	if err := ctrl.Service.DeleteFeeStructure(tenantID, id); err != nil {
		handleFeeError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Fee Structure Deleted", nil)
}

func (ctrl *FeeController) GenerateDemands(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	var req dto.GenerateDemandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	count, err := ctrl.Service.GenerateMonthlyDemands(tenantID, req)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Fee demands generated", gin.H{"count": count})
}

func (ctrl *FeeController) CollectFee(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	userID := c.MustGet("userID").(uuid.UUID)
	var req dto.CollectFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	receipt, err := ctrl.Service.CollectFee(tenantID, userID, req)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.SuccessWithMsg(c, "Payment collected", receipt)
}

func (ctrl *FeeController) GetStudentDues(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	studentID := uuid.MustParse(c.Param("student_id"))

	dues, err := ctrl.Service.GetStudentDues(tenantID, studentID)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.Success(c, dues)
}

func (ctrl *FeeController) GetStudentHistory(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	studentID := uuid.MustParse(c.Param("student_id"))

	txns, err := ctrl.Service.GetStudentHistory(tenantID, studentID)
	if err != nil {
		handleFeeError(c, err)
		return
	}
	response.Success(c, txns)
}
