package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"ncrypt/models"
	"ncrypt/utils"
	"ncrypt/utils/database"
	"ncrypt/utils/encryptor"
	"ncrypt/utils/jwt"
	"ncrypt/utils/logger"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/dgraph-io/badger/v4"
	"github.com/joho/godotenv"
)

type SystemService struct {
	database                    database.IDatabase
	database_name               string
	master_password_service     IMasterPasswordService
	SESSION_DURATION_IN_MINUTES int
}

func (obj *SystemService) Init() {

	logger.Log.Printf("Initializing system service")
	logger.Log.Printf("Setting up database")
	godotenv.Load("../.env")

	obj.database = database.InitBadgerDb()
	obj.database_name = "SYSTEM"
	obj.database.SetDatabase(obj.database_name)

	logger.Log.Printf("Setting up master password service")
	obj.master_password_service = InitBadgerMasterPasswordService()
	obj.master_password_service.Init()

	obj.SESSION_DURATION_IN_MINUTES = 20 //20 minutes is default
	logger.Log.Printf("System service initialized")

	// Code to launch UI - comment these lines to prevent launching of multiple UI instances while testing.
	system_data, err := obj.GetSystemData()

	isNewUser := "false"
	theme := "SYSTEM"

	if err != nil {
		if err == badger.ErrKeyNotFound {
			isNewUser = "true"
		}
	} else {
		theme = system_data.Theme
	}

	go obj.launchUI(os.Getenv("UI_EXECUTABLE_PATH"), []string{utils.PORT, isNewUser, theme})
	restoreAndBringWindowToFront("NCRYPT")
}

func (obj *SystemService) launchUI(commandPath string, args []string) {
	cmd := exec.Command(commandPath, args...)
	// Run the command and wait for it to complete
	err := cmd.Run()
	if err != nil {
		// Handle error
		return
	}
	obj.Logout()
	os.Exit(0)
}

func restoreAndBringWindowToFront(title string) error {
	user32 := syscall.NewLazyDLL("user32.dll")

	findWindow := user32.NewProc("FindWindowW")
	setForegroundWindow := user32.NewProc("SetForegroundWindow")

	showWindow := user32.NewProc("ShowWindow")
	const SW_RESTORE = 9

	// Convert string to UTF16
	u16Title, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}

	// Find the window by its title
	hwnd, _, err := findWindow.Call(0, uintptr(unsafe.Pointer(u16Title)))
	if hwnd == 0 {
		logger.Log.Printf("ERROR: window not found: %v", err.Error())
		return err
	}

	// Restore the window if minimized
	showWindow.Call(hwnd, SW_RESTORE)

	// Bring the window to the foreground
	_, _, err = setForegroundWindow.Call(hwnd)
	if err != nil {
		return err
	}

	return nil
}

func (obj *SystemService) initSystem(system_data models.SystemData) error {
	_, err := obj.GetSystemData()

	if err != nil && err == badger.ErrKeyNotFound {
		err = obj.setSystemData(system_data)
		if err != nil {
			logger.Log.Printf("ERROR: %s", err.Error())
			return err
		}
		err = nil
	}

	return err
}

func (obj *SystemService) setSystemData(system_data models.SystemData) error {
	logger.Log.Printf("Setting system data")
	err := obj.database.AddData(obj.database_name, system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
	}

	return err
}

func (obj *SystemService) GetSystemData() (*models.SystemData, error) {
	logger.Log.Printf("Getting system data")
	fetched_data, err := obj.database.GetData(obj.database_name)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return nil, err
	}

	var system_data models.SystemData
	system_data.FromMap(fetched_data.(map[string]interface{}))

	return &system_data, err
}

func (obj *SystemService) Setup(master_password string, auto_backup_setting map[string]interface{}) error {

	_, err := obj.master_password_service.GetMasterPassword()

	if err != nil {
		if err != badger.ErrKeyNotFound {
			logger.Log.Printf("ERROR: %s", err.Error())
			return err
		}
	}

	logger.Log.Printf("Setting up master password")
	err = obj.master_password_service.SetMasterPassword(master_password)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Printf("Setting up system data")

	new_auto_backup_setting := new(models.AutoBackupSetting).FromMap(auto_backup_setting)

	password_generator_preferance := new(models.PasswordGeneratorPreference)
	password_generator_preferance.HasDigits = false
	password_generator_preferance.HasUpperCase = false
	password_generator_preferance.HasSpecialChar = false
	password_generator_preferance.Length = 8

	err = obj.initSystem(models.SystemData{LoginCount: 0, LastLoginDateTime: "", CurrentLoginDateTime: "", IsLoggedIn: false, PasswordGeneratorPreference: *password_generator_preferance, AutoBackupSetting: *new_auto_backup_setting, SessionDurationInMinutes: obj.SESSION_DURATION_IN_MINUTES, Theme: "SYSTEM"})
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Printf("Initial setup completed")
	return nil
}

