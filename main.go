package main

import (
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

	gm := systems.CreateGameManager()

	for !rl.WindowShouldClose() {
		gm.Update()
		gm.Draw()
	}
}
