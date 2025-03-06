package entities

import "github.com/google/uuid"

type GameEntity interface {
	Activate(bool)
	Draw()
	Update(delta float32)
	GetID() uuid.UUID
}
