package controllers

import (
	"ncrypt/services"
	"ncrypt/utils/jwt"
	"ncrypt/utils/logger"
	"net/http"
	"strconv"

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

func (obj *SystemController) Setup(ctx *gin.Context) {
	request_data := make(map[string]any)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err := obj.service.Setup(request_data["master_password"].(string),
		request_data["automatic_backup"].(bool),
		request_data["backup_folder_path"].(string),
		request_data["backup_file_name"].(string))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
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
	has_digits := ctx.Query("hasDigits") == "true"
	has_upper_case := ctx.Query("hasUpperCase") == "true"
	has_special_char := ctx.Query("hasSpecialChar") == "true"
	length := 8

	if ctx.Query("length") != "" {
		l, err := strconv.Atoi(ctx.Query("length"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		length = l
	}

	password := obj.service.GeneratePassword(has_digits, has_upper_case, has_special_char, length)

	ctx.JSON(http.StatusOK, password)
}

func (obj *SystemController) Backup(ctx *gin.Context) {
	err := obj.service.Backup()

	if err != nil {
		logger.Log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) UpdateAutomaticBackup(ctx *gin.Context) {
	request_data := make(map[string]any)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err := obj.service.UpdateAutomaticBackup(
		request_data["automatic_backup"].(bool),
		request_data["backup_folder_path"].(string),
		request_data["backup_file_name"].(string))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("system")

	group.POST("/setup", obj.Setup)
	group.POST("/login", obj.Login)
	group.GET("/generate_password", obj.GeneratePassword)
	group.POST("/import", obj.Import)

	group.Use(jwt.ValidateAuthorization())
	group.PUT("/automatic_backup_setting", obj.UpdateAutomaticBackup)
	group.GET("/login_info", obj.GetLoginInfo)
	group.POST("/logout", obj.Logout)
	group.POST("/export", obj.Export)
	group.POST("/backup", obj.Backup)
}
