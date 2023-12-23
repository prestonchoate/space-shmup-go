package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	entities "github.com/prestonchoate/space-shmup/Entities"
)

var PlayerTexture rl.Texture2D
var BackgroundTexture rl.Texture2D
var ProjectileTexture rl.Texture2D
var EnemyTextures []rl.Texture2D

type GameState int

const (
	Loading GameState = iota
	Paused
	Playing
	Menu
	GameOver
)

type GameManager struct {
	entities        []entities.GameEntity
	windowHeight    int
	windowWidth     int
	fps             int32
	level           int
	state           GameState
	uiSystem        *UIManager
	collisionSystem *CollisionManager
}

func (gm *GameManager) Update() {
	if rl.IsKeyPressed(rl.KeyEscape) {
		if gm.state == Paused {
			gm.state = Playing
		} else if gm.state == Playing {
			gm.state = Paused
		}
	}

	for _, entity := range gm.entities {
		entity.Activate(gm.state == Playing)
		entity.Update()
		if em, ok := entity.(*entities.EnemyManager); ok {
			if em.GetEnemyCount() == 0 {
				gm.level++
				em.SpawnNewEnemies(gm.level)
			}
		}
	}
	gm.collisionSystem.Update()
	gm.uiSystem.Update()
}

func (gm *GameManager) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	for _, entity := range gm.entities {
		entity.Draw()
	}

	gm.uiSystem.Draw()

	rl.DrawFPS(10, 10)

	if gm.state == Paused {
		rl.DrawText("PAUSED", int32(gm.windowWidth/2)-rl.MeasureText("PAUSED", 40)/2, int32(gm.windowHeight/2), 40, rl.Gray)
	}
	rl.EndDrawing()
}

func (gm *GameManager) GameSetup() {
	rl.SetExitKey(0)
	gm.windowHeight = rl.GetScreenHeight()
	gm.windowWidth = rl.GetScreenWidth()
	gm.fps = rl.GetFPS()
	loadAssets()

	gm.entities = []entities.GameEntity{}

	bg := entities.CreateBackground(&BackgroundTexture, float32(gm.windowWidth), float32(gm.windowHeight))
	gm.entities = append(gm.entities, bg)

	player := entities.CreatePlayer(&PlayerTexture, &ProjectileTexture)
	gm.entities = append(gm.entities, player)

	em := entities.CreateEnemyManager(EnemyTextures)
	gm.entities = append(gm.entities, em)

	gm.collisionSystem = CreateCollisionManager(player, em)
	gm.uiSystem = CreateUIManager(player, em)

	gm.state = Playing
}

func loadAssets() {
	PlayerTexture = rl.LoadTexture("assets/player/1B.png")
	BackgroundTexture = rl.LoadTexture("assets/background.jpg")
	ProjectileTexture = rl.LoadTexture("assets/projectile/rocket.png")

	EnemyTextures = append(EnemyTextures, rl.LoadTexture("assets/enemies/Emissary.png"))
	rl.SetTextureWrap(BackgroundTexture, rl.RL_TEXTURE_WRAP_REPEAT)
	rl.SetTextureWrap(ProjectileTexture, rl.RL_TEXTURE_WRAP_REPEAT)
}
