package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type Enemy struct {
	id       uuid.UUID
	texture  rl.Texture2D
	speed    float32
	origin   rl.Vector2
	srcRect  rl.Rectangle
	destRect rl.Rectangle
	targetX  float32
}

func (e *Enemy) Draw() {
	rl.DrawTexturePro(e.texture, e.srcRect, e.destRect, e.origin, 0, rl.White)
	rl.DrawRectangleLines(int32(e.destRect.X), int32(e.destRect.Y), int32(e.destRect.Width), int32(e.destRect.Height), rl.Red)
}

func (e *Enemy) Update() {
	if e.targetX == 0.0 {
		e.targetX = float32(rl.GetRandomValue(0, int32(rl.GetScreenWidth())))
	}

	if e.destRect.X < e.targetX {
		e.destRect.X = e.destRect.X + e.speed
	}

	if e.destRect.X > e.targetX {
		e.destRect.X = e.destRect.X - e.speed
	}

	e.destRect.Y = e.destRect.Y + e.speed

	if e.destRect.Y >= float32(rl.GetScreenHeight()) {
		startY := rl.GetRandomValue(-300, -100)
		e.destRect.Y = float32(startY)
	}

	if e.destRect.X == e.targetX || e.destRect.X == e.targetX+1 || e.destRect.X == e.targetX-1 {
		e.targetX = float32(rl.GetRandomValue(0, int32(rl.GetScreenWidth())))
	}
}

func (e *Enemy) GetID() uuid.UUID {
	return e.id
}
