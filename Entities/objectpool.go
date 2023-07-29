package entities

import (
	"github.com/google/uuid"
)

type CreatePoolObject func() GameEntity

type ObjectPool[T GameEntity] struct {
	activePool   map[uuid.UUID]T
	inactivePool []T
	createFn     CreatePoolObject
}

func (op *ObjectPool[T]) Get() T {
	if len(op.inactivePool) == 0 {
		obj := op.createFn().(T)
		op.activePool[obj.GetID()] = obj
		return obj
	}

	obj := op.inactivePool[len(op.inactivePool)-1]
	op.inactivePool = op.inactivePool[:len(op.inactivePool)-1]
	op.activePool[obj.GetID()] = obj
	return obj
}

func (op *ObjectPool[T]) Return(obj T) {
	delete(op.activePool, obj.GetID())
	op.inactivePool = append(op.inactivePool, obj)
}
