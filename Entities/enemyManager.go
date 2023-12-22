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
}

func (em *EnemyManager) Draw() {
	for _, e := range em.enemies.activePool {
		e.Draw()
	}
}

func (em *EnemyManager) Update() {
	for _, e := range em.enemies.activePool {
		if e.texture.ID <= 0 {
			fmt.Println("Enemy texture is nil, setting to random texture")
			randIndex := rand.Intn(len(em.enemyTextures))
			e.texture = em.enemyTextures[randIndex]
			e.srcRect = rl.NewRectangle(0.0, 0.0, float32(e.texture.Width), float32(e.texture.Height))
			startX := rand.Intn(rl.GetScreenWidth())
			startY := rl.GetRandomValue(-3000, -100)
			e.destRect = rl.NewRectangle(float32(startX), float32(startY), float32(e.texture.Width), float32(e.texture.Height))
		}
		e.Update()
	}
}

func (em *EnemyManager) SpawnNewEnemies(level int) {
	totalCount := level * em.enemyCount
	if len(em.enemies.activePool) < totalCount {
		newSpawns := totalCount - len(em.enemies.activePool)
		fmt.Printf("Spawning %d new enemies\n", newSpawns)
		for i := 0; i < newSpawns; i++ {
			_ = em.enemies.Get()
		}
	}
}

func (em *EnemyManager) GetID() uuid.UUID {
	return em.id
}

func (em *EnemyManager) DestroyEnemy(e *Enemy) {
	e.destRect.Y = float32(rl.GetRandomValue(-3000, -100))
	em.enemies.Return(e)
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
	}

	return em
}

func createEnemy() GameEntity {
	return &Enemy{
		id:    uuid.New(),
		speed: 2.0,
	}
}
