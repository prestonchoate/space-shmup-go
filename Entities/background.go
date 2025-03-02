package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type Background struct {
	id       uuid.UUID
	texture  rl.Texture2D
	srcRect  rl.Rectangle
	destRect rl.Rectangle
	origin   rl.Vector2
}

func CreateBackground(tex *rl.Texture2D, width float32, height float32) *Background {
	bg := &Background{
		id:       uuid.New(),
		texture:  *(tex),
		srcRect:  rl.NewRectangle(0.0, 0.0, float32(tex.Width), float32(tex.Height)),
		destRect: rl.NewRectangle(0.0, 0.0, width, height),
	}

	return bg
}

func (bg *Background) Draw() {
	rl.DrawTexturePro(bg.texture, bg.srcRect, bg.destRect, bg.origin, 0, rl.White)
}

func (bg *Background) Update() {
	bg.srcRect.Y -= 1
}

func (bg *Background) GetID() uuid.UUID {
	return bg.id
}

func (bg *Background) Activate(active bool) {
	// Nothing to do here
}
