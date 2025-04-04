package entities

import (
	"log"
	"math/rand"
	"time"

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
	origin      rl.Vector2
	srcRect     rl.Rectangle
	destRect    rl.Rectangle
	keyMap      systems_data.InputMap
	projPool    ObjectPool[*Projectile]
	projTex     rl.Texture2D
	score       int
	active      bool
	damaged     bool
	damageTicks int
	fireDelay   float32
	stats       PlayerStats
}

func CreatePlayer(keys systems_data.InputMap) *Player {
	texture, ok := assets.GetAssetManagerInstance().GetTexture("assets/sprites/player/playerShip1_red.png")
	if !ok {
		log.Fatal("couldn't load player texture")
	}

	projTexture, ok := assets.GetAssetManagerInstance().GetTexture("assets/sprites/projectile/laserRed07.png")
	if !ok {
		log.Fatal("couldn't load projectile texture")
	}
	rl.SetTextureWrap(projTexture, rl.WrapRepeat)

	p := &Player{
		id:      uuid.New(),
		texture: texture,
		origin:  rl.Vector2{X: 0.0, Y: 0.0},
		srcRect: rl.NewRectangle(0.0, 0.0, float32(texture.Width), float32(texture.Height)),
		destRect: rl.NewRectangle(float32(texture.Width),
			float32(rl.GetScreenHeight()-int(texture.Height)),
			float32(texture.Width),
			float32(texture.Height)),
		keyMap: keys,
		projPool: ObjectPool[*Projectile]{
			activePool:   make(map[uuid.UUID]*Projectile),
			inactivePool: make([]*Projectile, 0, 200),
			createFn:     createProjectile,
		},
		projTex:     projTexture,
		active:      true,
		damageTicks: DEFAULT_DAMAGE_TICKS,
		fireDelay:   0.0,
		stats:       NewDefaultPlayerStats(),
	}

	// Subscribe to events
	events.GetEventManagerInstance().Subscribe(events_data.GameSettingsUpdated, p.handleSettingsUpdate)
	events.GetEventManagerInstance().Subscribe(events_data.StatUpgradeEvent, p.handleStatUpgrade)

	return p
}

func (p *Player) handleSettingsUpdate(event events.Event) {
	if data, ok := event.Data.(events_data.UpdateSettingsData); ok {
		p.keyMap = data.NewSettings.Keys
	}
}

func (p *Player) Reset() {
	p.stats = NewDefaultPlayerStats()
	p.active = true
	p.fireDelay = 0
	p.damageTicks = DEFAULT_DAMAGE_TICKS
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

	if p.stats.Health <= 0 {
		return
	}

	tint := rl.White
	if p.damaged {
		tint = rl.Red
	}

	for _, proj := range p.projPool.activePool {
		proj.Draw()
	}

	rl.DrawTexturePro(p.texture, p.srcRect, p.destRect, p.origin, 0, tint)
}

func (p *Player) Update(delta float32) {
	if !p.active {
		return
	}

	if p.stats.Health <= 0 {
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

	p.fireDelay += delta * p.stats.FireRate

	p.handlePlayerInput(delta)
	p.clampPlayerBounds()
	for _, proj := range p.projPool.activePool {
		proj.Update(delta)
		if proj.destRect.Y <= -(proj.destRect.Height) {
			p.projPool.Return(proj)
		}
	}

	if p.stats.Health <= 0 {
		events.GetEventManagerInstance().Emit("changeState", events_data.ChangeStateData{NewState: systems_data.GameOver})
	}
}

func (p *Player) GetID() uuid.UUID {
	return p.id
}

// TODO: Refactor to include delta time and normalize the movement speed
func (p *Player) handlePlayerInput(delta float32) {
	if rl.IsKeyDown(p.keyMap.KeyLeft) {
		p.destRect.X -= p.stats.Speed * delta
	}

	if rl.IsKeyDown(p.keyMap.KeyRight) {
		p.destRect.X += p.stats.Speed * delta
	}

	if rl.IsKeyDown(p.keyMap.KeyUp) {
		p.destRect.Y -= p.stats.Speed * delta
	}

	if rl.IsKeyDown(p.keyMap.KeyDown) {
		p.destRect.Y += p.stats.Speed * delta
	}

	if rl.IsKeyDown(p.keyMap.KeyFire) {
		p.fire()
	}

	if rl.IsKeyDown(rl.KeyLeftShift) && rl.IsKeyPressed(rl.KeyEnd) {
		p.TakeDamage(10000000)
	}

	if rl.IsKeyDown(rl.KeyLeftShift) && rl.IsKeyPressed(rl.KeyPageUp) {
		events.GetEventManagerInstance().Emit(events_data.AddMessage, events_data.AddMessageData{
			Message: "Sample Data" + time.Now().String(),
			Timer:   2.0,
			Color:   rl.DarkGreen,
		})
	}

	if rl.IsKeyDown(rl.KeyLeftShift) && rl.IsKeyPressed(rl.KeyEqual) {
		events.GetEventManagerInstance().Emit(events_data.StatUpgradeEvent, events_data.StatUpgradeData{
			StatType: "speed",
			Value:    10,
		})
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
	if p.fireDelay < 10 {
		return
	}

	// TODO: Change projectile to have a damage stat that is set by the player's damage stat when grabbed from the pool
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
		p.stats.Health -= dmg
		p.damaged = true
	}
	if p.stats.Health <= 0 {
		p.stats.Health = 0
		events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{NewState: systems_data.GameOver})
	}
}

func (p *Player) handleStatUpgrade(event events.Event) {
	if data, ok := event.Data.(events_data.StatUpgradeData); ok {
		switch data.StatType {
		case "speed":
			p.stats.Speed += data.Value
		case "damage":
			p.stats.Damage += data.Value
		case "health":
			// Increase current and max health
			p.stats.Health += int(data.Value)
			p.stats.MaxHealth += int(data.Value)
		case "fireRate":
			// Lower value means faster firing
			p.stats.FireRate -= data.Value
			// Set a minimum value to prevent too fast firing
			if p.stats.FireRate < 5 {
				p.stats.FireRate = 5
			}
		case "size":
			// Increase player size
			p.stats.Size += data.Value
			p.destRect.Width = float32(p.texture.Width) * p.stats.Size
			p.destRect.Height = float32(p.texture.Height) * p.stats.Size
		}

		// Emit a message to notify the player about the upgrade
		events.GetEventManagerInstance().Emit(events_data.AddMessage, events_data.AddMessageData{
			Message: data.StatType + " upgraded!",
			Timer:   2.0,
			Color:   rl.Green,
		})
	}
}

func (p *Player) AddScore(score int) {
	p.score += score
}

func (p *Player) DestroyProjectile(proj *Projectile) {
	p.projPool.Return(proj)
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
	return p.stats.Health
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
