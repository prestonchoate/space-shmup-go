package entities

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type UIManager struct {
	player       *Player
	enemyManager *EnemyManager
	health       int
	score        int
	enemyCount   int
}

func CreateUIManager(p *Player, em *EnemyManager) *UIManager {
	return &UIManager{
		player:       p,
		enemyManager: em,
	}
}

func (u *UIManager) Update() {
	u.health = u.player.health
	u.score = u.player.score
	u.enemyCount = len(u.enemyManager.enemies.activePool)
}

func (u *UIManager) Draw() {

	rl.DrawText(
		fmt.Sprintf("%10v%03d", "Health: ", u.health),
		int32(rl.GetScreenWidth()-200),
		0,
		20,
		rl.Gray)

	rl.DrawText(
		fmt.Sprintf("%10v%06d", "Score: ", u.score),
		int32(rl.GetScreenWidth()-200),
		30,
		20,
		rl.Gray)

	rl.DrawText(
		fmt.Sprintf("%10v%06d", "Enemies: ", u.enemyCount),
		int32(rl.GetScreenWidth()-200),
		60,
		20,
		rl.Gray)
}

func (u *UIManager) GetID() uuid.UUID {
	return uuid.Nil
}
