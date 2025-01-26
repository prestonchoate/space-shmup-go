package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	entities "github.com/prestonchoate/space-shmup/Entities"
)

var PlayerTexture rl.Texture2D
var BackgroundTexture rl.Texture2D
var ProjectileTexture rl.Texture2D
var EnemyTextures []rl.Texture2D
var gameManagerInstance *GameManager

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
	Player          *entities.Player
	EnemyManager    *entities.EnemyManager
}

func GetGameMangerInstance() *GameManager {
	if gameManagerInstance == nil {
		gameManagerInstance = createGameManager()
	}

	return gameManagerInstance
}

func (gm *GameManager) Update() {
	gm.handleButtonInputs()
	if gm.state == Playing && gm.assetsLoaded == false {
		gm.GameSetup()
	}

	for _, entity := range gm.entities {
		entity.Activate(gm.state == Playing)
		entity.Update()
	}

	if gm.Player.GetHealth() <= 0 {
		gm.state = GameOver
	}

	if gm.EnemyManager.GetEnemyCount() == 0 {
		gm.level++
		gm.EnemyManager.SpawnNewEnemies(gm.level)
	}

	if gm.state == Playing {
		gm.collisionSystem.Update()
	}

	gm.uiSystem.Update(UIUpdate{
		health:     gm.Player.GetHealth(),
		score:      gm.Player.GetScore(),
		enemyCount: gm.EnemyManager.GetEnemyCount(),
		state:      gm.state,
	})

	if gm.uiSystem.TransitionReady {
		// TODO: check if transitioning to the Playing state. If so call the GameSetup() function
		if gm.uiSystem.TransitionState == Playing {
			gm.GameSetup()
		}
		gm.state = gm.uiSystem.TransitionState
		gm.uiSystem.TransitionReady = false
	}
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
	gm.entities = append(gm.entities, gm.Player)
	gm.entities = append(gm.entities, gm.EnemyManager)

	gm.level = 0
	gm.state = Playing
}

func (gm *GameManager) loadAssets() {
	PlayerTexture = rl.LoadTexture("assets/player/1B.png")
	BackgroundTexture = rl.LoadTexture("assets/background.jpg")
	ProjectileTexture = rl.LoadTexture("assets/projectile/rocket.png")

	EnemyTextures = append(EnemyTextures, rl.LoadTexture("assets/enemies/Emissary.png"))
	rl.SetTextureWrap(BackgroundTexture, rl.WrapRepeat)
	rl.SetTextureWrap(ProjectileTexture, rl.WrapRepeat)
}

func createGameManager() *GameManager {
	gm := &GameManager{
		entities: []entities.GameEntity{},
	}
	gm.loadAssets()
	gm.assetsLoaded = true

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

	gm.Player = entities.CreatePlayer(&PlayerTexture, &ProjectileTexture, gm.currentSettings.keys)
	gm.EnemyManager = entities.CreateEnemyManager(EnemyTextures)

	gm.collisionSystem = CreateCollisionManager(gm.Player, gm.EnemyManager)
	gm.uiSystem = CreateUIManager(gm.Player, gm.EnemyManager)
	return gm
}
