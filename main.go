package main

import (
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	systems "github.com/prestonchoate/space-shmup/Systems"
)

const (
	WINDOW_WIDTH  = 1280
	WINDOW_HEIGHT = 720
	TARGET_FPS    = 120
)

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Space Shoot Em Up - Raylib Go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(TARGET_FPS)
	raygui.LoadStyle("assets/raygui-styles/style_cyber.rgs")

	gm := systems.GetGameMangerInstance()

	for !rl.WindowShouldClose() {
		gm.Update()
		gm.Draw()
	}
}
