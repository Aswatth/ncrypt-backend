package controllers

import (
	"ncrypt/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemController struct {
	service services.SystemService
}

func (obj *SystemController) Init(service services.SystemService) {
	obj.service = service
}

func (obj *SystemController) GetLoginInfo(ctx *gin.Context) {
	system_data, err := obj.service.GetSystemData()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, system_data)
}

func (obj *SystemController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("system")

	group.GET("/login_info", obj.GetLoginInfo)
}
