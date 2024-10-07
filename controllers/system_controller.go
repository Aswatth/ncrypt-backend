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

func (obj *SystemController) Init() {
	obj.service = *new(services.SystemService)
	obj.service.Init()
}

func (obj *SystemController) GetSystemData(ctx *gin.Context) {
	system_data, err := obj.service.GetSystemData()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, system_data)
}

func (obj *SystemController) Setup(ctx *gin.Context) {
	request_data := make(map[string]any)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	err := obj.service.Setup(request_data["master_password"].(string),
		request_data["auto_backup_setting"].(map[string]interface{}))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) SignIn(ctx *gin.Context) {
	request_data := make(map[string]string)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	token, err := obj.service.SignIn(request_data["master_password"])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (obj *SystemController) Logout(ctx *gin.Context) {
	if err := obj.service.Logout(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) Export(ctx *gin.Context) {
	request_data := make(map[string]string)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	if err := obj.service.Export(request_data["file_name"], request_data["path"]); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) Import(ctx *gin.Context) {
	request_data := make(map[string]string)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	if err := obj.service.Import(request_data["file_name"], request_data["path"], request_data["master_password"]); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) GeneratePassword(ctx *gin.Context) {
	password := obj.service.GeneratePassword()

	if password == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Error occured while generating password")
		return
	}

	ctx.JSON(http.StatusOK, password)
}

func (obj *SystemController) Backup(ctx *gin.Context) {
	err := obj.service.Backup()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) UpdateAutomaticBackup(ctx *gin.Context) {
	request_data := make(map[string]any)

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	err := obj.service.UpdateAutomaticBackup(request_data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) GetPasswordGeneratorPreference(ctx *gin.Context) {
	result, err := obj.service.GetPasswordGeneratorPreference()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (obj *SystemController) UpdatePasswordGeneratorPreference(ctx *gin.Context) {
	request_data := make(map[string]interface{})

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	err := obj.service.UpdatePasswordGeneratorPreference(request_data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) UpdateSessionDuration(ctx *gin.Context) {
	request_data := make(map[string]interface{})

	//Check if given JSON is valid
	if err := ctx.ShouldBindJSON(&request_data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	updated_token, err := obj.service.UpdateSessionDuration((int)(request_data["session_duration_in_minutes"].(float64)))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, updated_token)
}

func (obj *SystemController) ExtendSession(ctx *gin.Context) {
	new_token, err := obj.service.ExtendSession()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, new_token)
}

func (obj *SystemController) UpdateTheme(ctx *gin.Context) {
	var theme map[string]string

	if err := ctx.ShouldBindBodyWithJSON(&theme); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}
	err := obj.service.UpdateTheme(theme["theme"])

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		logger.Log.Printf("ERROR: %s", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (obj *SystemController) GetTheme(ctx *gin.Context) {

	theme := obj.service.GetTheme()

	ctx.JSON(http.StatusOK, theme)
}

func (obj *SystemController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("system")

	group.POST("/setup", obj.Setup)
	group.POST("/signin", obj.SignIn)
	group.GET("/generate_password", obj.GeneratePassword)
	group.POST("/import", obj.Import)
	group.GET("/theme", obj.GetTheme)

	group.Use(jwt.ValidateAuthorization())
	group.PUT("/automatic_backup_setting", obj.UpdateAutomaticBackup)

	group.GET("/password_generator_preference", obj.GetPasswordGeneratorPreference)
	group.PUT("/password_generator_preference", obj.UpdatePasswordGeneratorPreference)
	group.PUT("/session_duration", obj.UpdateSessionDuration)
	group.GET("/session_duration", obj.ExtendSession)

	group.PUT("/theme", obj.UpdateTheme)

	group.GET("/data", obj.GetSystemData)
	group.POST("/logout", obj.Logout)
	group.POST("/export", obj.Export)
	group.POST("/backup", obj.Backup)
}
