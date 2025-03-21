package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type Enemy struct {
	id        uuid.UUID
	texture   rl.Texture2D
	speed     float32
	origin    rl.Vector2
	srcRect   rl.Rectangle
	destRect  rl.Rectangle
	targetX   float32
	damage    int
	score     int
	scoreTick int
	active    bool
	prevX     float32
}

func (e *Enemy) Draw() {
	rl.DrawTexturePro(e.texture, e.srcRect, e.destRect, e.origin, 0, rl.White)
}

func (e *Enemy) Update(delta float32) {
	e.prevX = e.destRect.X
	if e.targetX == 0.0 {
		e.targetX = float32(rl.GetRandomValue(0, int32(rl.GetScreenWidth())))
	}

	if e.destRect.X < e.targetX {
		e.destRect.X = e.destRect.X + e.speed*delta
	}

	if e.destRect.X > e.targetX {
		e.destRect.X = e.destRect.X - e.speed*delta
	}

	e.destRect.Y = e.destRect.Y + e.speed*delta

	if e.destRect.Y >= float32(rl.GetScreenHeight()) {
		startY := rl.GetRandomValue(-300, -100)
		e.destRect.Y = float32(startY)
	}

	// if the enemy has reached the target, get a new target
	if e.destRect.X == e.targetX {
		e.targetX = float32(rl.GetRandomValue(e.texture.Width, int32(rl.GetScreenWidth())-e.texture.Width))
	}

	if e.scoreTick > 0 {
		e.scoreTick--
	}

	if e.scoreTick <= 0 {
		e.score--
		e.scoreTick = 120
	}

	if e.prevX == e.destRect.X { // enemy is stuck
		//e.destRect.X += float32(rl.GetRandomValue(-1, 1) * int32(e.speed))
		e.targetX = float32(rl.GetRandomValue(e.texture.Width, int32(rl.GetScreenWidth())-e.texture.Width))
	}

}

func (e *Enemy) GetID() uuid.UUID {
	return e.id
}

func (e *Enemy) Activate(active bool) {
	e.active = active
}

func (e *Enemy) GetDamage() int {
	return e.damage
}

func (e *Enemy) GetScore() int {
	return e.score
}

func (e *Enemy) GetRect() rl.Rectangle {
	return e.destRect
}
