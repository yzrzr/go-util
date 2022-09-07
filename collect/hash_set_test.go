package collect

import (
	"reflect"
	"sort"
	"testing"
)

func getWantSet[E comparable](args ...E) Set[E] {
	m := make(map[E]struct{})
	for _, v := range args {
		m[v] = struct{}{}
	}
	return &hashSet[E]{
		data: m,
	}
}

func TestSetOf(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want Set[int]
	}{
		{"SetOf-1", []int{2, 2, 3, 4, 5, 4, 6}, getWantSet[int](2, 3, 4, 5, 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetOf(tt.args...); !got.Equals(tt.want) {
				t.Errorf("SetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashSet_Add(t *testing.T) {
	set := NewSet[int]()
	tests := []struct {
		name    string
		arg     int
		want    bool
		wantSet Set[int]
	}{
		{"Add-1", 10, true, getWantSet[int](10)},
		{"Add-2", 20, true, getWantSet[int](10, 20)},
		{"Add-3", 10, false, getWantSet[int](10, 20)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := set.Add(tt.arg)
			if got != tt.want {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
			if !set.Equals(tt.wantSet) {
				t.Errorf("Add() set = %v, want %v", set, tt.wantSet)
			}
		})
	}
}

func Test_hashSet_AddAll(t *testing.T) {
	set := NewSet[int]()
	tests := []struct {
		name    string
		arg     Collection[int]
		wantErr bool
		wantSet Set[int]
	}{
		{"AddAll-1", SetOf(1, 2), false, getWantSet[int](1, 2)},
		{"AddAll-2", SetOf(1, 2), false, getWantSet[int](1, 2)},
		{"AddAll-3", SetOf(1, 2, 3, 4), false, getWantSet[int](1, 2, 3, 4)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set.AddAll(tt.arg)
			if !set.Equals(tt.wantSet) {
				t.Errorf("AddAll() set = %v, want %v", set, tt.wantSet)
			}
		})
	}
}

func Test_hashSet_Clear(t *testing.T) {
	set := SetOf[int](1, 2, 3, 4)
	set.Clear()
	if !set.IsEmpty() {
		t.Errorf("Clear() isEmpty = false, wantErr true")
	}
}

func Test_hashSet_Contains(t *testing.T) {
	set := SetOf[int](1, 2, 3, 4)
	tests := []struct {
		name string
		arg  int
		want bool
	}{
		{"Contains-1", 2, true},
		{"Contains-2", 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.Contains(tt.arg); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashSet_ContainsAll(t *testing.T) {
	set := SetOf[int](1, 2, 3, 4)
	tests := []struct {
		name string
		arg  Collection[int]
		want bool
	}{
		{"ContainsAll-1", SetOf[int](1, 2, 3, 4), true},
		{"ContainsAll-2", SetOf[int](1, 2, 3, 4, 5), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.ContainsAll(tt.arg); got != tt.want {
				t.Errorf("ContainsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashSet_Equals(t *testing.T) {
	tests := []struct {
		name string
		set  Set[int]
		arg  Set[int]
		want bool
	}{
		{"Equals-1", SetOf[int](1, 2, 3, 4), SetOf[int](1, 2, 3, 4), true},
		{"Equals-2", SetOf[int](1, 2, 3, 4), SetOf[int](1, 2, 3), false},
		{"Equals-3", SetOf[int](), SetOf[int](), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Equals(tt.arg); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashSet_ForEach(t *testing.T) {
	total := 0
	tests := []struct {
		name      string
		set       Set[int]
		arg       Consumer[int]
		wantErr   bool
		wantTotal int
	}{
		{"ForEach-1", SetOf[int](1, 2, 3, 4, 4), func(e int) error { total += e; return nil }, false, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.set.ForEach(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("ForEach() error = %v, wantErr %v", err, tt.wantErr)
			}
			if total != tt.wantTotal {
				t.Errorf("ForEach() total = %v, wantTotal %v", total, tt.wantTotal)
			}
		})
	}
}

func Test_hashSet_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		set  Set[int]
		want bool
	}{
		{"IsEmpty-1", SetOf[int](1, 2, 3, 4), false},
		{"IsEmpty-2", SetOf[int](), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashSet_Remove(t *testing.T) {
	set := SetOf[int](10, 20, 30, 40)
	tests := []struct {
		name    string
		arg     int
		want    bool
		wantSet Set[int]
	}{
		{"Remove-1", 20, true, SetOf[int](10, 30, 40)},
		{"Remove-2", 20, false, SetOf[int](10, 30, 40)},
		{"Remove-3", 10, true, SetOf[int](30, 40)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.Remove(tt.arg); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
			if !set.Equals(tt.wantSet) {
				t.Errorf("Remove() set = %v, want %v", set, tt.wantSet)
			}
		})
	}
}

func Test_hashSet_RemoveAll(t *testing.T) {
	set := SetOf[int](50, 60, 70, 80, 10, 20, 30, 40)
	tests := []struct {
		name    string
		arg     Collection[int]
		want    int
		wantSet Set[int]
	}{
		{"RemoveAll-1", SetOf[int](10, 50), 2, SetOf[int](60, 70, 80, 20, 30, 40)},
		{"RemoveAll-2", SetOf[int](10, 50), 0, SetOf[int](60, 70, 80, 20, 30, 40)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.RemoveAll(tt.arg); got != tt.want {
				t.Errorf("RemoveAll() = %v, want %v", got, tt.want)
			}
			if !set.Equals(tt.wantSet) {
				t.Errorf("RemoveAll() set = %v, want %v", set, tt.wantSet)
			}
		})
	}
}

func Test_hashSet_RemoveIf(t *testing.T) {
	set := SetOf[int](1, 2, 3, 4, 5, 6, 7, 8)
	tests := []struct {
		name    string
		arg     Predicate[int]
		want    int
		wantSet Set[int]
	}{
		{"RemoveAll-1", func(e int) bool { return e%2 == 0 }, 4, SetOf[int](1, 3, 5, 7)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.RemoveIf(tt.arg); got != tt.want {
				t.Errorf("RemoveIf() = %v, want %v", got, tt.want)
			}
			if !set.Equals(tt.wantSet) {
				t.Errorf("RemoveIf() set = %v, want %v", set, tt.wantSet)
			}
		})
	}
}

func Test_hashSet_RetainAll(t *testing.T) {
	set := SetOf[int](1, 2, 3, 4, 5, 6, 7, 8)
	tests := []struct {
		name    string
		arg     Collection[int]
		want    int
		wantSet Set[int]
	}{
		{"RetainAll-1", SetOf[int](1, 2, 3, 4, 5, 6, 7, 8), 0, SetOf[int](1, 2, 3, 4, 5, 6, 7, 8)},
		{"RetainAll-2", SetOf[int](2, 6, 7), 5, SetOf[int](2, 6, 7)},
		{"RetainAll-3", SetOf[int](), 3, SetOf[int]()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.RetainAll(tt.arg); got != tt.want {
				t.Errorf("RetainAll() = %v, want %v", got, tt.want)
			}
			if !set.Equals(tt.wantSet) {
				t.Errorf("RetainAll() set = %v, want %v", set, tt.wantSet)
			}
		})
	}
}

func Test_hashSet_Size(t *testing.T) {
	set := SetOf[int](1, 2, 3, 4, 5, 6, 7, 8)
	if set.Size() != 8 {
		t.Errorf("Size() = %v, want %v", set.Size(), 8)
	}
	set.Clear()
	if set.Size() != 0 {
		t.Errorf("Size() = %v, want %v", set.Size(), 0)
	}
}

func Test_hashSet_ToArray(t *testing.T) {
	tests := []struct {
		name string
		set  Set[int]
		want []int
	}{
		{"ToArray-1", SetOf[int](1, 2, 3, 4, 5, 6, 7, 8), []int{1, 2, 3, 4, 5, 6, 7, 8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.set.ToArray()
			sort.Ints(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
