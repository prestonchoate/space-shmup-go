package entities

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type InputMap struct {
	keyLeft  int32
	keyRight int32
	keyUp    int32
	keyDown  int32
	keyFire  int32
}

type Player struct {
	id       uuid.UUID
	texture  rl.Texture2D
	speed    float32
	origin   rl.Vector2
	srcRect  rl.Rectangle
	destRect rl.Rectangle
	keyMap   InputMap
	projPool ObjectPool[*Projectile]
}

func CreatePlayer(tex *rl.Texture2D) *Player {
	return &Player{
		id:      uuid.New(),
		texture: *(tex),
		speed:   2.5,
		origin:  rl.Vector2{X: 0.0, Y: 0.0},
		srcRect: rl.NewRectangle(0.0, 0.0, float32(tex.Width), float32(tex.Height)),
		destRect: rl.NewRectangle(float32(tex.Width),
			float32(rl.GetScreenHeight()-int(tex.Height)),
			float32(tex.Width),
			float32(tex.Height)),
		keyMap: InputMap{
			keyLeft:  rl.KeyLeft,
			keyUp:    rl.KeyUp,
			keyRight: rl.KeyRight,
			keyDown:  rl.KeyDown,
			keyFire:  rl.KeySpace,
		},
		projPool: ObjectPool[*Projectile]{
			activePool:   make(map[uuid.UUID]*Projectile),
			inactivePool: make([]*Projectile, 0, 200),
			createFn:     createProjectile,
		},
	}
}

func (p *Player) Draw() {
	rl.DrawTexturePro(p.texture, p.srcRect, p.destRect, p.origin, 0, rl.White)
	rl.DrawRectangleLines(int32(p.destRect.X), int32(p.destRect.Y), int32(p.destRect.Width), int32(p.destRect.Height), rl.Green)
	for _, proj := range p.projPool.activePool {
		proj.Draw()
		rl.DrawRectangleLines(int32(proj.destRect.X), int32(proj.destRect.Y), int32(proj.destRect.Width), int32(proj.destRect.Height), rl.Blue)
	}
}

func (p *Player) Update() {
	p.handlePlayerInput()
	p.clampPlayerBounds()
	for _, proj := range p.projPool.activePool {
		proj.Update()
		if proj.destRect.Y <= -(proj.destRect.Height) {
			p.projPool.Return(proj)
		}
	}
}

func (p *Player) GetID() uuid.UUID {
	return p.id
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
	proj := p.projPool.Get()
	proj.destRect.X = p.destRect.X + (float32(p.texture.Width) / 3.75)
	proj.destRect.Y = p.destRect.Y
}

func (p *Player) Destroy() {
	fmt.Println("Destroying player")
}

func (p *Player) DestroyProjectile(proj *Projectile) {
	p.projPool.Return(proj)
}

func createProjectile() GameEntity {
	frameCount := 3
	frameSize := int(ProjectileTexture.Width) / 3
	scale := 3.0
	speed := 7.5

	return &Projectile{
		id:         uuid.New(),
		texture:    ProjectileTexture,
		speed:      float32(speed),
		frameCount: frameCount,
		frameSize:  frameSize,
		scale:      float32(scale),
		srcRect:    rl.NewRectangle(0.0, 0.0, float32(frameSize), float32(ProjectileTexture.Height)),
		destRect:   rl.NewRectangle(0.0, 0.0, float32(frameSize)*float32(scale), float32(ProjectileTexture.Height)*float32(scale)),
		framespeed: 8,
	}
}
