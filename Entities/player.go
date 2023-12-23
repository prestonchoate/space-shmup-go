package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

const DEFAULT_DAMAGE_TICKS = 15

type InputMap struct {
	KeyLeft  int32
	KeyRight int32
	KeyUp    int32
	KeyDown  int32
	KeyFire  int32
}

type Player struct {
	id          uuid.UUID
	texture     rl.Texture2D
	speed       float32
	origin      rl.Vector2
	srcRect     rl.Rectangle
	destRect    rl.Rectangle
	keyMap      InputMap
	projPool    ObjectPool[*Projectile]
	projTex     rl.Texture2D
	health      int
	score       int
	active      bool
	damaged     bool
	damageTicks int
}

func CreatePlayer(tex *rl.Texture2D, projTex *rl.Texture2D, keys InputMap) *Player {
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
		keyMap: keys,
		projPool: ObjectPool[*Projectile]{
			activePool:   make(map[uuid.UUID]*Projectile),
			inactivePool: make([]*Projectile, 0, 200),
			createFn:     createProjectile,
		},
		projTex:     *projTex,
		health:      100,
		active:      true,
		damageTicks: DEFAULT_DAMAGE_TICKS,
	}
}

func (p *Player) Draw() {
	if !p.active {
		return
	}

	if p.health <= 0 {
		return
	}
	tint := rl.White
	if p.damaged {
		tint = rl.Red
	}

	rl.DrawTexturePro(p.texture, p.srcRect, p.destRect, p.origin, 0, tint)
	for _, proj := range p.projPool.activePool {
		proj.Draw()
	}
}

func (p *Player) Update() {
	if !p.active {
		return
	}

	if p.health <= 0 {
		// Move player off screen
		p.destRect.X = -1000
		return
	}

	if p.damaged {
		p.damageTicks--
		if p.damageTicks <= 0 {
			p.damaged = false
			p.damageTicks = DEFAULT_DAMAGE_TICKS
		}
	}

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
	if rl.IsKeyDown(p.keyMap.KeyLeft) {
		p.destRect.X -= p.speed
	}

	if rl.IsKeyDown(p.keyMap.KeyRight) {
		p.destRect.X += p.speed
	}

	if rl.IsKeyDown(p.keyMap.KeyUp) {
		p.destRect.Y -= p.speed
	}

	if rl.IsKeyDown(p.keyMap.KeyDown) {
		p.destRect.Y += p.speed
	}

	if rl.IsKeyPressed(p.keyMap.KeyFire) {
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
	if proj.texture.ID == 0 {
		proj.texture = p.projTex
		proj.frameSize = int(p.projTex.Width) / 3
		proj.srcRect = rl.NewRectangle(0.0, 0.0, float32(proj.frameSize), float32(p.projTex.Height))
		proj.destRect = rl.NewRectangle(0.0, 0.0, float32(proj.frameSize)*float32(proj.scale), float32(p.projTex.Height)*float32(proj.scale))
	}
	proj.destRect.X = p.destRect.X + (float32(p.texture.Width) / 3.75)
	proj.destRect.Y = p.destRect.Y
}

func (p *Player) TakeDamage(dmg int) {
	if !p.damaged {
		p.health -= dmg
		p.damaged = true
	}
	if p.health <= 0 {
		p.health = 0
	}
}

func (p *Player) AddScore(score int) {
	p.score += score
}

func (p *Player) DestroyProjectile(proj *Projectile) {
	p.projPool.Return(proj)
	p.score += 10
}

func (p *Player) Activate(active bool) {
	p.active = active
}

func (p *Player) GetRect() rl.Rectangle {
	return p.destRect
}

func (p *Player) GetProjeciles() map[uuid.UUID]*Projectile {
	return p.projPool.activePool
}

func (p *Player) GetHealth() int {
	return p.health
}

func (p *Player) GetScore() int {
	return p.score
}

func createProjectile() GameEntity {
	frameCount := 3
	scale := 3.0
	speed := 7.5

	return &Projectile{
		id:         uuid.New(),
		speed:      float32(speed),
		frameCount: frameCount,
		scale:      float32(scale),
		framespeed: 8,
	}
}
