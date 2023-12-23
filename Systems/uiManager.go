package systems

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	entities "github.com/prestonchoate/space-shmup/Entities"
)

type UIManager struct {
	player       *entities.Player
	enemyManager *entities.EnemyManager
	health       int
	score        int
	enemyCount   int
}

func CreateUIManager(p *entities.Player, em *entities.EnemyManager) *UIManager {
	return &UIManager{
		player:       p,
		enemyManager: em,
	}
}

func (u *UIManager) Update() {
	u.health = u.player.GetHealth()
	u.score = u.player.GetScore()
	u.enemyCount = u.enemyManager.GetEnemyCount()
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