func (obj *SystemService) SignIn(password string) (string, error) {
	logger.Log.Printf("Logging in")
	result, err := obj.master_password_service.Validate(password)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	if !result {
		logger.Log.Printf("ERROR: invalid password")
		return "", errors.New("invalid password")
	}

	system_data, err := obj.GetSystemData()
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	system_data.IsLoggedIn = true
	system_data.LoginCount += 1
	system_data.CurrentLoginDateTime = time.Now().Format(time.RFC3339)

	err = obj.setSystemData(*system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	logger.Log.Printf("Logged in")

	token, err := jwt.GenerateToken(system_data.SessionDurationInMinutes)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return "", err
	}

	return token, nil
}

func (obj *SystemService) Logout() error {
	logger.Log.Printf("Logging out")

	system_data, err := obj.GetSystemData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	if !system_data.IsLoggedIn {
		my_error := errors.New("not logged in to logout")
		logger.Log.Printf("ERROR: %s", my_error.Error())
		return my_error
	}

	system_data.IsLoggedIn = false
	system_data.LastLoginDateTime = system_data.CurrentLoginDateTime

	err = obj.setSystemData(*system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Printf("Logged out")
	return err
}

func (obj *SystemService) Export(file_name string, file_path string) error {
	logger.Log.Println("Exporting data...")

	if !strings.HasSuffix(file_name, ".ncrypt") {
		return errors.New("incorrect file format")
	}

	export_data := new(ExportData)

	logger.Log.Println("Fetching master password")
	//Get master password data
	master_password, err := obj.master_password_service.GetMasterPassword()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	} else {
		export_data.MASTER_PASSWORD = master_password
	}

	system_data, err := obj.GetSystemData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	} else {
		export_data.SYSTEM_DATA = *system_data
	}

	var wg sync.WaitGroup
	err_channel := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Log.Println("Fetching login data")

		//Get login data
		login_service := InitBadgerLoginService()
		login_service.Init()
		login_data_list, err := login_service.GetAllLoginData()

		if err != nil {
			logger.Log.Printf("ERROR: %s", err.Error())
			err_channel <- err
		} else {
			export_data.LOGIN_DATA = login_data_list
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Log.Println("Fetching notes")

		note_service := InitBadgerNoteService()
		note_service.Init()
		notes, err := note_service.GetAllNotes()

		if err != nil {
			logger.Log.Printf("ERROR: %s", err.Error())
			err_channel <- err
		} else {
			export_data.NOTE_DATA = notes
		}
	}()

	go func() {
		wg.Wait()
		close(err_channel)
	}()

	for err := range err_channel {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Exporting to " + file_path + "\\" + file_name)
	var path string
	if file_path != "" {
		path = file_path + "\\" + file_name
	} else {
		path = file_name
	}

	file, err := os.Create(path)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	defer file.Close()

	export_data_bytes, err := json.Marshal(export_data)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Encrypting export data")
	//Encrpyt data using master_password
	encrypted_export_data, err := encryptor.Encrypt(base64.StdEncoding.EncodeToString(export_data_bytes), master_password)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	encrypted_export_data_bytes, err := base64.StdEncoding.DecodeString(encrypted_export_data)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Saving to file")
	_, err = file.Write(encrypted_export_data_bytes)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Printf("Export complete!")
	return nil
}

func (obj *SystemService) Import(file_name string, file_path string, master_password string) error {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))

	logger.Log.Println("Importing data")
	file, err := os.Open(file_path + "\\" + file_name)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}
	defer file.Close()

	logger.Log.Println("Reading import file")
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Println("Decrypting import content")
	decrypted_data, err := encryptor.Decrypt(base64.StdEncoding.EncodeToString(data), encryptor.CreateHash(master_password))
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	decrypted_data_bytes, err := base64.StdEncoding.DecodeString(decrypted_data)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return errors.New("incorrect master password or corrupted file")
	}

	logger.Log.Println("Importing system data")
	imported_data := new(ExportData)
	json.Unmarshal(decrypted_data_bytes, &imported_data)

	//Import system data
	logger.Log.Println("Importing system data")
	err = obj.setSystemData(imported_data.SYSTEM_DATA)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	//Import master password
	logger.Log.Println("Importing master password")
	err = obj.master_password_service.importData(imported_data.MASTER_PASSWORD)
	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	var wg sync.WaitGroup
	err_channel := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		//Import login data
		logger.Log.Println("Importing login data")
		login_service := InitBadgerLoginService()
		login_service.Init()
		err = login_service.importData(imported_data.LOGIN_DATA)
		if err != nil {
			logger.Log.Printf("ERROR: %s", err.Error())
			// return err
			err_channel <- err
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		//Import note data
		logger.Log.Println("Importing note data")
		note_service := InitBadgerNoteService()
		note_service.Init()
		err = note_service.importData(imported_data.NOTE_DATA)
		if err != nil {
			logger.Log.Printf("ERROR: %s", err.Error())
			// return err
			err_channel <- err
		}
	}()

	go func() {
		wg.Wait()
		close(err_channel)
	}()

	for err := range err_channel {
		if err != nil {
			return err
		}
	}

	logger.Log.Println("DONE")
	return nil
}

