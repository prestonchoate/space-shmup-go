package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct {
	texture    rl.Texture2D
	speed      float32
	srcRect    rl.Rectangle
	destRect   rl.Rectangle
	frameSize  int
	frameCount int
	origin     rl.Vector2
}

func (p *Projectile) Setup() {
	p.texture = rl.LoadTexture("assets/projectile/rocket.png")
	p.speed = 7.5
	p.frameCount = 3
	p.frameSize = int(p.texture.Width) / p.frameCount
	p.srcRect.Width = float32(p.frameSize)
	p.srcRect.Height = float32(p.texture.Height)

	p.destRect.Width = float32(p.frameSize)
	p.destRect.Height = float32(p.texture.Height)
}

func (p *Projectile) Draw() {
	rl.DrawTexturePro(p.texture, p.srcRect, p.destRect, p.origin, 0, rl.White)
}

func (p *Projectile) Update() {
	p.destRect.Y -= p.speed
}
