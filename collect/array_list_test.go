/*
 *
 * Copyright 2022 go-util authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package collect

import (
	"errors"
	"reflect"
	"testing"
)

func newArrayList[E comparable](list ...E) List[E] {
	arr := NewArrayList[E](len(list))
	//arr := NewLinkedList[E]()
	for _, v := range list {
		arr.Add(v)
	}
	return arr
}

func checkData[E comparable](t *testing.T, list List[E], data []E) {
	for k, v := range data {
		e, err := list.Get(k)
		if err != nil {
			t.Errorf("Get() error = %v, wantErr %v", err, nil)
		}
		if e != v {
			t.Errorf("index %d got = %v, want %v", k, e, v)
		}
	}
}

func checkDataIndex[E comparable](t *testing.T, list List[E], index int, want E) {
	e, err := list.Get(index)
	if err != nil {
		t.Errorf("Get() error = %v, wantErr %v", err, nil)
	}
	if e != want {
		t.Errorf("index %d got = %v, want %v", index, e, want)
	}
}

func Test_arrayList_Add(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	list := newArrayList[int](s...)
	checkData(t, list, s)
}

func Test_arrayList_AddAll(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	list1 := newArrayList[int]([]int{1, 2, 3, 4, 5}...)
	list2 := newArrayList[int]([]int{10, 9, 8, 7}...)
	list1.AddAll(list2)
	checkData(t, list1, s)
}

func Test_arrayList_AddAt(t *testing.T) {
	list := newArrayList[int](10, 20)
	type args[E comparable] struct {
		index int
		e     E
	}
	tests := []struct {
		name    string
		args    args[int]
		wantErr bool
	}{
		{"addAt-1", args[int]{2, 1}, true},
		{"addAt-2", args[int]{0, 20}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := list.AddAt(tt.args.index, tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddAt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				checkDataIndex[int](t, list, tt.args.index, tt.args.e)
			}
		})
	}
}

func Test_arrayList_Clear(t *testing.T) {
	list := newArrayList[int](10, 20)
	list.Clear()
	if list.Size() != 0 {
		t.Errorf("Clear() size = %d, want 0", list.Size())
	}
}

func Test_arrayList_Contains(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	list := newArrayList[int](s...)
	tests := []struct {
		name string
		args int
		want bool
	}{
		{"Contains-1", 1, true},
		{"Contains-2", 3, true},
		{"Contains-3", 7, true},
		{"Contains-4", 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.Contains(tt.args); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_ContainsAll(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	list := newArrayList[int](s...)
	tests := []struct {
		name string
		args Collection[int]
		want bool
	}{
		{"ContainsAll-1", newArrayList[int](1, 2, 7), true},
		{"ContainsAll-2", newArrayList[int](7, 10, 6), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.ContainsAll(tt.args); got != tt.want {
				t.Errorf("ContainsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_Equals(t *testing.T) {
	s := []int{5, 10, 9}
	list := newArrayList[int](s...)
	tests := []struct {
		name string
		args Collection[int]
		want bool
	}{
		{"Equals-1", newArrayList[int](1, 2, 7), false},
		{"Equals-2", newArrayList[int](5, 10, 9), true},
		{"Equals-3", newArrayList[int](10, 5, 9), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.Equals(tt.args); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_ForEach(t *testing.T) {
	s := []int{5, 10, 9}
	list := newArrayList[int](s...)
	tests := []struct {
		name    string
		args    Consumer[int]
		wantErr bool
	}{
		{"ForEach-1", func(e int) error { return nil }, false},
		{"ForEach-2", func(e int) error {
			if e%2 == 0 {
				return errors.New("not an odd number")
			}
			return nil
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := list.ForEach(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ForEach() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_arrayList_Get(t *testing.T) {
	s := []int{5, 10, 9}
	list := newArrayList[int](s...)
	tests := []struct {
		name    string
		arg     int
		want    int
		wantErr bool
	}{
		{"Get-1", 2, 9, false},
		{"Get-2", -1, 0, true},
		{"Get-3", 3, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := list.Get(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_IndexOf(t *testing.T) {
	s := []int{5, 5, 9, 6, 6}
	list := newArrayList[int](s...)
	tests := []struct {
		name string
		arg  int
		want int
	}{
		{"IndexOf-1", 5, 0},
		{"IndexOf-2", 6, 3},
		{"IndexOf-3", 8, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.IndexOf(tt.arg); got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		list List[int]
		want bool
	}{
		{"IsEmpty-1", newArrayList[int](), true},
		{"IsEmpty-2", newArrayList[int](1), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_LastIndexOf(t *testing.T) {
	s := []int{5, 5, 9, 6, 6}
	list := newArrayList[int](s...)
	tests := []struct {
		name string
		arg  int
		want int
	}{
		{"LastIndexOf-1", 5, 1},
		{"LastIndexOf-2", 6, 4},
		{"LastIndexOf-3", 8, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.LastIndexOf(tt.arg); got != tt.want {
				t.Errorf("LastIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_Remove(t *testing.T) {
	s := []int{5, 9, 5, 6, 6, 5}
	tests := []struct {
		name  string
		org   List[int]
		arg   int
		want1 bool
		want2 List[int]
	}{
		{"Remove-1", newArrayList[int](s...), 5, true, newArrayList[int](9, 5, 6, 6, 5)},
		{"Remove-2", newArrayList[int](s...), 9, true, newArrayList[int](5, 5, 6, 6, 5)},
		{"Remove-3", newArrayList[int](s...), 10, false, newArrayList[int](5, 9, 5, 6, 6, 5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.org.Remove(tt.arg); got != tt.want1 {
				t.Errorf("Remove() = %v, want %v", got, tt.want1)
			}
			if !tt.org.Equals(tt.want2) {
				t.Errorf("Remove() = %v, want %v", tt.org, tt.want2)
			}
		})
	}
}

func Test_arrayList_RemoveN(t *testing.T) {
	s := []int{5, 9, 5, 6, 6, 5}
	type args struct {
		e int
		n int
	}
	tests := []struct {
		name  string
		org   List[int]
		args  args
		want1 int
		want2 List[int]
	}{
		{"RemoveN-1", newArrayList[int](s...), args{5, 1}, 1, newArrayList[int](9, 5, 6, 6, 5)},
		{"RemoveN-2", newArrayList[int](s...), args{6, -1}, 2, newArrayList[int](5, 9, 5, 5)},
		{"RemoveN-3", newArrayList[int](s...), args{5, 8}, 3, newArrayList[int](9, 6, 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.org.RemoveN(tt.args.e, tt.args.n); got != tt.want1 {
				t.Errorf("Remove() = %v, want %v", got, tt.want1)
			}
			if !tt.org.Equals(tt.want2) {
				t.Errorf("Remove() = %v, want %v", tt.org, tt.want2)
			}
		})
	}
}

func Test_arrayList_RemoveAll(t *testing.T) {
	s := []int{5, 9, 5, 6, 6, 5}
	tests := []struct {
		name  string
		org   List[int]
		args  List[int]
		want1 int
		want2 List[int]
	}{
		{"RemoveAll-1", newArrayList[int](s...), newArrayList[int](5), 3, newArrayList[int](9, 6, 6)},
		{"RemoveAll-2", newArrayList[int](s...), newArrayList[int](5, 6), 5, newArrayList[int](9)},
		{"RemoveAll-3", newArrayList[int](s...), newArrayList[int](8), 0, newArrayList[int](s...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.org.RemoveAll(tt.args); got != tt.want1 {
				t.Errorf("RemoveAll() = %v, want %v", got, tt.want1)
			}
			if !tt.org.Equals(tt.want2) {
				t.Errorf("Remove() = %v, want %v", tt.org, tt.want2)
			}
		})
	}
}

func Test_arrayList_RemoveAt(t *testing.T) {
	list := newArrayList[int](5, 9, 5, 6, 6, 5)
	tests := []struct {
		name    string
		arg     int
		want    int
		wantErr bool
	}{
		{"RemoveAt-1", 1, 9, false},
		{"RemoveAt-2", 2, 6, false},
		{"RemoveAt-3", 10, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := list.RemoveAt(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveAt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_RemoveIf(t *testing.T) {
	tests := []struct {
		name  string
		list  List[int]
		arg   Predicate[int]
		want1 int
		want2 List[int]
	}{
		{"RemoveIf-1", newArrayList[int](5, 9, 5, 6, 6, 5),
			func(e int) bool { return e%2 == 0 }, 2, newArrayList[int](5, 9, 5, 5)},
		{"RemoveIf-2", newArrayList[int](5, 9, 5, 6, 6, 5),
			func(e int) bool { return e%2 == 1 }, 4, newArrayList[int](6, 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.RemoveIf(tt.arg); got != tt.want1 {
				t.Errorf("RemoveIf() = %v, want %v", got, tt.want1)
			}
			if !tt.list.Equals(tt.want2) {
				t.Errorf("RemoveIf() = %v, want %v", tt.list, tt.want2)
			}
		})
	}
}

func Test_arrayList_RemoveIfN(t *testing.T) {
	type args struct {
		filter Predicate[int]
		n      int
	}
	tests := []struct {
		name  string
		list  List[int]
		args  args
		want1 int
		want2 List[int]
	}{
		{"RemoveIfN-1", newArrayList[int](5, 9, 5, 6, 6, 5),
			args{func(e int) bool { return e%2 == 0 }, -1}, 2, newArrayList[int](5, 9, 5, 5)},
		{"RemoveIfN-2", newArrayList[int](5, 9, 5, 6, 6, 5),
			args{func(e int) bool { return e%2 == 1 }, 3}, 3, newArrayList[int](6, 6, 5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.RemoveIfN(tt.args.filter, tt.args.n); got != tt.want1 {
				t.Errorf("RemoveIf() = %v, want %v", got, tt.want1)
			}
			if !tt.list.Equals(tt.want2) {
				t.Errorf("RemoveIf() = %v, want %v", tt.list, tt.want2)
			}
		})
	}
}

func Test_arrayList_ReplaceAll(t *testing.T) {
	tests := []struct {
		name string
		list List[int]
		arg  UnaryOperator[int]
		want List[int]
	}{
		{"ReplaceAll-1", newArrayList[int](5, 9, 5, 6, 6, 5),
			func(e int) int { return e + 1 }, newArrayList[int](6, 10, 6, 7, 7, 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.list.ReplaceAll(tt.arg)
			if !tt.list.Equals(tt.want) {
				t.Errorf("ReplaceAll() = %v, want %v", tt.list, tt.want)
			}
		})
	}
}

func Test_arrayList_RetainAll(t *testing.T) {
	tests := []struct {
		name  string
		list  List[int]
		arg   List[int]
		want1 int
		want2 List[int]
	}{
		{"RetainAll-1", newArrayList[int](5, 9, 5, 6, 6, 5),
			newArrayList[int](6), 4, newArrayList[int](6, 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.RetainAll(tt.arg); got != tt.want1 {
				t.Errorf("RetainAll() = %v, want %v", got, tt.want1)
			}
			if !tt.list.Equals(tt.want2) {
				t.Errorf("RetainAll() = %v, want %v", tt.list, tt.want2)
			}
		})
	}
}

func Test_arrayList_Set(t *testing.T) {
	list := newArrayList[int](5, 9, 5, 6, 6, 5)
	type args struct {
		index int
		e     int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"Set-1", args{2, 10}, 5, false},
		{"Set-2", args{10, 10}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := list.Set(tt.args.index, tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_Size(t *testing.T) {
	list := newArrayList[int](5, 9, 5)
	if list.Size() != 3 {
		t.Errorf("Size() = %v, want %v", list.Size(), 3)
	}
	list.Add(10)
	if list.Size() != 4 {
		t.Errorf("Size() = %v, want %v", list.Size(), 4)
	}
	list.AddAll(newArrayList[int](6, 6, 9))
	if list.Size() != 7 {
		t.Errorf("Size() = %v, want %v", list.Size(), 7)
	}
	list.RemoveN(6, -1)
	if list.Size() != 5 {
		t.Errorf("Size() = %v, want %v", list.Size(), 5)
	}
}

func Test_arrayList_Sort(t *testing.T) {
	list := newArrayList[int](3, 9, 5)
	tests := []struct {
		name string
		arg  SortLess[int]
		want List[int]
	}{
		{"Sort-1", func(e1, e2 int) bool {
			return e1 < e2
		}, newArrayList[int](3, 5, 9)},
		{"Sort-2", func(e1, e2 int) bool {
			return e1 > e2
		}, newArrayList[int](9, 5, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list.Sort(tt.arg)
			if !list.Equals(tt.want) {
				t.Errorf("Sort() got = %v, want %v", list, tt.want)
			}
		})
	}
}

func Test_arrayList_SubList(t *testing.T) {
	list := newArrayList[int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	type args struct {
		fromIndex int
		toIndex   int
	}
	tests := []struct {
		name string
		args args
		want List[int]
	}{
		{"SubList-1", args{1, 3}, newArrayList[int](2, 3)},
		{"SubList-2", args{5, list.Size()}, newArrayList[int](6, 7, 8, 9, 10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.SubList(tt.args.fromIndex, tt.args.toIndex); !got.Equals(tt.want) {
				t.Errorf("SubList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayList_ToArray(t *testing.T) {
	list := newArrayList[int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	if got := list.ToArray(); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Errorf("ToArray() = %v, want %v", got, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	}
}
