package controller

import (
	"net/http"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/dto"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CertificateController struct {
	Service *service.CertificateService
}

func NewCertificateController(svc *service.CertificateService) *CertificateController {
	return &CertificateController{Service: svc}
}

func (ctrl *CertificateController) GetIDCards(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	classID := c.Query("class_id")
	if classID == "" {
		response.Error(c, http.StatusBadRequest, "class_id is required")
		return
	}

	data, err := ctrl.Service.GetIDCardData(tenantID, classID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, data)
}

func (ctrl *CertificateController) GenerateTC(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	issuerID := c.MustGet("userID").(uuid.UUID)

	var req dto.IssueTCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := ctrl.Service.GenerateTC(tenantID, issuerID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Transfer Certificate Generated", data)
}

func (ctrl *CertificateController) GenerateBonafide(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(uuid.UUID)
	issuerID := c.MustGet("userID").(uuid.UUID)

	var req dto.IssueBonafideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := ctrl.Service.GenerateBonafide(tenantID, issuerID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMsg(c, "Bonafide Certificate Generated", data)
}
