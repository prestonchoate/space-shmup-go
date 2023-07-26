package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
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

	gm := GameManager{}
	gm.GameSetup()

	for !rl.WindowShouldClose() {
		gm.UpdateEntities()
		gm.DrawEntities()
	}
}
