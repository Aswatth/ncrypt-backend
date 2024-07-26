package controllers

import (
	"ncrypt/models"
	"ncrypt/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DataController struct {
	service services.DataService
}

func (obj *DataController) Init(service *services.DataService) {
	obj.service = *service
}

func (obj *DataController) AddData(ctx *gin.Context) {
	var new_data models.Data

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&new_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err := obj.service.AddData(new_data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *DataController) GetData(ctx *gin.Context) {
	data, err := obj.service.GetAllData()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (obj *DataController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/data")

	group.POST("", obj.AddData)
	group.GET("", obj.GetData)
}
