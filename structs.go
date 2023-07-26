package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
}

type Projectile struct {
	texturePath string
	speed       float32
}

type Enemy struct {
}

type EnemyManager struct {
	enemies []Enemy
}

type Background struct {
	texturePath string
	srcRect     rl.Rectangle
	destRect    rl.Rectangle
	origin      rl.Vector2
}

type GameManager struct {
	player       Player
	enemyManager EnemyManager
	background   Background
}
