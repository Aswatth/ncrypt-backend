package services

import (
	"fmt"
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
	err := service.setSystemData(models.SystemData{LoginCount: 10, LastLoginDateTime: now, IsLoggedIn: false, CurrentLoginDateTime: now, AutomaticBackup: false, AutomaticBackupLocation: "", BackupFileName: "", SessionTimeInMinutes: 20})

	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestGetSystemData(t *testing.T) {
	service := new(SystemService)
	service.Init()

	now := time.Now().Format(time.RFC3339)
	initial_data := models.SystemData{LoginCount: 1, LastLoginDateTime: now, IsLoggedIn: false, CurrentLoginDateTime: now, AutomaticBackup: false, AutomaticBackupLocation: "", BackupFileName: "", SessionTimeInMinutes: 20}

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

	err := service.Setup("12345", false, "", "")

	if err != nil {
		t.Error(err.Error())
	}

	fetched_system_data, err := service.GetSystemData()

	if err != nil {
		t.Error(err.Error())
	}

	if fetched_system_data.AutomaticBackup != false || fetched_system_data.AutomaticBackupLocation != "" || fetched_system_data.BackupFileName != "" {
		t.Error("Automatic backup should be empty")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestSetup_WithBackup(t *testing.T) {
	system_service_test_cleanup()

	service := new(SystemService)
	service.Init()

	err := service.Setup("12345", true, "", "my_backup")

	if err != nil {
		t.Error(err.Error())
	}

	fetched_system_data, err := service.GetSystemData()

	if err != nil {
		t.Error(err.Error())
	}

	if fetched_system_data.AutomaticBackup != true {
		t.Error("Automatic backup should be empty")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestSignIn(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"

	err := service.Setup(password, false, "", "")

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

	err := service.Setup(password, false, "", "")

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

	err := service.Setup(password, false, "", "")

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

	err := service.Setup(password, false, "", "")

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
	err := service.Setup(password, false, "", "")

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
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestExport_WithCustomPath(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, false, "", "")

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
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestExport_WithInvalidPath(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, false, "", "")

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
	err := service.Setup(password, false, "", "")

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.txt", "")
	if err == nil {
		t.Error("should result in an error due to incorrect format")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestGeneratePassword(t *testing.T) {
	service := new(SystemService)
	service.Init()

	t.Run("Default case", func(t *testing.T) {
		t.Parallel()
		password := service.GeneratePassword(false, false, false, 8)

		if password == "" {
			t.Error("Generated password cannot be empty")
		}
	},
	)

	t.Run("With digits", func(t *testing.T) {
		t.Parallel()
		password := service.GeneratePassword(true, false, false, 8)

		if password == "" {
			t.Error("Generated password cannot be empty")
		}
	})

	t.Run("With uppercase characters", func(t *testing.T) {
		t.Parallel()
		password := service.GeneratePassword(false, true, false, 8)

		if password == "" {
			t.Error("Generated password cannot be empty")
		}
	})

	t.Run("With special characters", func(t *testing.T) {
		t.Parallel()
		password := service.GeneratePassword(false, false, true, 8)

		if password == "" {
			t.Error("Generated password cannot be empty")
		}
	})

	t.Run("With digits and upper case", func(t *testing.T) {
		t.Parallel()
		password := service.GeneratePassword(true, true, false, 8)

		if password == "" {
			t.Error("Generated password cannot be empty")
		}
	})

	t.Run("With uppercase and special characters", func(t *testing.T) {
		t.Parallel()
		password := service.GeneratePassword(false, true, true, 8)

		if password == "" {
			t.Error("Generated password cannot be empty")
		}
	})

	t.Run("With digits and special characters", func(t *testing.T) {
		t.Parallel()
		password := service.GeneratePassword(true, false, true, 8)

		if password == "" {
			t.Error("Generated password cannot be empty")
		}
	})

	t.Run("With digits, uppercase and special characters", func(t *testing.T) {
		t.Parallel()
		password := service.GeneratePassword(true, true, true, 8)

		if password == "" {
			t.Error("Generated password cannot be empty")
		}
	})

	for i := 8; i <= 16; i++ {
		t.Run("With length "+fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			password := service.GeneratePassword(false, false, true, 16)

			if password == "" {
				t.Error("Generated password cannot be empty")
			}
		})
	}
}

func TestImport(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, false, "", "")

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.ncrypt", "")
	if err != nil {
		t.Error(err.Error())
	}

	//Changing data
	service.UpdateAutomaticBackup(true, "", "backup")

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

		if system_data.AutomaticBackup != true || system_data.BackupFileName != "backup" || system_data.AutomaticBackupLocation != "" {
			t.Error("incorrect data")
		}
	} else {
		err = os.Remove("test_export.ncrypt")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestImport_WithIncorrectMasterPassword(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, false, "", "")

	if err != nil {
		t.Error(err.Error())
	}

	err = service.Export("test_export.ncrypt", "")
	if err != nil {
		t.Error(err.Error())
	}

	//Changing data
	service.UpdateAutomaticBackup(true, "", "backup")

	if file_info, err := os.Stat("test_export.ncrypt"); os.IsNotExist(err) {
		file_name := file_info.Name()
		err = service.Import(file_name, "", "123")

		if err == nil {
			t.Error("should result in an error")
		}
	} else {
		err = os.Remove("test_export.ncrypt")
	}

	t.Cleanup(system_service_test_cleanup)
}

// Test backup
func TestBackup(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, true, "", "test_backup")

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
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_SettingToTrue(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, false, "", "")

	if err != nil {
		t.Error(err.Error())
	}

	err = service.UpdateAutomaticBackup(true, "D:", "my_backup.ncrypt")

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()
	if err != nil {
		t.Error(err.Error())
	}
	if system_data.AutomaticBackup != true || system_data.AutomaticBackupLocation != "D:" || system_data.BackupFileName != "my_backup.ncrypt" {
		t.Error("inavlid information")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_SettingToFalse(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, true, "D:", "my_backup.ncrypt")

	if err != nil {
		t.Error(err.Error())
	}

	err = service.UpdateAutomaticBackup(false, "D:", "my_backup.ncrypt")

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()
	if err != nil {
		t.Error(err.Error())
	}
	if system_data.AutomaticBackup != false || system_data.AutomaticBackupLocation != "D:" || system_data.BackupFileName != "my_backup.ncrypt" {
		t.Error("inavlid information")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_ChangingFileName(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, true, "D:", "my_backup.ncrypt")

	if err != nil {
		t.Error(err.Error())
	}

	err = service.UpdateAutomaticBackup(true, "D:", "backup.ncrypt")

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()
	if err != nil {
		t.Error(err.Error())
	}
	if system_data.AutomaticBackup != true || system_data.AutomaticBackupLocation != "D:" || system_data.BackupFileName != "backup.ncrypt" {
		t.Error("inavlid information")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_ChangingPath(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, true, "D:", "my_backup.ncrypt")

	if err != nil {
		t.Error(err.Error())
	}

	err = service.UpdateAutomaticBackup(true, "C:", "my_backup.ncrypt")

	if err != nil {
		t.Error(err.Error())
	}

	system_data, err := service.GetSystemData()
	if err != nil {
		t.Error(err.Error())
	}
	if system_data.AutomaticBackup != true || system_data.AutomaticBackupLocation != "C:" || system_data.BackupFileName != "my_backup.ncrypt" {
		t.Error("inavlid information")
	}

	t.Cleanup(system_service_test_cleanup)
}

func TestUpadteAutomaticBackupData_EmptyFileName(t *testing.T) {
	service := new(SystemService)
	service.Init()

	password := "12345"
	err := service.Setup(password, false, "", "")

	if err != nil {
		t.Error(err.Error())
	}

	err = service.UpdateAutomaticBackup(true, "", "")

	if err == nil {
		t.Error("Should cause an error")
	}

	t.Cleanup(system_service_test_cleanup)
}

func system_service_test_cleanup() {
	os.RemoveAll(os.Getenv("STORAGE_FOLDER"))
}
