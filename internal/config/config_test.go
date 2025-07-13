package config

import (
	"iter"
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetEffectiveValue(t *testing.T) {
	t.Parallel()

	root := NewNode[string](nil, "root_value")
	child := NewNode[string](root, "child_value")
	grandchild := NewNode[string](child, "grand_value")

	values := collect(t, grandchild.All())
	slices.Reverse(values)

	require.Equal(t, []string{"root_value", "child_value", "grand_value"}, values)

	val, err := grandchild.GetEffectiveValue()
	require.NoError(t, err)
	require.Equal(t, "grand_value", val)

	grandchild.ClearValue()

	val, err = grandchild.GetEffectiveValue()
	require.NoError(t, err)
	require.Equal(t, "child_value", val)

	child.ClearValue()

	val, err = grandchild.GetEffectiveValue()
	require.NoError(t, err)
	require.Equal(t, "root_value", val)

	root.ClearValue()

	_, err = grandchild.GetEffectiveValue()
	require.ErrorIs(t, err, ErrEffectiveValueNotFound)
	require.Equal(t, "root_value", val)
}

func TestThreadSafety(t *testing.T) {
	t.Parallel()

	root := NewNode[int](nil, 100)
	child := atomic.Pointer[Node[int]]{}
	child.Store(NewNode[int](root, 42))

	child.Load().ClearValue()

	wg := new(sync.WaitGroup)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(rand.N[time.Duration](10 * time.Millisecond))

			NewNode(child.Load(), i)
			NewNode(child.Load(), i)
			grand := NewNode(child.Load(), i)

			val, err := grand.GetEffectiveValue()

			require.NoError(t, err)
			require.Equal(t, i, val)

			child.Load().ClearValue()
			child.Store(grand)

			collect(t, child.Load().All())
		}()
	}

	wg.Wait()
}

func collect[T any](t *testing.T, iterator iter.Seq[*Node[T]]) []T {
	t.Helper()

	s := make([]T, 0)
	for v := range iterator {
		value, err := v.GetEffectiveValue()
		require.NoError(t, err)

		s = append(s, value)
	}

	return s
}
