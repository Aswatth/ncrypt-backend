package controllers

import (
	"fmt"
	"ncrypt/models"
	"ncrypt/services"
	"ncrypt/utils/jwt"
	"ncrypt/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginDataController struct {
	service services.ILoginDataService
}

func (obj *LoginDataController) Init() {
	obj.service = services.InitBadgerLoginService()
	obj.service.Init()
}

func (obj *LoginDataController) AddLoginData(ctx *gin.Context) {
	var new_login_data models.Login

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&new_login_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	if err := obj.service.AddLoginData(&new_login_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *LoginDataController) GetLoginData(ctx *gin.Context) {
	name := ctx.Query("name")

	//Get all
	if name == "" {
		if data, err := obj.service.GetAllLoginData(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			logger.Log.Printf("ERROR: %s", err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, data)
		}
	} else if data, err := obj.service.GetLoginData(name); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	} else {
		ctx.JSON(http.StatusOK, data)
	}
}

func (obj *LoginDataController) GetAccountPassword(ctx *gin.Context) {
	login_data_name := ctx.Param("name")
	account_username := ctx.Query("username")

	fmt.Println(login_data_name, account_username)

	password, err := obj.service.GetDecryptedAccountPassword(login_data_name, account_username)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, password)
}

func (obj *LoginDataController) DeleteLoginData(ctx *gin.Context) {
	name := ctx.Param("name")

	if err := obj.service.DeleteLoginData(name); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	} else {
		ctx.Status(http.StatusOK)
	}
}

func (obj *LoginDataController) UpdateLoginData(ctx *gin.Context) {
	name := ctx.Param("name")

	var login_data models.Login

	if err := ctx.ShouldBindJSON(&login_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	if err := obj.service.UpdateLoginData(name, &login_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	} else {
		ctx.Status(http.StatusOK)
	}
}

func (obj *LoginDataController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/login")

	group.Use(jwt.ValidateAuthorization())
	group.POST("", obj.AddLoginData)

	group.GET("", obj.GetLoginData)
	group.GET("/:name", obj.GetAccountPassword)

	group.DELETE("/:name", obj.DeleteLoginData)
	group.PUT("/:name", obj.UpdateLoginData)
}
