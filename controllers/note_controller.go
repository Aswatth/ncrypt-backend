package controllers

import (
	"ncrypt/models"
	"ncrypt/services"
	"ncrypt/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	service services.INoteService
}

func (obj *NoteController) Init() {
	obj.service = services.InitBadgerNoteService()
	obj.service.Init()
}

func (obj *NoteController) AddNote(ctx *gin.Context) {
	var new_note models.Note

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&new_note); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	if err := obj.service.AddNote(&new_note); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *NoteController) GetNote(ctx *gin.Context) {
	created_date_time := ctx.Query("created_date_time")

	if(created_date_time == "") {
		if data, err := obj.service.GetAllNotes(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			logger.Log.Printf("ERROR: %s", err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, data)
		}	
	} else {
		if data, err := obj.service.GetNote(created_date_time); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			logger.Log.Printf("ERROR: %s", err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, data)
		}
	}
}

func (obj *NoteController) GetContent(ctx *gin.Context) {
	created_date_time := ctx.Param("created_date_time")

	if data, err := obj.service.GetDecryptedContent(created_date_time); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	} else {
		ctx.JSON(http.StatusOK, data)
	}
}

func (obj *NoteController) DeleteNote(ctx *gin.Context) {
	created_date_time := ctx.Param("created_date_time")

	if err := obj.service.DeleteNote(created_date_time); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	} else {
		ctx.Status(http.StatusOK)
	}
}

func (obj *NoteController) UpdateNote(ctx *gin.Context) {
	created_date_time := ctx.Param("created_date_time")

	var new_note models.Note

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&new_note); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	if err := obj.service.UpdateNote(created_date_time, new_note); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *NoteController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/note")

	group.POST("", obj.AddNote)
	group.GET("", obj.GetNote)
	group.GET("/:created_date_time", obj.GetContent)
	group.DELETE("/:created_date_time", obj.DeleteNote)
	group.PUT("/:created_date_time", obj.UpdateNote)
}
