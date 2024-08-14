package controllers

import (
	"ncrypt/services"
	"ncrypt/utils/jwt"
	"ncrypt/utils/logger"
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

	token, err := obj.service.Login(request_data["master_password"])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, token)

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
		logger.Log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := obj.service.Export(request_data["file_name"], request_data["path"]); err != nil {
		logger.Log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) Import(ctx *gin.Context) {
	request_data := make(map[string]string)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		logger.Log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := obj.service.Import(request_data["file_name"], request_data["path"], request_data["master_password"]); err != nil {
		logger.Log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) GeneratePassword(ctx *gin.Context) {
	request_data := make(map[string]interface{})

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	password := obj.service.GeneratePassword(request_data["has_digits"].(bool), request_data["has_upper_case"].(bool), request_data["has_special_char"].(bool), int(request_data["length"].(float64)))

	ctx.JSON(http.StatusOK, password)
}

func (obj *SystemController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("system")

	group.GET("/login_info", obj.GetLoginInfo)
	group.POST("/login", obj.Login)
	group.GET("/generate_password", obj.GeneratePassword)
	group.POST("/import", obj.Import)

	group.Use(jwt.ValidateAuthorization())
	group.POST("/logout", obj.Logout)
	group.POST("/export", obj.Export)
}
