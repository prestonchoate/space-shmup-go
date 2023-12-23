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
	Start GameState = iota
	Loading
	Paused
	Playing
	Menu
	GameOver
)

type GameSettings struct {
	targetFPS    int32
	screenWidth  int
	screenHeight int
	fullscreen   bool
	keys         entities.InputMap
}

type GameManager struct {
	entities        []entities.GameEntity
	windowHeight    int
	windowWidth     int
	fps             int32
	level           int
	state           GameState
	uiSystem        *UIManager
	collisionSystem *CollisionManager
	assetsLoaded    bool
	currentSettings GameSettings
}

func (gm *GameManager) Update() {
	gm.handleButtonInputs()
	for _, entity := range gm.entities {
		entity.Activate(gm.state == Playing)
		entity.Update()
		if em, ok := entity.(*entities.EnemyManager); ok {
			if em.GetEnemyCount() == 0 {
				gm.level++
				em.SpawnNewEnemies(gm.level)
			}
		}
		if p, ok := entity.(*entities.Player); ok {
			if p.GetHealth() <= 0 {
				gm.state = GameOver
			}
		}
	}
	if gm.state == Playing {
		gm.collisionSystem.Update()
	}

	gm.uiSystem.Update(gm.state)
}

func (gm *GameManager) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	for _, entity := range gm.entities {
		entity.Draw()
	}
	gm.uiSystem.HandleGameStateRender(gm.state)
	rl.DrawFPS(10, 10)
	rl.EndDrawing()
}

func (gm *GameManager) handleButtonInputs() {
	switch gm.state {
	case Start:
		if rl.IsKeyPressed(rl.KeyEnter) {
			gm.GameSetup()
		}
	case Playing:
		if rl.IsKeyPressed(rl.KeyEscape) {
			gm.state = Paused
		}
	case Paused:
		if rl.IsKeyPressed(rl.KeyEscape) {
			gm.state = Playing
		}
	case GameOver:
		if rl.IsKeyPressed(rl.KeyEnter) {
			gm.GameSetup()
		}
	}
}

func (gm *GameManager) GameSetup() {
	rl.SetExitKey(0)
	gm.windowHeight = rl.GetScreenHeight()
	gm.windowWidth = rl.GetScreenWidth()
	gm.fps = rl.GetFPS()
	gm.state = Loading
	if !gm.assetsLoaded {
		gm.loadAssets()
		gm.assetsLoaded = true
	}

	gm.entities = []entities.GameEntity{}

	bg := entities.CreateBackground(&BackgroundTexture, float32(gm.windowWidth), float32(gm.windowHeight))
	gm.entities = append(gm.entities, bg)

	player := entities.CreatePlayer(&PlayerTexture, &ProjectileTexture, gm.currentSettings.keys)
	gm.entities = append(gm.entities, player)

	em := entities.CreateEnemyManager(EnemyTextures)
	gm.entities = append(gm.entities, em)

	gm.collisionSystem = CreateCollisionManager(player, em)
	gm.uiSystem = CreateUIManager(player, em)
	gm.level = 0
	gm.state = Playing
}

func (gm *GameManager) loadAssets() {
	PlayerTexture = rl.LoadTexture("assets/player/1B.png")
	BackgroundTexture = rl.LoadTexture("assets/background.jpg")
	ProjectileTexture = rl.LoadTexture("assets/projectile/rocket.png")

	EnemyTextures = append(EnemyTextures, rl.LoadTexture("assets/enemies/Emissary.png"))
	rl.SetTextureWrap(BackgroundTexture, rl.RL_TEXTURE_WRAP_REPEAT)
	rl.SetTextureWrap(ProjectileTexture, rl.RL_TEXTURE_WRAP_REPEAT)
}

func CreateGameManager() *GameManager {
	gm := &GameManager{
		entities: []entities.GameEntity{},
	}

	settings := GameSettings{
		targetFPS:    120,
		screenWidth:  rl.GetScreenWidth(),
		screenHeight: rl.GetScreenHeight(),
		fullscreen:   false,
		keys: entities.InputMap{
			KeyLeft:  rl.KeyA,
			KeyUp:    rl.KeyW,
			KeyRight: rl.KeyD,
			KeyDown:  rl.KeyS,
			KeyFire:  rl.KeySpace,
		},
	}

	gm.currentSettings = settings

	return gm
}
