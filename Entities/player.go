package entities

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	assets "github.com/prestonchoate/space-shmup/Systems/Assets"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
	events "github.com/prestonchoate/space-shmup/Systems/Events"
	events_data "github.com/prestonchoate/space-shmup/Systems/Events/Data"
	"github.com/prestonchoate/space-shmup/Systems/saveManager"
)

const DEFAULT_DAMAGE_TICKS = 15
const DEFAULT_FIRE_RATE = 35

type Player struct {
	id          uuid.UUID
	texture     rl.Texture2D
	speed       float32
	origin      rl.Vector2
	srcRect     rl.Rectangle
	destRect    rl.Rectangle
	keyMap      systems_data.InputMap
	projPool    ObjectPool[*Projectile]
	projTex     rl.Texture2D
	health      int
	score       int
	active      bool
	damaged     bool
	damageTicks int
	fireRate    float32
	fireDelay   float32
}

func CreatePlayer(tex *rl.Texture2D, projTex *rl.Texture2D, keys systems_data.InputMap) *Player {
	p := &Player{
		id:      uuid.New(),
		texture: *(tex),
		speed:   350,
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
		fireRate:    DEFAULT_FIRE_RATE,
	}
	events.GetEventManagerInstance().Subscribe(events_data.GameSettingsUpdated, p.handleSettingsUpdate)
	return p
}

func (p *Player) handleSettingsUpdate(event events.Event) {
	if data, ok := event.Data.(events_data.UpdateSettingsData); ok {
		p.keyMap = data.NewSettings.Keys
	}
}

func (p *Player) Reset() {
	p.health = 100
	p.active = true
	p.destRect = rl.NewRectangle(float32(p.texture.Width),
		float32(rl.GetScreenHeight()-int(p.texture.Height)),
		float32(p.texture.Width),
		float32(p.texture.Height))
	p.projPool.Reset()
}

func (p *Player) Draw() {
	if !p.active {
		return
	}

	if p.health <= 0 {
		return
	}

	rl.DrawText(fmt.Sprint("Fire Rate: ", p.fireRate), 10, 30, 10, rl.Blue)

	tint := rl.White
	if p.damaged {
		tint = rl.Red
	}

	rl.DrawTexturePro(p.texture, p.srcRect, p.destRect, p.origin, 0, tint)
	for _, proj := range p.projPool.activePool {
		proj.Draw()
	}
}

func (p *Player) Update(delta float32) {
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

	p.fireDelay += delta * p.fireRate

	p.handlePlayerInput(delta)
	p.clampPlayerBounds()
	for _, proj := range p.projPool.activePool {
		proj.Update(delta)
		if proj.destRect.Y <= -(proj.destRect.Height) {
			p.projPool.Return(proj)
		}
	}

	if p.health <= 0 {
		events.GetEventManagerInstance().Emit("changeState", events_data.ChangeStateData{NewState: systems_data.GameOver})
	}
}

func (p *Player) GetID() uuid.UUID {
	return p.id
}

// TODO: Refactor to include delta time and normalize the movement speed
func (p *Player) handlePlayerInput(delta float32) {
	if rl.IsKeyDown(p.keyMap.KeyLeft) {
		p.destRect.X -= p.speed * delta
	}

	if rl.IsKeyDown(p.keyMap.KeyRight) {
		p.destRect.X += p.speed * delta
	}

	if rl.IsKeyDown(p.keyMap.KeyUp) {
		p.destRect.Y -= p.speed * delta
	}

	if rl.IsKeyDown(p.keyMap.KeyDown) {
		p.destRect.Y += p.speed * delta
	}

	if rl.IsKeyDown(p.keyMap.KeyFire) {
		p.fire()
	}

	if rl.IsKeyDown(rl.KeyLeftShift) && rl.IsKeyPressed(rl.KeyEnd) {
		p.TakeDamage(10000000)
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

// TODO: Limit fire rate based on player stat
func (p *Player) fire() {
	if p.fireDelay < 10 {
		return
	}

	proj := p.projPool.Get()
	if proj.texture.ID == 0 {
		proj.texture = p.projTex
		proj.frameSize = int(p.projTex.Width) / 3
		proj.srcRect = rl.NewRectangle(0.0, 0.0, float32(proj.frameSize), float32(p.projTex.Height))
		proj.destRect = rl.NewRectangle(0.0, 0.0, float32(proj.frameSize)*float32(proj.scale), float32(p.projTex.Height)*float32(proj.scale))
	}
	proj.destRect.X = p.destRect.X + (float32(p.texture.Width) / 3.75)
	proj.destRect.Y = p.destRect.Y
	sound, ok := assets.GetAssetManagerInstance().GetSound("assets/sfx/laser.wav")
	if ok {
		sfxVolume := saveManager.GetInstance().Data.Settings.SfxVolume
		rl.PlaySound(sound)
		rl.SetSoundVolume(sound, sfxVolume)
		pitchAdj := rand.Float32() / 2
		rl.SetSoundPitch(sound, 1+pitchAdj)
	}
	p.fireDelay = 0.0
}

func (p *Player) TakeDamage(dmg int) {
	if !p.damaged {
		p.health -= dmg
		p.damaged = true
	}
	if p.health <= 0 {
		p.health = 0
		events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{NewState: systems_data.GameOver})
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
	speed := 750

	return &Projectile{
		id:         uuid.New(),
		speed:      float32(speed),
		frameCount: frameCount,
		scale:      float32(scale),
		framespeed: 8,
	}
}
