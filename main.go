package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WINDOW_WIDTH  = 1280
	WINDOW_HEIGHT = 720
	TARGET_FPS    = 60
)

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Space Shoot Em Up - Raylib Go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(TARGET_FPS)

	bg := rl.LoadTexture("assets/background.jpg")
	bg_width := bg.Width
	bg_height := bg.Height

	bg_src_rec := rl.Rectangle{X: 0.0, Y: 0.0, Width: float32(bg_width), Height: float32(bg_height)}
	bg_dest_rec := rl.Rectangle{X: 0.0, Y: 0.0, Width: WINDOW_WIDTH, Height: WINDOW_HEIGHT}
	bg_origin := rl.Vector2{X: 0.0, Y: 0.0}

	rl.SetTextureWrap(bg, rl.RL_TEXTURE_WRAP_REPEAT)

	for !rl.WindowShouldClose() {
		bg_src_rec.Y -= 1
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Draw Background
		rl.DrawTexturePro(bg, bg_src_rec, bg_dest_rec, bg_origin, 0, rl.White)

		rl.EndDrawing()
	}
}
