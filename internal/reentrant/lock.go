package reentrant

import (
	"errors"
)

// ErrUnlockFromAnotherGoroutine - ошибка, если Unlock вызывает горутина, не владевшая мьютексом
var ErrUnlockFromAnotherGoroutine = errors.New("unlock from non-owner goroutine")

// New создаёт новый экземпляр рекурсивного (реентерабельного) мьютекса
func New() *Mutex {
	panic("not implemented")
}

// Mutex реализует реентрантный мьютекс
type Mutex struct{}

// Lock захватывает мьютекс. Если текущая горутина уже владеет им,
// повторный вызов безопасен (реентерабельность)
func (r *Mutex) Lock() {
	panic("not implemented")
}

// Unlock освобождает мьютекс. Только горутина-владелец может вызвать Unlock, если другая горутина
// попытается разблокировать - будет паника
func (r *Mutex) Unlock() {
	panic("not implemented")
}
