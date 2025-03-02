package main

import (
	"embed"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	systems "github.com/prestonchoate/space-shmup/Systems"
	assets "github.com/prestonchoate/space-shmup/Systems/Assets"
)

const (
	TARGET_FPS = 120
)

//go:embed assets/* assets/**/*
var assetsFS embed.FS

func main() {
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.InitWindow(0, 0, "Space Shoot Em Up - Raylib Go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(TARGET_FPS)
	data, err := assetsFS.ReadFile("assets/raygui-styles/style_cyber.rgs")
	if err == nil {
		raygui.LoadStyleFromMemory(data)
	}

	if !rl.IsWindowFullscreen() {
		rl.ToggleFullscreen()
	}

	assets.GetAssetManagerInstance().LoadAssets(assetsFS)
	gm := systems.GetGameMangerInstance()

	for !rl.WindowShouldClose() && !gm.ShouldExit() {
		gm.Update()
		gm.Draw()
	}
}
