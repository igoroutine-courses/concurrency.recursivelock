package reentrant

import (
	"math/rand/v2"

	"github.com/stretchr/testify/require"

	"sync"
	"testing"
)

var _ sync.Locker = (*Mutex)(nil)

func TestMutexDoubleLock(t *testing.T) {
	t.Parallel()

	mx := New()

	mx.Lock()
	mx.Lock()
	mx.Unlock()
}

func TestUnlockFromNonOwner(t *testing.T) {
	t.Parallel()

	mx := New()
	require.PanicsWithError(t, ErrUnlockFromAnotherGoroutine.Error(), func() {
		mx.Unlock()
	})
}

func TestUnlockFromAnotherGoroutine(t *testing.T) {
	t.Parallel()

	mx := New()

	mx.Lock()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		require.PanicsWithError(t, ErrUnlockFromAnotherGoroutine.Error(), func() {
			mx.Unlock()
		})
	}()

	wg.Wait()
}

func TestMutualExclusion(t *testing.T) {
	t.Parallel()

	v := make(map[int]int)
	mx := New()

	wg := new(sync.WaitGroup)
	for range 1_000 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mx.Lock()
			v[rand.N[int](10e9)]++
			defer mx.Unlock()

		}()
	}

	wg.Wait()
}

func TestMutexPerformance(t *testing.T) {
	stdLibMx := testing.Benchmark(func(b *testing.B) {
		b.SetParallelism(10)

		v := make(map[int]int)
		mx := new(sync.Mutex)

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mx.Lock()
				v[rand.N[int](10e9)]++
				mx.Unlock()
			}
		})
	})

	myMx := testing.Benchmark(func(b *testing.B) {
		b.SetParallelism(10)

		v := make(map[int]int)
		mx := New()

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mx.Lock()
				v[rand.N[int](10e9)]++
				mx.Unlock()
			}
		})
	})

	require.LessOrEqual(t, float64(stdLibMx.NsPerOp())/float64(myMx.NsPerOp()), 4.0)
}
