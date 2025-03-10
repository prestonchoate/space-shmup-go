package systems

import (
	"fmt"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	entities "github.com/prestonchoate/space-shmup/Entities"
	assets "github.com/prestonchoate/space-shmup/Systems/Assets"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
	events "github.com/prestonchoate/space-shmup/Systems/Events"
	events_data "github.com/prestonchoate/space-shmup/Systems/Events/Data"
	"github.com/prestonchoate/space-shmup/Systems/saveManager"
)

var PlayerTexture rl.Texture2D
var BackgroundTexture rl.Texture2D
var ProjectileTexture rl.Texture2D
var EnemyTextures []rl.Texture2D
var gameManagerInstance *GameManager

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
	currentSettings systems_data.GameSettings
	Player          *entities.Player
	EnemyManager    *entities.EnemyManager
	backgroundMusic *rl.Sound // TODO: Change this to rl.Music when the audio issue is resolved
	returnState     systems_data.GameState
}

func GetGameMangerInstance() *GameManager {
	if gameManagerInstance == nil {
		gameManagerInstance = createGameManager()
	}

	return gameManagerInstance
}

func (gm *GameManager) Update() {

	dt := rl.GetFrameTime()

	//rl.UpdateMusicStream(*gm.backgroundMusic)

	if !rl.IsSoundPlaying(*gm.backgroundMusic) {
		rl.PlaySound(*gm.backgroundMusic)
		rl.SetSoundVolume(*gm.backgroundMusic, gm.currentSettings.MusicVolume)
	}

	gm.handleButtonInputs()
	if gm.state == systems_data.Playing && gm.assetsLoaded == false {
		gm.GameSetup()
	}

	for _, entity := range gm.entities {
		entity.Activate(gm.state == systems_data.Playing)
		entity.Update(dt)
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
	gm.Player = entities.CreatePlayer(&PlayerTexture, &ProjectileTexture, gm.currentSettings.Keys)
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
	gm.entities = append(gm.entities, entities.CreateBackground())
	gm.entities = append(gm.entities, gm.Player)
	gm.entities = append(gm.entities, gm.EnemyManager)

	gm.level = 0
	gm.state = systems_data.Playing
}

func (gm *GameManager) loadAssets() {
	// TODO: allow all entities to grab their own textures from the AssetManager so that there is a single place to manage that entity's data
	am := assets.GetAssetManagerInstance()

	pt, ok := am.GetTexture("assets/sprites/player/1B.png")
	if !ok {
		log.Fatal("Player texture not available in asset manager")
	}

	projTex, ok := am.GetTexture("assets/sprites/projectile/rocket.png")
	if !ok {
		log.Fatal("Projectile texture not available in asset manager")
	}

	et, ok := am.GetTexture("assets/sprites/enemies/Emissary.png")
	if !ok {
		log.Fatal("Enemy texture not available in asset manager")
	}

	PlayerTexture = pt
	ProjectileTexture = projTex

	EnemyTextures = append(EnemyTextures, et)
	rl.SetTextureWrap(BackgroundTexture, rl.WrapRepeat)
	rl.SetTextureWrap(ProjectileTexture, rl.WrapRepeat)
}

func createGameManager() *GameManager {
	gm := &GameManager{
		state:    systems_data.Start,
		entities: []entities.GameEntity{},
	}
	gm.loadAssets()
	gm.assetsLoaded = true
	bg := entities.CreateBackground()
	gm.entities = append(gm.entities, bg)
	bgMusic, ok := assets.GetAssetManagerInstance().GetSound("assets/music/deep-space-barrier-121195.ogg")
	if ok {
		gm.backgroundMusic = &bgMusic
		//gm.backgroundMusic.Looping = true
		//gm.backgroundMusic.Stream.SampleRate = 44100
		//rl.PlayMusicStream(*gm.backgroundMusic)

		rl.PlaySound(*gm.backgroundMusic)

	} else {
		log.Println("Game Manager: Failed to load bg music from asset manager")
	}

	gm.currentSettings = saveManager.GetInstance().Data.Settings

	rl.SetSoundVolume(*gm.backgroundMusic, gm.currentSettings.MusicVolume)

	gm.Player = entities.CreatePlayer(&PlayerTexture, &ProjectileTexture, gm.currentSettings.Keys)
	gm.EnemyManager = entities.CreateEnemyManager(EnemyTextures)

	gm.collisionSystem = CreateCollisionManager(gm.Player, gm.EnemyManager)
	gm.uiSystem = CreateUIManager()

	events.GetEventManagerInstance().Subscribe(events_data.ChangeGameState, gm.handleChangeStateEvent)
	events.GetEventManagerInstance().Subscribe(events_data.GameSettingsUpdated, gm.handleUpdatedSettings)
	events.GetEventManagerInstance().Subscribe(events_data.ReturnGameState, gm.handleReturnStateEvent)
	return gm
}

func (gm *GameManager) ShouldExit() bool {
	return gm.state == systems_data.Exit
}

func (gm *GameManager) handleChangeStateEvent(e events.Event) {
	if data, ok := e.Data.(events_data.ChangeStateData); ok {
		fmt.Println("Processing change state event. Changing to: ", data.NewState)
		gm.returnState = gm.state
		gm.state = data.NewState
		if gm.state == systems_data.Playing {
			gm.GameSetup()
		} else if gm.state == systems_data.Restart {
			gm.Reset()
		}
	}
}

func (gm *GameManager) handleReturnStateEvent(e events.Event) {
	if gm.returnState == systems_data.None {
		return
	}
	gm.state = gm.returnState
	gm.returnState = systems_data.None
}

func (gm *GameManager) handleUpdatedSettings(e events.Event) {
	if data, ok := e.Data.(events_data.UpdateSettingsData); ok {
		fmt.Println("Game Manager: Updating settings")
		gm.currentSettings = data.NewSettings
		rl.SetTargetFPS(data.NewSettings.TargetFPS)
		if rl.IsWindowFullscreen() != data.NewSettings.Fullscreen {
			rl.ToggleFullscreen()
		}
		if !rl.IsWindowFullscreen() {
			rl.SetWindowSize(data.NewSettings.ScreenWidth, data.NewSettings.ScreenHeight)
		}

		rl.SetSoundVolume(*gm.backgroundMusic, data.NewSettings.MusicVolume)
	}
}
