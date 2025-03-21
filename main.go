package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	systems "github.com/prestonchoate/space-shmup/Systems"
	assets "github.com/prestonchoate/space-shmup/Systems/Assets"
	saveManager "github.com/prestonchoate/space-shmup/Systems/saveManager"
)

//go:embed assets/* assets/**/*
var assetsFS embed.FS

func main() {
	checkAndCreateLogPaths(filepath.Join(getHomePath(), "Games", "space-shmup-go"), "game.log")
	file, err := os.OpenFile(filepath.Join(getHomePath(), "Games", "space-shmup-go", "game.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	rl.SetTraceLogLevel(rl.LogError)
	log.Printf("\n\n\n\nSTARTING NEW GAME!\n")

	sm := saveManager.GetInstance()

	rl.InitAudioDevice()
	rl.SetAudioStreamBufferSizeDefault(16384)
	defer rl.CloseAudioDevice()

	rl.InitWindow(int32(sm.Data.Settings.ScreenWidth), int32(sm.Data.Settings.ScreenHeight), "Space Shoot Em Up - Raylib Go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(sm.Data.Settings.TargetFPS)
	data, err := assetsFS.ReadFile("assets/raygui-styles/style_cyber.rgs")
	if err == nil {
		raygui.LoadStyleFromMemory(data)
	}

	if rl.IsWindowFullscreen() != sm.Data.Settings.Fullscreen {
		rl.ToggleFullscreen()
	}

	assets.GetAssetManagerInstance().LoadAssets(assetsFS)
	gm := systems.GetGameMangerInstance()

	for !rl.WindowShouldClose() && !gm.ShouldExit() {
		gm.Update()
		gm.Draw()
	}
}

func getHomePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Save Manager: failed to get user home directory: %v\n", err)
		return ""
	}
	return homeDir

}

func checkAndCreateLogPaths(path string, fileName string) bool {
	settingsPath := filepath.Join(path, fileName)
	dirPath := filepath.Dir(settingsPath)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Printf("failed to create log directory: %v\n", err)
		return false
	}

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		file, err := os.Create(settingsPath)
		if err != nil {
			log.Printf("failed to create log file: %v\n", err)
			return false
		}
		file.Close()
		log.Printf("created log file\n")
	} else if err != nil {
		log.Printf("could not check for log file: %v\n", err)
		return false
	}

	return true

}
