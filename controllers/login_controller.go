package controllers

import (
	"ncrypt/models"
	"ncrypt/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	service services.LoginService
}

func (obj *LoginController) Init(serivce *services.LoginService) {
	obj.service = *serivce
}

func (obj *LoginController) CreateLogin(ctx *gin.Context) {
	var new_login_data models.Login

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&new_login_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := obj.service.AddLoginData(&new_login_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *LoginController) GetLoginData(ctx *gin.Context) {
	name := ctx.Query("name")

	//Get all
	if name == "" {
		if data, err := obj.service.GetAllLoginData(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, data)
		}
	} else if data, err := obj.service.GetLoginData(name); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	} else {
		ctx.JSON(http.StatusOK, data)
	}
}

func (obj *LoginController) DeleteLoginData(ctx *gin.Context) {
	name := ctx.Param("name")

	if err := obj.service.DeleteLogin(name); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	} else {
		ctx.Status(http.StatusOK)
	}
}

func (obj *LoginController) UpdateLoginData(ctx *gin.Context) {
	name := ctx.Param("name")

	var login_data models.Login

	if err := ctx.ShouldBindJSON(&login_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := obj.service.UpdateLoginData(name, &login_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	} else {
		ctx.Status(http.StatusOK)
	}
}

func (obj *LoginController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/login")

	group.POST("", obj.CreateLogin)
	group.GET("", obj.GetLoginData)
	group.DELETE("/:name", obj.DeleteLoginData)
	group.PUT("/:name", obj.UpdateLoginData)
}