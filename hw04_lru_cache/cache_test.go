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

	t.Run("direct element push-out", func(t *testing.T) {
		c := NewCache(3)

		// add 5 elements
		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		wasInCache = c.Set("eee", 500)
		require.False(t, wasInCache)

		require.Equal(t, 3, c.Len())

		// check all elements
		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)

		val, ok = c.Get("eee")
		require.True(t, ok)
		require.Equal(t, 500, val)
	})

	t.Run("long-unused element push-out", func(t *testing.T) {
		c := NewCache(4)

		// add first 4 elements
		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		require.Equal(t, 4, c.Len())

		// change priority
		wasInCache = c.Set("aaa", 101)
		require.True(t, wasInCache)

		val, ok := c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		// add new elements
		wasInCache = c.Set("eee", 500)
		require.False(t, wasInCache)

		wasInCache = c.Set("fff", 600)
		require.False(t, wasInCache)

		// check all elements
		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 101, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("ddd")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("eee")
		require.True(t, ok)
		require.Equal(t, 500, val)

		val, ok = c.Get("fff")
		require.True(t, ok)
		require.Equal(t, 600, val)
	})
}

func TestConcurrentCache(t *testing.T) {
	t.Run("concurrency", func(t *testing.T) {
		numConcurrent := 9999
		c := NewCache(10_000)
		var wg sync.WaitGroup
		wg.Add(numConcurrent)

		for i := 1; i <= numConcurrent; i++ {
			go func(value int) {
				defer wg.Done()
				c.Set(Key(strconv.Itoa(value)), value)
			}(i)
		}

		wg.Wait()
		require.Equal(t, 9999, c.Len())
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	maxValue := 1_000_000
	go func() {
		defer wg.Done()
		for i := 0; i < maxValue; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < maxValue; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(maxValue))))
		}
	}()

	wg.Wait()
	require.Equal(t, 10, c.Len())
}
