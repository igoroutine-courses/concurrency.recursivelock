package config

import (
	"errors"
	"iter"
)

// ErrEffectiveValueNotFound
// Ошибка, возвращаемая, если ни в текущем узле, ни у предков не найдено значение
var ErrEffectiveValueNotFound = errors.New("effective value not found")

// Node представляет узел конфигурационного дерева.
// Использует обобщения (type parameter T), поддерживает потокобезопасность
type Node[T any] struct{}

// NewNode создает новый узел конфигурационного дерева.
// Если передан родитель, добавляет новый узел в список его детей
func NewNode[T any](parent *Node[T], value T) *Node[T] {
	panic("not implemented")
}

// ClearValue очищает локальное значение в текущем узле
func (n *Node[T]) ClearValue() {
	panic("not implemented")
}

// All возвращает все предыдущие узлы
func (n *Node[T]) All() iter.Seq[*Node[T]] {
	panic("not implemented")
}

// GetEffectiveValue возвращает эффективное значение текущего узла.
// Если локальное значение отсутствует, рекурсивно ищет значение у предков
func (n *Node[T]) GetEffectiveValue() (T, error) {
	panic("not implemented")
}