type ExportData struct {
	SYSTEM_DATA     models.SystemData `json:"SYSTEM" bson:"SYSTEM"`
	LOGIN_DATA      []models.Login    `json:"LOGIN_DATA" bson:"LOGIN_DATA"`
	NOTE_DATA       []models.Note     `json:"NOTE_DATA" bson:"NOTE_DATA"`
	MASTER_PASSWORD string            `json:"MASTER_PASSWORD" bson:"MASTER_PASSWORD"`
}

func (obj *SystemService) GeneratePassword() string {
	password_preference, err := obj.GetPasswordGeneratorPreference()

	if err != nil {
		return ""
	}

	return utils.GeneratePassword(password_preference.HasDigits, password_preference.HasUpperCase, password_preference.HasSpecialChar, password_preference.Length)
}

func (obj *SystemService) Backup() error {
	logger.Log.Printf("Backing up data")
	logger.Log.Printf("Getting system data")
	system_data, err := obj.GetSystemData()

	logger.Log.Printf("Checking for automatic backup setting")
	auto_backup_setting := system_data.AutoBackupSetting

	if err == nil {
		if auto_backup_setting.IsEnabled {
			logger.Log.Printf("Automatic backup is enabled")
			file_name := auto_backup_setting.BackupFileName + "_" + time.Now().Format(time.RFC3339) + ".ncrypt"
			file_name = strings.Replace(file_name, ":", "-", -1)
			logger.Log.Printf("Exporting data")
			err = obj.Export(file_name, auto_backup_setting.BackupLocation)
			if err != nil {
				logger.Log.Printf("ERROR: %s", err.Error())
			}
		} else {
			logger.Log.Printf("Automatic backup is not enabled")
		}
	} else {
		logger.Log.Printf("ERROR: %s", err.Error())
	}

	return err
}

func (obj *SystemService) UpdateAutomaticBackup(updated_auto_backup_setting map[string]interface{}) error {
	logger.Log.Printf("Updating automatic backup data")

	system_data, err := obj.GetSystemData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	auto_backup_setting := new(models.AutoBackupSetting).FromMap(updated_auto_backup_setting)

	if auto_backup_setting.IsEnabled {
		if auto_backup_setting.BackupFileName == "" {
			return errors.New("file name cannot be empty")
		}
	}

	system_data.AutoBackupSetting = *auto_backup_setting

	err = obj.setSystemData(*system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	logger.Log.Printf("Update completed successfully")
	return nil
}

func (obj *SystemService) GetPasswordGeneratorPreference() (*models.PasswordGeneratorPreference, error) {
	system_data, err := obj.GetSystemData()

	if err != nil {
		return nil, err
	}

	return &system_data.PasswordGeneratorPreference, err
}

func (obj *SystemService) UpdatePasswordGeneratorPreference(data map[string]interface{}) error {
	system_data, err := obj.GetSystemData()

	if err != nil {
		return err
	}

	password_generator_preference := new(models.PasswordGeneratorPreference)
	password_generator_preference.FromMap(data)

	system_data.PasswordGeneratorPreference = *password_generator_preference
	err = obj.setSystemData(*system_data)

	if err != nil {
		return err
	}

	return nil
}

func (obj *SystemService) UpdateSessionDuration(session_duration_in_minutes int) (string, error) {
	system_data, err := obj.GetSystemData()

	if err != nil {
		return "", err
	}

	system_data.SessionDurationInMinutes = session_duration_in_minutes

	err = obj.setSystemData(*system_data)

	if err != nil {
		return "", err
	}

	return jwt.GenerateToken(session_duration_in_minutes)
}

func (obj *SystemService) ExtendSession() (string, error) {
	system_data, err := obj.GetSystemData()

	if err != nil {
		return "", err
	}

	return jwt.GenerateToken(system_data.SessionDurationInMinutes)
}

func (obj *SystemService) UpdateTheme(theme string) error {
	logger.Log.Printf("Setting theme")
	system_data, err := obj.GetSystemData()

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	system_data.Theme = theme

	err = obj.setSystemData(*system_data)

	if err != nil {
		logger.Log.Printf("ERROR: %s", err.Error())
		return err
	}

	return err
}
