package main

import (
	"embed"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	systems "github.com/prestonchoate/space-shmup/Systems"
	assets "github.com/prestonchoate/space-shmup/Systems/Assets"
	saveManager "github.com/prestonchoate/space-shmup/Systems/saveManager"
)

//go:embed assets/* assets/**/*
var assetsFS embed.FS

func main() {
	saveManager := saveManager.GetInstance()

	rl.SetTraceLogLevel(rl.LogError)

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.InitWindow(int32(saveManager.Data.Settings.ScreenWidth), int32(saveManager.Data.Settings.ScreenHeight), "Space Shoot Em Up - Raylib Go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(saveManager.Data.Settings.TargetFPS)
	data, err := assetsFS.ReadFile("assets/raygui-styles/style_cyber.rgs")
	if err == nil {
		raygui.LoadStyleFromMemory(data)
	}

	if rl.IsWindowFullscreen() != saveManager.Data.Settings.Fullscreen {
		rl.ToggleFullscreen()
	}

	assets.GetAssetManagerInstance().LoadAssets(assetsFS)
	gm := systems.GetGameMangerInstance()

	for !rl.WindowShouldClose() && !gm.ShouldExit() {
		gm.Update()
		gm.Draw()
	}
}
