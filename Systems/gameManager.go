package systems

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	entities "github.com/prestonchoate/space-shmup/Entities"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
	events "github.com/prestonchoate/space-shmup/Systems/Events"
	events_data "github.com/prestonchoate/space-shmup/Systems/Events/Data"
)

var PlayerTexture rl.Texture2D
var BackgroundTexture rl.Texture2D
var ProjectileTexture rl.Texture2D
var EnemyTextures []rl.Texture2D
var gameManagerInstance *GameManager


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
	state           systems_data.GameState
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
	if gm.state == systems_data.Playing && gm.assetsLoaded == false {
		gm.GameSetup()
	}

	for _, entity := range gm.entities {
		entity.Activate(gm.state == systems_data.Playing)
		entity.Update()
	}

	// TODO: Handle this in EnemyManager and emit an event when all enemies are cleared
	if gm.EnemyManager.GetEnemyCount() == 0 {
		gm.level++
		gm.EnemyManager.SpawnNewEnemies(gm.level)
	}

	if gm.state == systems_data.Playing {
		gm.collisionSystem.Update()
	}

	gm.uiSystem.Update(UIUpdate{
		health:     gm.Player.GetHealth(),
		score:      gm.Player.GetScore(),
		enemyCount: gm.EnemyManager.GetEnemyCount(),
		state:      gm.state,
	})
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
	case systems_data.Playing:
		if rl.IsKeyPressed(rl.KeyEscape) {
			gm.state = systems_data.Paused
		}
	case systems_data.Paused:
		if rl.IsKeyPressed(rl.KeyEscape) {
			gm.state = systems_data.Playing
		}
	}
}

func (gm *GameManager) Reset() {
	gm.Player = entities.CreatePlayer(&PlayerTexture, &ProjectileTexture, gm.currentSettings.keys)
	gm.EnemyManager = entities.CreateEnemyManager(EnemyTextures)
	gm.collisionSystem.player = gm.Player
	gm.collisionSystem.enemyManager = gm.EnemyManager
	gm.GameSetup()
}

func (gm *GameManager) GameSetup() {
	rl.SetExitKey(0)
	gm.windowHeight = rl.GetScreenHeight()
	gm.windowWidth = rl.GetScreenWidth()
	gm.fps = rl.GetFPS()
	gm.state = systems_data.Loading
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
	gm.state = systems_data.Playing
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
	gm.uiSystem = CreateUIManager()

	events.GetEventManagerInstance().Subscribe("changeState", gm.handleChangeStateEvent)
	return gm
}

func (gm *GameManager) handleChangeStateEvent(e events.Event) {
	if data, ok := e.Data.(events_data.ChangeStateData); ok {
		fmt.Println("Processing change state event. Changing to: ", data.NewState)
		gm.state = data.NewState
		if gm.state == systems_data.Playing {
			gm.GameSetup()
		} else if gm.state == systems_data.Restart {
			gm.Reset()
		}
	}
}
