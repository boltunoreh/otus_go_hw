package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic by size", func(t *testing.T) {
		c := NewCache(5)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)
		c.Set("eee", 500)
		c.Set("fff", 600)
		c.Set("ggg", 700)
		c.Set("hhh", 800)

		_, ok := c.Get("aaa")
		require.False(t, ok)
		_, ok = c.Get("bbb")
		require.False(t, ok)
		_, ok = c.Get("ccc")
		require.False(t, ok)

		_, ok = c.Get("ddd")
		require.True(t, ok)
		_, ok = c.Get("eee")
		require.True(t, ok)
		_, ok = c.Get("fff")
		require.True(t, ok)
		_, ok = c.Get("ggg")
		require.True(t, ok)
		_, ok = c.Get("hhh")
		require.True(t, ok)
	})

	t.Run("purge logic by usage", func(t *testing.T) {
		c := NewCache(5)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)
		c.Set("eee", 500)

		_, _ = c.Get("aaa")
		_, _ = c.Get("aaa")
		_, _ = c.Get("aaa")
		_, _ = c.Get("ccc")
		_, _ = c.Get("ccc")

		c.Set("fff", 600)
		c.Set("ggg", 700)
		c.Set("hhh", 800)

		_, ok := c.Get("aaa")
		require.True(t, ok)
		_, ok = c.Get("ccc")
		require.True(t, ok)
		_, ok = c.Get("fff")
		require.True(t, ok)
		_, ok = c.Get("ggg")
		require.True(t, ok)
		_, ok = c.Get("hhh")
		require.True(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
		_, ok = c.Get("ddd")
		require.False(t, ok)
		_, ok = c.Get("eee")
		require.False(t, ok)
	})

	t.Run("purge logic by usage", func(t *testing.T) {
		c := NewCache(5)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)
		c.Set("eee", 500)

		c.Clear()

		_, ok := c.Get("aaa")
		require.False(t, ok)
		_, ok = c.Get("bbb")
		require.False(t, ok)
		_, ok = c.Get("ccc")
		require.False(t, ok)
		_, ok = c.Get("ddd")
		require.False(t, ok)
		_, ok = c.Get("eee")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
