package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	texture  rl.Texture2D
	speed    float32
	origin   rl.Vector2
	srcRect  rl.Rectangle
	destRect rl.Rectangle
}

func (p *Player) Setup() {
	p.speed = 10.0
	p.srcRect.Width = float32(p.texture.Width)
	p.srcRect.Height = float32(p.texture.Height)

	p.destRect.Width = float32(p.texture.Width)
	p.destRect.Height = float32(p.texture.Height)
	p.destRect.X = float32(p.texture.Width)
	p.destRect.Y = float32(WINDOW_HEIGHT - p.texture.Width)
}

func (p *Player) Draw() {
	rl.DrawTexturePro(p.texture, p.srcRect, p.destRect, p.origin, 0, rl.White)
}

func (p *Player) Update() {

}
