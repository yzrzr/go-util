package pool

import (
	"testing"
)

func TestNewPool(t *testing.T) {
	p := NewPool[int](func() int {
		return 10
	})
	if p.Get() != 10 {
		t.Errorf("Get() = %v, want 10", p.Get())
	}
}
