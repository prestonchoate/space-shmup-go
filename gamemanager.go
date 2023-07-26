package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameEntity interface {
	Draw()
	Update()
}

type GameManager struct {
	entities []GameEntity
}

func (gm *GameManager) UpdateEntities() {
	for _, entity := range gm.entities {
		entity.Update()
	}
}

func (gm *GameManager) DrawEntities() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	for _, entity := range gm.entities {
		entity.Draw()
	}

	rl.DrawFPS(10, 10)
	rl.EndDrawing()
}

func (gm *GameManager) GameSetup() {
	gm.entities = []GameEntity{}

	bg := &Background{}
	bg.texture = rl.LoadTexture("assets/background.jpg")
	bg.Setup()
	gm.entities = append(gm.entities, bg)

	player := &Player{}
	player.texture = rl.LoadTexture("assets/player/1B.png")
	player.Setup()
	gm.entities = append(gm.entities, player)
}
