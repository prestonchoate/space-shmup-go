package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

var PlayerTexture rl.Texture2D
var BackgroundTexture rl.Texture2D
var ProjectileTexture rl.Texture2D
var EnemyTextures []rl.Texture2D

type GameEntity interface {
	Draw()
	Update()
	GetID() uuid.UUID
}

type GameManager struct {
	entities     []GameEntity
	windowHeight int
	windowWidth  int
	fps          int32
}

func (gm *GameManager) UpdateEntities() {
	for _, entity := range gm.entities {
		entity.Update()
	}
}

func (gm *GameManager) DrawEntities() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	for _, entity := range gm.entities {
		entity.Draw()
	}

	rl.DrawFPS(10, 10)
	rl.EndDrawing()
}

func (gm *GameManager) GameSetup() {
	gm.windowHeight = rl.GetScreenHeight()
	gm.windowWidth = rl.GetScreenWidth()
	gm.fps = rl.GetFPS()
	loadAssets()

	gm.entities = []GameEntity{}

	bg := &Background{
		id:       uuid.New(),
		texture:  BackgroundTexture,
		srcRect:  rl.NewRectangle(0.0, 0.0, float32(BackgroundTexture.Width), float32(BackgroundTexture.Height)),
		destRect: rl.NewRectangle(0.0, 0.0, float32(gm.windowWidth), float32(gm.windowHeight)),
	}
	gm.entities = append(gm.entities, bg)

	player := &Player{
		id:       uuid.New(),
		texture:  PlayerTexture,
		speed:    2.5,
		srcRect:  rl.NewRectangle(0.0, 0.0, float32(PlayerTexture.Width), float32(PlayerTexture.Height)),
		destRect: rl.NewRectangle(float32(PlayerTexture.Width), float32(gm.windowWidth-int(PlayerTexture.Width)), float32(PlayerTexture.Width), float32(PlayerTexture.Height)),
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
	gm.entities = append(gm.entities, player)

	em := CreateEnemyManager(EnemyTextures)
	gm.entities = append(gm.entities, em)

	em.SpawnNewEnemies(1)
}

func loadAssets() {
	PlayerTexture = rl.LoadTexture("assets/player/1B.png")
	BackgroundTexture = rl.LoadTexture("assets/background.jpg")
	ProjectileTexture = rl.LoadTexture("assets/projectile/rocket.png")

	EnemyTextures = append(EnemyTextures, rl.LoadTexture("assets/enemies/Emissary.png"))
	rl.SetTextureWrap(BackgroundTexture, rl.RL_TEXTURE_WRAP_REPEAT)
	rl.SetTextureWrap(ProjectileTexture, rl.RL_TEXTURE_WRAP_REPEAT)
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
