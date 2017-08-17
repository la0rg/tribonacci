// Package tribo provides a tribonacci sequence generator.
// Numbers could be requested separately by their serial number.
//
// Results are presented by tribo.Result struct.
//
// The package supports caching with configurable cache size
// and is safe for concurrent usage.
//
// Example usage:
//
// 	tribo := tribo.New(100000)
//	result, _ := tribo.Get(10)
//	fmt.Println(result.Value) // Output: 44
package tribo

import (
	"context"
	"errors"
	"math/big"
	"sync"
)

// Tribo is a type that provides interface to get tribonacci numbers
// and keeps a cache of them.
type Tribo struct {
	mx             sync.Mutex
	cache          []*big.Int
	cacheSizeLimit int
}

// New creates new instance of Tribo struct with the given cache size limit.
// Size limit should be a positive integer more than 3 (default cache size),
// otherwise cache will never be updated.
func New(cacheSizeLimit int) *Tribo {
	return &Tribo{
		cache:          []*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(1)},
		cacheSizeLimit: cacheSizeLimit,
	}
}

// Get returns tribonacci sequence value by its serial number.
// N has to be a positive integer, otherwise an error would be returned.
func (t *Tribo) Get(ctx context.Context, n int) (*Result, error) {
	if n <= 0 {
		return nil, errors.New("N should be a positive integer")
	}

	// Number is in the cache
	value, ok := t.cacheValue(n)
	if ok {
		return value, nil
	}

	// Calculate number without blocking of the cache
	tempCache := t.cacheCopy()
	var n1, n2, n3 = lastThree(tempCache)

	compute := make(chan struct{}, 1)
	// Compute and cache numbers to the required one
	for l := len(tempCache); n > l; l++ {
		go func() {
			// Compute the next tribonacci number
			// Tn+3 = Tn+2 + Tn+1 + Tn => Tn = Tn-1 + Tn-2 + Tn-3
			n1.Add(n1, n2).Add(n1, n3)
			n1, n2, n3 = n2, n3, n1

			// Do not exceed memory usage by too large cache
			if l < t.cacheSizeLimit {
				tempCache = append(tempCache, big.NewInt(0).Set(n3))
			}
			compute <- struct{}{}
		}()

		// Waits for finish of the computation or cancellation of the context
		select {
		case <-compute:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	go t.appendToCache(tempCache)
	return &Result{n3.String()}, nil
}

func lastThree(s []*big.Int) (*big.Int, *big.Int, *big.Int) {
	l := len(s)
	if l < 3 {
		panic("Tribo struct was not properly initialized")
	}
	// big.Int requires deep copy
	var r1, r2, r3 big.Int
	return r1.Set(s[l-3]), r2.Set(s[l-2]), r3.Set(s[l-1])
}

// blocking operations

func (t *Tribo) cacheValue(n int) (*Result, bool) {
	t.mx.Lock()
	defer t.mx.Unlock()
	if n <= len(t.cache) {
		return &Result{t.cache[n-1].String()}, true
	}
	return nil, false
}

func (t *Tribo) cacheCopy() []*big.Int {
	t.mx.Lock()
	defer t.mx.Unlock()
	c := make([]*big.Int, len(t.cache))
	copy(c, t.cache)
	return c
}

func (t *Tribo) appendToCache(temp []*big.Int) {
	t.mx.Lock()
	defer t.mx.Unlock()
	if len(temp) > len(t.cache) {
		t.cache = append(t.cache, temp[len(t.cache):]...)
	}
}

func (t *Tribo) getCacheLength() int {
	t.mx.Lock()
	defer t.mx.Unlock()
	return len(t.cache)
}
