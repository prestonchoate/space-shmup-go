package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type InputMap struct {
	keyLeft  int32
	keyRight int32
	keyUp    int32
	keyDown  int32
	keyFire  int32
}

type Player struct {
	texture     rl.Texture2D
	speed       float32
	origin      rl.Vector2
	srcRect     rl.Rectangle
	destRect    rl.Rectangle
	keyMap      InputMap
	projectiles []*Projectile
}

func (p *Player) Setup() {
	p.speed = 2.5
	p.srcRect.Width = float32(p.texture.Width)
	p.srcRect.Height = float32(p.texture.Height)

	p.destRect.Width = float32(p.texture.Width)
	p.destRect.Height = float32(p.texture.Height)
	p.destRect.X = float32(p.texture.Width)
	p.destRect.Y = float32(int32(rl.GetScreenHeight()) - p.texture.Width)

	p.keyMap.keyLeft = rl.KeyLeft
	p.keyMap.keyRight = rl.KeyRight
	p.keyMap.keyUp = rl.KeyUp
	p.keyMap.keyDown = rl.KeyDown
	p.keyMap.keyFire = rl.KeySpace
}

func (p *Player) Draw() {
	rl.DrawTexturePro(p.texture, p.srcRect, p.destRect, p.origin, 0, rl.White)
	for _, proj := range p.projectiles {
		proj.Draw()
	}
}

func (p *Player) Update() {
	p.handlePlayerInput()
	p.clampPlayerBounds()
	for _, proj := range p.projectiles {
		proj.Update()
	}
}

func (p *Player) handlePlayerInput() {
	if rl.IsKeyDown(p.keyMap.keyLeft) {
		p.destRect.X -= p.speed
	}

	if rl.IsKeyDown(p.keyMap.keyRight) {
		p.destRect.X += p.speed
	}

	if rl.IsKeyDown(p.keyMap.keyUp) {
		p.destRect.Y -= p.speed
	}

	if rl.IsKeyDown(p.keyMap.keyDown) {
		p.destRect.Y += p.speed
	}

	if rl.IsKeyPressed(p.keyMap.keyFire) {
		p.fire()
	}
}

func (p *Player) clampPlayerBounds() {
	minWidth := float32(0.0)
	maxWidth := float32(int32(rl.GetScreenWidth()) - p.texture.Width)
	minHeight := float32(0.0)
	maxHeight := float32(int32(rl.GetScreenHeight()) - p.texture.Height)

	if p.destRect.X < minWidth {
		p.destRect.X = minWidth
	}

	if p.destRect.X > maxWidth {
		p.destRect.X = maxWidth
	}

	if p.destRect.Y < minHeight {
		p.destRect.Y = minHeight
	}

	if p.destRect.Y > maxHeight {
		p.destRect.Y = maxHeight
	}
}

func (p *Player) fire() {
	proj := &Projectile{}
	proj.Setup()
	proj.destRect.X = p.destRect.X
	proj.destRect.Y = p.destRect.Y
	p.projectiles = append(p.projectiles, proj)
}
