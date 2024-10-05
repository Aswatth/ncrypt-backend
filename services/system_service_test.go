package services

import (
	"ncrypt/models"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSetSystemData(t *testing.T) {
	service := new(SystemService)
	service.Init()

	now := time.Now().Format(time.RFC3339)
	err := service.setSystemData(models.SystemData{LoginCount: 10, LastLoginDateTime: now, IsLoggedIn: false, CurrentLoginDateTime: now, AutoBackupSetting: models.AutoBackupSetting{IsEnabled: false, BackupLocation: "", BackupFileName: ""}, SessionDurationInMinutes: 20})

	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestGetSystemData(t *testing.T) {
	service := new(SystemService)
	service.Init()

	now := time.Now().Format(time.RFC3339)
	initial_data := models.SystemData{LoginCount: 1, LastLoginDateTime: now, IsLoggedIn: false, CurrentLoginDateTime: now, AutoBackupSetting: models.AutoBackupSetting{IsEnabled: false, BackupLocation: "", BackupFileName: ""}, SessionDurationInMinutes: 20}

	err := service.setSystemData(initial_data)

	if err != nil {
		t.Error(err.Error())
	}

	_, err = service.GetSystemData()

	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestSetup_WithNoBackup(t *testing.T) {
	system_service_test_cleanup()

	service := new(SystemService)
	service.Init()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup("12345", auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_system_data, err := service.GetSystemData()

	if err != nil {
		t.Error(err.Error())
	}

	if fetched_system_data.AutoBackupSetting.IsEnabled != false || fetched_system_data.AutoBackupSetting.BackupLocation != "" || fetched_system_data.AutoBackupSetting.BackupFileName != "" {
		t.Error("Automatic backup should be empty")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestSetup_WithBackup(t *testing.T) {
	system_service_test_cleanup()

	service := new(SystemService)
	service.Init()

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "../services/"
	auto_backup_setting["backup_file_name"] = "test_backup.ncrypt"

	err := service.Setup("12345", auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	fetched_system_data, err := service.GetSystemData()

	if err != nil {
		t.Error(err.Error())
	}

	if fetched_system_data.AutoBackupSetting.IsEnabled != true {
		t.Error("Automatic backup should be empty")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestSignIn(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	token, err := service.SignIn(password)

	if err != nil {
		t.Error(err.Error())
	}

	if token == "" {
		t.Error("Empty token")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestSignIn_WithIncorrectPassword(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"

	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	token, err := service.SignIn("123")

	if err != nil && err.Error() != "invalid password" {
		t.Error(err.Error())
	}

	if token != "" {
		t.Error("Should not return token")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestLogout(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	token, err := service.SignIn(password)

	if err != nil {
		t.Error(err.Error())
	}

	if token == "" {
		t.Error("Token not found")
	}

	err = service.Logout()

	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestLogout_WithoutSignin(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Logout()

	if err == nil {
		t.Error("Should result in an un-authorized error")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestExport_CurrentFolder(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.ncrypt", "")
	if err != nil {
		t.Error(err.Error())
	}

	if _, err := os.Stat("test_export.ncrypt"); os.IsNotExist(err) {
		t.Error("Exported file not found")
	} else {
		err = os.Remove("test_export.ncrypt")
		if err != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestExport_WithCustomPath(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.ncrypt", "..\\models")
	if err != nil {
		t.Error(err.Error())
	}

	if _, err := os.Stat("..\\models\\test_export.ncrypt"); os.IsNotExist(err) {
		t.Error("Exported file not found")
	} else {
		err = os.Remove("..\\models\\test_export.ncrypt")
		if err != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestExport_WithInvalidPath(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.ncrypt", "..\\test")
	if err == nil {
		t.Error("should result in an error as folder is not found")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestExport_IncorrectFormat(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.txt", "")
	if err == nil {
		t.Error("should result in an error due to incorrect format")
	}

	t.Cleanup(system_service_test_cleanup)
}

// func TestGeneratePassword(t *testing.T) {
// 	service := new(SystemService)
// 	service.Init()

// 	t.Run("Default case", func(t *testing.T) {
// 		t.Parallel()
// 		password := service.GeneratePassword(false, false, false, 8)

// 		if password == "" {
// 			t.Error("Generated password cannot be empty")
// 		}
// 	},
// 	)

// 	t.Run("With digits", func(t *testing.T) {
// 		t.Parallel()
// 		password := service.GeneratePassword(true, false, false, 8)

// 		if password == "" {
// 			t.Error("Generated password cannot be empty")
// 		}
// 	})

// 	t.Run("With uppercase characters", func(t *testing.T) {
// 		t.Parallel()
// 		password := service.GeneratePassword(false, true, false, 8)

// 		if password == "" {
// 			t.Error("Generated password cannot be empty")
// 		}
// 	})

// 	t.Run("With special characters", func(t *testing.T) {
// 		t.Parallel()
// 		password := service.GeneratePassword(false, false, true, 8)

// 		if password == "" {
// 			t.Error("Generated password cannot be empty")
// 		}
// 	})

// 	t.Run("With digits and upper case", func(t *testing.T) {
// 		t.Parallel()
// 		password := service.GeneratePassword(true, true, false, 8)

// 		if password == "" {
// 			t.Error("Generated password cannot be empty")
// 		}
// 	})

// 	t.Run("With uppercase and special characters", func(t *testing.T) {
// 		t.Parallel()
// 		password := service.GeneratePassword(false, true, true, 8)

// 		if password == "" {
// 			t.Error("Generated password cannot be empty")
// 		}
// 	})

// 	t.Run("With digits and special characters", func(t *testing.T) {
// 		t.Parallel()
// 		password := service.GeneratePassword(true, false, true, 8)

// 		if password == "" {
// 			t.Error("Generated password cannot be empty")
// 		}
// 	})

// 	t.Run("With digits, uppercase and special characters", func(t *testing.T) {
// 		t.Parallel()
// 		password := service.GeneratePassword(true, true, true, 8)

// 		if password == "" {
// 			t.Error("Generated password cannot be empty")
// 		}
// 	})

// 	for i := 8; i <= 16; i++ {
// 		t.Run("With length "+fmt.Sprintf("%d", i), func(t *testing.T) {
// 			t.Parallel()
// 			password := service.GeneratePassword(false, false, true, 16)

// 			if password == "" {
// 				t.Error("Generated password cannot be empty")
// 			}
// 		})
// 	}
// }

func TestImport(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.ncrypt", "")
	if err != nil {
		t.Error(err.Error())
	}

	//Changing data
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = "backup"

	service.UpdateAutomaticBackup(auto_backup_setting)

	if file_info, err := os.Stat("test_export.ncrypt"); os.IsNotExist(err) {
		file_name := file_info.Name()
		err = service.Import(file_name, "", password)

		if err != nil {
			t.Error(err.Error())
		}

		system_data, err := service.GetSystemData()

		if err != nil {
			t.Error(err.Error())
		}

		if system_data.AutoBackupSetting.IsEnabled != true || system_data.AutoBackupSetting.BackupFileName != "backup" || system_data.AutoBackupSetting.BackupLocation != "" {
			t.Error("incorrect data")
		}
	} else {
		err = os.Remove("test_export.ncrypt")
		if err != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestImport_WithIncorrectMasterPassword(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.ncrypt", "")
	if err != nil {
		t.Error(err.Error())
	}

	//Changing data
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = "backup"

	service.UpdateAutomaticBackup(auto_backup_setting)

	if file_info, err := os.Stat("test_export.ncrypt"); os.IsNotExist(err) {
		file_name := file_info.Name()
		err = service.Import(file_name, "", "123")

		if err == nil {
			t.Error("should result in an error")
		}
	} else {
		err = os.Remove("test_export.ncrypt")
		if err != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(system_service_test_cleanup)
}

// Test backup
func TestBackup(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = "test_backup"

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Backup()
	if err != nil {
		t.Error(err.Error())
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err.Error())
		return
	}

	files, err := os.ReadDir(cwd)
	if err != nil {
		t.Error(err.Error())
	}

	// Iterate over files and find a match
	var backup_file_name string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "test_backup") {
			backup_file_name = file.Name()
			break
		}
	}

	if _, err := os.Stat(backup_file_name); os.IsNotExist(err) {
		t.Error("Backup file not found")
	} else {
		err = os.Remove(backup_file_name)
		if err != nil {
			t.Error(err.Error())
		}
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_SettingToTrue(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err = service.UpdateAutomaticBackup(auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()
	if err != nil {
		t.Error(err.Error())
	}
	if system_data.AutoBackupSetting.IsEnabled != true || system_data.AutoBackupSetting.BackupLocation != "D:" || system_data.AutoBackupSetting.BackupFileName != "my_backup.ncrypt" {
		t.Error("inavlid information")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_SettingToFalse(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting["is_enabled"] = false

	err = service.UpdateAutomaticBackup(auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()
	if err != nil {
		t.Error(err.Error())
	}
	if system_data.AutoBackupSetting.IsEnabled != false || system_data.AutoBackupSetting.BackupLocation != "D:" || system_data.AutoBackupSetting.BackupFileName != "my_backup.ncrypt" {
		t.Error("inavlid information")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_ChangingFileName(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting["backup_file_name"] = "backup.ncrypt"

	err = service.UpdateAutomaticBackup(auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()
	if err != nil {
		t.Error(err.Error())
	}
	if system_data.AutoBackupSetting.IsEnabled != true || system_data.AutoBackupSetting.BackupLocation != "D:" || system_data.AutoBackupSetting.BackupFileName != "backup.ncrypt" {
		t.Error("inavlid information")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_ChangingPath(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = true
	auto_backup_setting["backup_location"] = "D:"
	auto_backup_setting["backup_file_name"] = "my_backup.ncrypt"

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting["backup_location"] = "C:"

	err = service.UpdateAutomaticBackup(auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()
	if err != nil {
		t.Error(err.Error())
	}
	if system_data.AutoBackupSetting.IsEnabled != true || system_data.AutoBackupSetting.BackupLocation != "C:" || system_data.AutoBackupSetting.BackupFileName != "my_backup.ncrypt" {
		t.Error("inavlid information")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_EmptyFileName(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	auto_backup_setting := make(map[string]interface{})
	auto_backup_setting["is_enabled"] = false
	auto_backup_setting["backup_location"] = ""
	auto_backup_setting["backup_file_name"] = ""

	err := service.Setup(password, auto_backup_setting)

	if err != nil {
		t.Error(err.Error())
	}

	auto_backup_setting["is_enabled"] = true

	err = service.UpdateAutomaticBackup(auto_backup_setting)

	if err == nil {
		t.Error("Should cause an error")
	}

	t.Cleanup(system_service_test_cleanup)
}

func system_service_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
}
