package entities

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type EnemyManager struct {
	enemies       ObjectPool[*Enemy]
	id            uuid.UUID
	enemyTextures []rl.Texture2D
	enemyCount    int
	active        bool
}

func (em *EnemyManager) Reset() {
	em.enemyCount = 10
	em.enemies.Reset()
}

func (em *EnemyManager) Draw() {
	if !em.active {
		return
	}
	for _, e := range em.enemies.activePool {
		e.Draw()
	}
}

func (em *EnemyManager) Update(delta float32) {
	if !em.active {
		return
	}
	for _, e := range em.enemies.activePool {
		if e.texture.ID <= 0 {
			fmt.Println("Enemy texture is nil, setting to random texture")
			randIndex := rand.Intn(len(em.enemyTextures))
			e.texture = em.enemyTextures[randIndex]
			e.srcRect = rl.NewRectangle(0.0, 0.0, float32(e.texture.Width), float32(e.texture.Height))
			startX := rl.GetRandomValue(e.texture.Width+10, int32(rl.GetScreenWidth())-e.texture.Width-10)
			startY := rl.GetRandomValue(-300, -100)
			e.destRect = rl.NewRectangle(float32(startX), float32(startY), float32(e.texture.Width), float32(e.texture.Height))
		}
		e.Update(delta)
	}
}

func (em *EnemyManager) SpawnNewEnemies(level int) {
	totalCount := level * em.enemyCount
	if len(em.enemies.activePool) < totalCount {
		newSpawns := totalCount - len(em.enemies.activePool)
		fmt.Printf("Spawning %d new enemies\n", newSpawns)
		for range newSpawns {
			_ = em.enemies.Get()
		}
	}
}

func (em *EnemyManager) GetID() uuid.UUID {
	return em.id
}

func (em *EnemyManager) Activate(active bool) {
	em.active = active
}

func (em *EnemyManager) DestroyEnemy(e *Enemy) {
	e.destRect.Y = float32(rl.GetRandomValue(-300, -100))
	e.score = int(rl.GetRandomValue(10, 250))
	e.scoreTick = 120
	e.speed = float32(rl.GetRandomValue(200, 300))
	em.enemies.Return(e)
}

func (em *EnemyManager) GetEnemies() map[uuid.UUID]*Enemy {
	return em.enemies.activePool
}

func (em *EnemyManager) GetEnemyCount() int {
	return len(em.enemies.activePool)
}

func CreateEnemyManager(textures []rl.Texture2D) *EnemyManager {
	em := &EnemyManager{
		id:            uuid.UUID{},
		enemyTextures: textures,
		enemyCount:    10,
		enemies: ObjectPool[*Enemy]{
			activePool:   make(map[uuid.UUID]*Enemy),
			inactivePool: make([]*Enemy, 0, 20000),
			createFn:     createEnemy,
		},
		active: true,
	}

	return em
}

func createEnemy() GameEntity {
	return &Enemy{
		id:        uuid.New(),
		speed:     float32(rl.GetRandomValue(200, 300)),
		score:     int(rl.GetRandomValue(10, 250)),
		damage:    10,
		scoreTick: 120,
	}
}
