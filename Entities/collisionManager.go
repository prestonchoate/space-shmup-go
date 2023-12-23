package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type CollisionManager struct {
	player       *Player
	enemyManager *EnemyManager
}

func CreateCollisionManager(p *Player, em *EnemyManager) *CollisionManager {
	return &CollisionManager{
		player:       p,
		enemyManager: em,
	}
}

func (c *CollisionManager) Update() {
	c.checkPlayerCollision()
	c.checkProjectileCollision()
}

func (c *CollisionManager) Draw() {
	// No Draw method needed
}

func (c *CollisionManager) GetID() uuid.UUID {
	return uuid.UUID{}
}

func (c *CollisionManager) checkPlayerCollision() {
	for _, e := range c.enemyManager.enemies.activePool {
		if checkCollisionRecs(c.player.destRect, e.destRect) {
			c.enemyManager.DestroyEnemy(e)
			c.player.Damage(e.damage)
		}
	}
}

func (c *CollisionManager) checkProjectileCollision() {
	for _, proj := range c.player.projPool.activePool {
		for _, e := range c.enemyManager.enemies.activePool {
			if checkCollisionRecs(proj.destRect, e.destRect) {
				c.enemyManager.DestroyEnemy(e)
				c.player.DestroyProjectile(proj)
				c.player.AddScore(e.score)
				break
			}
		}
	}
}

func checkCollisionRecs(rec1, rec2 rl.Rectangle) bool {
	return rl.CheckCollisionRecs(rec1, rec2)
}
