package collect

import (
	"errors"
	"reflect"
	"slices"
	"testing"
)

func Test_setIterator_next(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	it := SetOf(s...).Iterator()
	var s2 []int
	for _, v := range s {
		if !it.HasNext() {
			t.Error("HasNext() = false, want true")
			return
		}
		next, err := it.Next()
		if err != nil {
			t.Errorf("Next() = %v, want %v", next, v)
		}
		s2 = append(s2, next)
	}
	if it.HasNext() {
		t.Error("HasNext() = true, want false")
	}
	_, err := it.Next()
	if !errors.Is(err, ErrNoSuchElement) {
		t.Error("Next() = nil, want ErrNoSuchElement")
	}
	slices.Sort(s)
	slices.Sort(s2)
	if !reflect.DeepEqual(s, s2) {
		t.Errorf("set = %v, want %v", s, s)
	}
}

func Test_setIterator_remove(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	list := SetOf(s...)
	it := list.Iterator()
	err := it.Remove()
	if !errors.Is(err, ErrIllegalState) {
		t.Errorf("Remove() = %v, want ErrNoSuchElement", err)
	}
	for it.HasNext() {
		v, err := it.Next()
		if err != nil {
			t.Errorf("Next() = %v, want nil", err)
		}
		if v%2 == 0 {
			err = it.Remove()
			if err != nil {
				t.Errorf("Remove() = %v, want nil", err)
			}
		}
	}
	_, err = it.Next()
	if !errors.Is(err, ErrNoSuchElement) {
		t.Error("Next() = nil, want ErrNoSuchElement")
	}
	setArr := list.ToArray()
	slices.Sort(setArr)
	if !reflect.DeepEqual(setArr, []int{1, 3, 5, 7, 9}) {
		t.Errorf("list() = %v, want %v", list, []int{1, 3, 5, 7, 9})
	}
}

func Test_setIterator_close(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	list := SetOf(s...)
	it := list.Iterator()
	it.Close()
	if it.HasNext() {
		t.Error("HasNext() = true, want false")
	}
	_, err := it.Next()
	if !errors.Is(err, ErrIteratorClose) {
		t.Error("Next() = nil, want ErrIteratorClose")
	}
	err = it.ForEachRemaining(func(e int) error {
		return nil
	})
	if !errors.Is(err, ErrIteratorClose) {
		t.Error("ForEachRemaining() = nil, want ErrIteratorClose")
	}
	if !errors.Is(it.Remove(), ErrIteratorClose) {
		t.Error("Remove() = nil, want ErrIteratorClose")
	}
}

func Test_setIterator_foreach(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	list := SetOf(s...)
	it := list.Iterator()
	set2 := NewSet[int]()
	err := it.ForEachRemaining(func(e int) error {
		set2.Add(e * 2)
		return nil
	})
	if err != nil {
		t.Errorf("ForEachRemaining() = %v, want nil", err)
	}
	setArr := set2.ToArray()
	slices.Sort(setArr)
	if !reflect.DeepEqual(setArr, []int{2, 4, 6, 8, 10, 14, 16, 18, 20}) {
		t.Errorf("set2() = %v, want %v", set2, []int{2, 4, 6, 8, 10, 14, 16, 18, 20})
	}
	err = it.ForEachRemaining(func(e int) error {
		return errors.New("xx")
	})
	if err == nil {
		t.Errorf("ForEachRemaining() = %v, want error", err)
	}
}
