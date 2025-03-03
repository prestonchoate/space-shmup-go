package saveManager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
)

var instance *SaveManager

type SaveManager struct {
	FilePath string
	FileName string
	Data     SaveData
}

type SaveData struct {
	Settings systems_data.GameSettings `json:"settings"`
}

func GetInstance() *SaveManager {
	if instance == nil {
		instance = &SaveManager{
			FilePath: "Games/Saves/space-shmup-go/",
			FileName: "settings.json",
		}
		ok := instance.checkAndCreatePaths()
		if !ok {
			return instance
		}

		ok = instance.loadData()
		if !ok {
			return instance
		}
	}
	return instance
}

func (sm *SaveManager) checkAndCreatePaths() bool {
	settingsPath := sm.getFullFilePath()
	dirPath := filepath.Dir(settingsPath)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Printf("Save Manager: failed to create settings directory: %v\n", err)
		return false
	}

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		file, err := os.Create(settingsPath)
		if err != nil {
			log.Printf("Save Manager: failed to create settings file: %v\n", err)
			return false
		}
		file.Close()
		log.Printf("Save Manager: created settings file\n")
	} else if err != nil {
		log.Printf("Save Manager: could not check for settings file: %v\n", err)
		return false
	}

	return true

}

func (sm *SaveManager) getFullFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Save Manager: failed to get user home directory: %v\n", err)
		return ""
	}
	return filepath.Join(homeDir, sm.FilePath, sm.FileName)
}

func (sm *SaveManager) loadData() bool {
	log.Printf("Save Manager: Attempting to load data from %v\n", sm.getFullFilePath())
	data, err := os.ReadFile(sm.getFullFilePath())
	if err != nil {
		log.Printf("Save Manager: failed to read settings file: %v\n", err)
		return false
	}

	if len(data) == 0 {
		log.Println("Save Manager: No data in settings file. Loading defaults")
		return instance.createDefaultSettings()
	}

	err = json.Unmarshal(data, &sm.Data)
	if err != nil {
		log.Printf("Save Manager: failed to parse settings file. Loading default data: %v\n", err)
		return sm.createDefaultSettings()
	}

	log.Printf("Save Manager: successfully loaded data from settings file!\nSettings: %+v\n", sm.Data.Settings)
	return true
}

func (sm *SaveManager) createDefaultSettings() bool {
	settings := systems_data.GameSettings{
		TargetFPS:    120,
		ScreenWidth:  rl.GetScreenWidth(),
		ScreenHeight: rl.GetScreenHeight(),
		Fullscreen:   true,
		Keys: systems_data.InputMap{
			KeyLeft:  rl.KeyA,
			KeyUp:    rl.KeyW,
			KeyRight: rl.KeyD,
			KeyDown:  rl.KeyS,
			KeyFire:  rl.KeySpace,
		},
	}

	sm.Data.Settings = settings

	return sm.saveSettings(sm.Data)
}

func (sm *SaveManager) saveSettings(settings SaveData) bool {
	data, err := json.Marshal(settings)

	if err != nil {
		log.Printf("Save Manager: failed to marshal json settings data: %v\n", err)
		return false
	}

	file, err := os.OpenFile(sm.getFullFilePath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Save Manager: failed to open settings file for writes: %v\n", err)
		return false
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Printf("Save Manager: failed to save settings to file: %v\n", err)
		return false
	}

	fmt.Printf("Save Manager: wrote settings to save file!\n")

	return true
}
