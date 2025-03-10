package entities

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	assets "github.com/prestonchoate/space-shmup/Systems/Assets"
	"github.com/prestonchoate/space-shmup/Systems/saveManager"
)

type Background struct {
	id       uuid.UUID
	texture  rl.Texture2D
	srcRect  rl.Rectangle
	destRect rl.Rectangle
	origin   rl.Vector2
}

func CreateBackground() *Background {
	am := assets.GetAssetManagerInstance()
	sm := saveManager.GetInstance()
	bt, ok := am.GetTexture("assets/sprites/backgrounds/background.jpg")
	if !ok {
		log.Fatal("Background textrue not available in asset manager")
	}
	bg := &Background{
		id:       uuid.New(),
		texture:  bt,
		srcRect:  rl.NewRectangle(0.0, 0.0, float32(bt.Width), float32(bt.Height)),
		destRect: rl.NewRectangle(0.0, 0.0, float32(sm.Data.Settings.ScreenWidth), float32(sm.Data.Settings.ScreenHeight)),
	}

	return bg
}

func (bg *Background) Draw() {
	rl.DrawTexturePro(bg.texture, bg.srcRect, bg.destRect, bg.origin, 0, rl.White)
}

func (bg *Background) Update(delta float32) {
	settings := saveManager.GetInstance().Data.Settings
	bg.destRect.Width = float32(settings.ScreenWidth)
	bg.destRect.Height = float32(settings.ScreenHeight)
	bg.srcRect.Y -= 100 * delta
}

func (bg *Background) GetID() uuid.UUID {
	return bg.id
}

func (bg *Background) Activate(active bool) {
	// Nothing to do here
}
