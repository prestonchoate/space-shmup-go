package entities

import "github.com/google/uuid"

type GameEntity interface {
	Activate(bool)
	Draw()
	Update()
	GetID() uuid.UUID
}
