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
	level        int
}

func (gm *GameManager) UpdateEntities() {
	for _, entity := range gm.entities {
		entity.Update()
		if em, ok := entity.(*EnemyManager); ok {
			if len(em.enemies.activePool) == 0 {
				gm.level++
				em.SpawnNewEnemies(gm.level)
			}
		}
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

	player := CreatePlayer(&PlayerTexture)
	gm.entities = append(gm.entities, player)

	em := CreateEnemyManager(EnemyTextures)
	gm.entities = append(gm.entities, em)

	cm := CreateCollisionManager(player, em)
	gm.entities = append(gm.entities, cm)

	ui := CreateUIManager(player, em)
	gm.entities = append(gm.entities, ui)
}

func loadAssets() {
	PlayerTexture = rl.LoadTexture("assets/player/1B.png")
	BackgroundTexture = rl.LoadTexture("assets/background.jpg")
	ProjectileTexture = rl.LoadTexture("assets/projectile/rocket.png")

	EnemyTextures = append(EnemyTextures, rl.LoadTexture("assets/enemies/Emissary.png"))
	rl.SetTextureWrap(BackgroundTexture, rl.RL_TEXTURE_WRAP_REPEAT)
	rl.SetTextureWrap(ProjectileTexture, rl.RL_TEXTURE_WRAP_REPEAT)
}
