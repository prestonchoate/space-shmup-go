package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Projectile struct {
	texture rl.Texture2D
	speed   float32
}

type Enemy struct {
}

type EnemyManager struct {
	enemies []Enemy
}
