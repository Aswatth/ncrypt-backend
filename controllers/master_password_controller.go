package controllers

import (
	"ncrypt/services"
	"ncrypt/utils/jwt"
	"ncrypt/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MasterPasswordController struct {
	service services.IMasterPasswordService
}

func (obj *MasterPasswordController) Init() {
	logger.Log.Printf("Initializing master password controller")
	obj.service = services.InitBadgerMasterPasswordService()
	obj.service.Init()
	logger.Log.Printf("Initialization complete!")
}

func (obj *MasterPasswordController) SetPassword(ctx *gin.Context) {
	var data map[string]string

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	} else if err := obj.service.SetMasterPassword(data["master_password"]); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *MasterPasswordController) UpdatePassword(ctx *gin.Context) {
	var data map[string]string

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	} else if err := obj.service.UpdateMasterPassword(data["old_master_password"], data["new_master_password"]); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *MasterPasswordController) ValidatePassword(ctx *gin.Context) {
	var data map[string]string

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	result, err := obj.service.Validate(data["master_password"])

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}
	if !result {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "incorrect password")
		logger.Log.Printf("ERROR: %s", "incorrect password")
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *MasterPasswordController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("master_password")

	group.Use(jwt.ValidateAuthorization())
	group.POST("/validate", obj.ValidatePassword)
	group.POST("", obj.SetPassword)
	group.PUT("", obj.UpdatePassword)
}
