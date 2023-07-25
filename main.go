package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WINDOW_WIDTH  = 1920
	WINDOW_HEIGHT = 1080
	TARGET_FPS    = 60
)

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Space Shoot Em Up - Raylib Go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(TARGET_FPS)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Just the beginning!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}
}
