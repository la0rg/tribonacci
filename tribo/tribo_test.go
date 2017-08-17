package tribo

import (
	"context"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func BenchmarkGetCachedValue(b *testing.B) {
	tribo := New(100000)
	for n := 0; n < b.N; n++ {
		tribo.Get(context.Background(), 10000)
	}
}

func BenchmarkGetValueOutOfCache(b *testing.B) {
	tribo := New(100000)
	for n := 0; n < b.N; n++ {
		tribo.Get(context.Background(), 1000000)
	}
}

func TestNew(t *testing.T) {
	tribo := New(10)
	if tribo.cacheSizeLimit != 10 {
		t.Errorf("cacheSizeLimit should contain 10, not %d", tribo.cacheSizeLimit)
	}
	if len(tribo.cache) != 3 {
		t.Errorf("Cache length should be 3, not %d", len(tribo.cache))
	}
}

func TestCacheLimiting(t *testing.T) {
	tribo := New(100)

	tribo.Get(context.Background(), 50)
	// Waits for cache to be updated in separate goroutine
	time.Sleep(2000)
	if l := tribo.getCacheLength(); l != 50 {
		t.Errorf("Cache length should be 50, not %d", l)
	}

	tribo.Get(context.Background(), 1000)
	time.Sleep(2000)
	if l := tribo.getCacheLength(); l > 100 {
		t.Errorf("Cache length should not be bigger then cacheSizeLimit; cacheLength: %d", l)
	}
}

func TestGet(t *testing.T) {
	tribo := New(20)
	ctx := context.Background()
	ctxTimeout, cn := context.WithTimeout(context.Background(), time.Second)
	defer cn()
	tests := []struct {
		input   int
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{-5, ctx, "", true},
		{0, ctx, "", true},
		{1, ctx, "0", false},
		{2, ctx, "0", false},
		{3, ctx, "1", false},
		{14, ctx, "504", false},
		{15, ctx, "927", false},
		{16, ctx, "1705", false},
		{4, ctx, "1", false},
		{5, ctx, "2", false},
		{6, ctx, "4", false},
		{7, ctx, "7", false},
		{8, ctx, "13", false},
		{23, ctx, "121415", false},
		{24, ctx, "223317", false},
		{25, ctx, "410744", false},
		{26, ctx, "755476", false},
		{27, ctx, "1389537", false},
		{28, ctx, "2555757", false},
		{29, ctx, "4700770", false},
		{30, ctx, "8646064", false},
		{9, ctx, "24", false},
		{10, ctx, "44", false},
		{11, ctx, "81", false},
		{12, ctx, "149", false},
		{13, ctx, "274", false},
		{17, ctx, "3136", false},
		{18, ctx, "5768", false},
		{19, ctx, "10609", false},
		{20, ctx, "19513", false},
		{21, ctx, "35890", false},
		{22, ctx, "66012", false},
		{1000000, ctxTimeout, "", true},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.input), func(t *testing.T) {
			got, err := tribo.Get(tt.ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tribo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value != tt.want {
				t.Errorf("Tribo.Get(%v) = %v, want %v", tt.input, got.Value, tt.want)
			}
		})
	}
}

func TestLastThree(t *testing.T) {
	var a, b, c = big.NewInt(1), big.NewInt(2), big.NewInt(3)
	var d, e, f = lastThree([]*big.Int{a, b, c, c, b, a})
	if d.Cmp(c) != 0 || e.Cmp(b) != 0 || f.Cmp(a) != 0 {
		t.Errorf("lastThree %v %v %v, want %v %v %v", d, e, f, c, b, a)
	}
	if d == c || e == b || f == a {
		t.Error("lastThree should return deep copy of the last three integer")
	}
}

func TestLastThreePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("lastThree did not panic with nil slice")
		}
	}()
	lastThree(nil)
}
