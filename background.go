package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Background struct {
	texture  rl.Texture2D
	srcRect  rl.Rectangle
	destRect rl.Rectangle
	origin   rl.Vector2
}

func (bg *Background) Setup() {
	bg.srcRect.Width = float32(bg.texture.Width)
	bg.srcRect.Height = float32(bg.texture.Height)
	bg.destRect = rl.Rectangle{X: 0.0, Y: 0.0, Width: WINDOW_WIDTH, Height: WINDOW_HEIGHT}

	rl.SetTextureWrap(bg.texture, rl.RL_TEXTURE_WRAP_REPEAT)
}

func (bg *Background) Draw() {
	rl.DrawTexturePro(bg.texture, bg.srcRect, bg.destRect, bg.origin, 0, rl.White)
}

func (bg *Background) Update() {
	bg.srcRect.Y -= 1
}
