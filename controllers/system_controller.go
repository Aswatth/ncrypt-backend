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

func (obj *SystemController) Login(ctx *gin.Context) {
	request_data := make(map[string]string)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := obj.service.Login(request_data["master_password"]); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) Logout(ctx *gin.Context) {
	if err := obj.service.Logout(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) Export(ctx *gin.Context) {
	request_data := make(map[string]string)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := obj.service.Export(request_data["file_name"], request_data["path"]); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) Import(ctx *gin.Context) {
	request_data := make(map[string]string)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := obj.service.Import(request_data["file_name"], request_data["path"], request_data["master_password"]); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("system")

	group.GET("/login_info", obj.GetLoginInfo)
	group.POST("/login", obj.Login)
	group.POST("/logout", obj.Logout)
	group.POST("/export", obj.Export)
	group.POST("/import", obj.Import)
}
