package routes

import (
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/controller"
	coreMiddleware "github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/middleware"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/versioning"
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, docCtrl *controller.DocumentController) {
	engine.Use(versioning.EnforceVersion())

	engine.GET("/ping", docCtrl.Ping)
	engine.GET("/", docCtrl.Ping)

	protected := engine.Group("")

	protected.Use(coreMiddleware.RequireGatewayAuth())
	{
		docs := protected.Group("/docs")
		{
			docs.POST("/sign-upload", coreMiddleware.RequireGatewayPermission("utility:files:create"), docCtrl.SignUpload)
			docs.POST("/approve", coreMiddleware.RequireGatewayPermission("utility:files:update"), docCtrl.MoveDocument)
			docs.GET("/access", coreMiddleware.RequireGatewayPermission("utility:files:read"), docCtrl.GetAccessURL)
			docs.POST("/soft-delete", coreMiddleware.RequireGatewayPermission("utility:files:delete"), docCtrl.SoftDelete)
		}
	}
}
