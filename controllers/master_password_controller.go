package controllers

import (
	"ncrypt/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MasterPasswordController struct {
	service services.MasterPasswordService
}

func (obj *MasterPasswordController) Init(service *services.MasterPasswordService) {
	obj.service = *service
}

func (obj *MasterPasswordController) SetPassword(ctx *gin.Context) {
	var data map[string]string

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	} else if err := obj.service.SetMasterPassword(data["master_password"]); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *MasterPasswordController) Validate(ctx *gin.Context) {
	var data map[string]string

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	} else {
		result, err := obj.service.ValidateMasterPassword(data["master_password"])

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if !result {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, "invalid password")
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (obj *MasterPasswordController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("master_password")

	group.POST("", obj.SetPassword)
	group.POST("/validate", obj.Validate)
}
