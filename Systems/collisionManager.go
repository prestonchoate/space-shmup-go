package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	entities "github.com/prestonchoate/space-shmup/Entities"
)

type CollisionManager struct {
	player       *entities.Player
	enemyManager *entities.EnemyManager
}

func CreateCollisionManager(p *entities.Player, em *entities.EnemyManager) *CollisionManager {
	return &CollisionManager{
		player:       p,
		enemyManager: em,
	}
}

func (c *CollisionManager) Update() {
	c.checkPlayerCollision()
	c.checkProjectileCollision()
}

func (c *CollisionManager) checkPlayerCollision() {
	for _, e := range c.enemyManager.GetEnemies() {
		if checkCollisionRecs(c.player.GetRect(), e.GetRect()) {
			c.enemyManager.DestroyEnemy(e)
			c.player.TakeDamage(e.GetDamage())
		}
	}
}

func (c *CollisionManager) checkProjectileCollision() {
	for _, proj := range c.player.GetProjeciles() {
		for _, e := range c.enemyManager.GetEnemies() {
			if checkCollisionRecs(proj.GetRect(), e.GetRect()) {
				c.enemyManager.DestroyEnemy(e)
				c.player.DestroyProjectile(proj)
				c.player.AddScore(e.GetScore())
				break
			}
		}
	}
}

func checkCollisionRecs(rec1, rec2 rl.Rectangle) bool {
	return rl.CheckCollisionRecs(rec1, rec2)
}
