package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct {
	texture      rl.Texture2D
	speed        float32
	srcRect      rl.Rectangle
	destRect     rl.Rectangle
	frameSize    int
	frameCount   int
	origin       rl.Vector2
	frameLimiter int
	framespeed   int
	scale        float32
}

func (p *Projectile) Setup() {
	p.texture = rl.LoadTexture("assets/projectile/rocket.png")
	p.speed = 7.5
	p.frameCount = 3
	p.frameSize = int(p.texture.Width) / p.frameCount
	p.scale = 3.0

	p.srcRect.Width = float32(p.frameSize)
	p.srcRect.Height = float32(p.texture.Height)

	p.destRect.Width = float32(p.frameSize) * p.scale
	p.destRect.Height = float32(p.texture.Height) * p.scale
	rl.SetTextureWrap(p.texture, rl.RL_TEXTURE_WRAP_REPEAT)
	p.framespeed = 8
}

func (p *Projectile) Draw() {
	rl.DrawTexturePro(p.texture, p.srcRect, p.destRect, p.origin, 0, rl.White)
}

func (p *Projectile) Update() {
	p.destRect.Y -= p.speed
	p.frameLimiter++
	if p.frameLimiter >= (int(rl.GetFPS()) / p.framespeed) {
		p.frameLimiter = 0
		p.srcRect.X += float32(p.frameSize)
	}
}