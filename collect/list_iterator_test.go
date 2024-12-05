package collect

import (
	"errors"
	"reflect"
	"testing"
)

func Test_listIterator_next(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	slen := len(s)
	for _, c := range configList {
		it := newArrayList(c, s...).ListIterator()
		for i, v := range s {
			if !it.HasNext() {
				t.Error("HasNext() = false, want true")
				return
			}
			if it.NextIndex() != i {
				t.Errorf("NextIndex() = %v, want %v", it.NextIndex(), i)
			}
			next, err := it.Next()
			if err != nil || v != next {
				t.Errorf("Next() = %v, want %v", next, v)
			}
		}
		if it.HasNext() {
			t.Error("HasNext() = true, want false")
		}
		_, err := it.Next()
		if !errors.Is(err, ErrNoSuchElement) {
			t.Error("Next() = nil, want ErrNoSuchElement")
		}
		for i := slen - 1; i >= 0; i-- {
			if !it.HasPrevious() {
				t.Error("HasPrevious() = false, want true")
			}
			if it.PreviousIndex() != i {
				t.Errorf("PreviousIndex() = %v, want %v", it.PreviousIndex(), i)
			}
			prev, err := it.Previous()
			if err != nil || s[i] != prev {
				t.Errorf("Previous() = %v, want %v", prev, s[i])
			}
		}
		if it.HasPrevious() {
			t.Error("HasPrevious() = true, want false")
		}
		_, err = it.Previous()
		if !errors.Is(err, ErrNoSuchElement) {
			t.Error("Previous() = nil, want ErrNoSuchElement")
		}
		v1, _ := it.Next()
		v2, _ := it.Previous()
		if v1 != v2 {
			t.Errorf("Next() = %v, Previous() = %v, want equal", v1, v2)
		}
		_, _ = it.Next()
		v1, _ = it.Next()
		v2, _ = it.Previous()
		if v1 != v2 {
			t.Errorf("Next() = %v, Previous() = %v, want equal", v1, v2)
		}
		if v1 != 2 {
			t.Errorf("Next() = %v, Previous() = %v, want equal", v1, v2)
		}
		if c.Safe {
			it.Close()
		}
	}
}

func Test_listIterator_remove(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	for _, c := range configList {
		list := newArrayList(c, s...)
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
		if c.Safe {
			it.Close()
		}
		if !reflect.DeepEqual(list.ToArray(), []int{1, 3, 5, 9, 7}) {
			t.Errorf("list() = %v, want %v", list, []int{1, 3, 5, 9, 7})
		}
	}
}

func Test_listIterator_close(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	for _, c := range configList {
		list := newArrayList(c, s...)
		it := list.ListIterator()
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
		if it.HasPrevious() {
			t.Error("HasPrevious() = true, want false")
		}
		_, err = it.Previous()
		if !errors.Is(err, ErrIteratorClose) {
			t.Error("Previous() = nil, want ErrIteratorClose")
		}
	}
}

func Test_listIterator_foreach(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	var endError = errors.New("end error")
	for _, c := range configList {
		list := newArrayList(c, s...)
		it := list.Iterator()
		list2 := newArrayList[int](DefaultListConfig)
		err := it.ForEachRemaining(func(e int) error {
			list2.Add(e * 2)
			return nil
		})
		if err != nil {
			t.Errorf("ForEachRemaining() = %v, want nil", err)
		}
		if !reflect.DeepEqual(list2.ToArray(), []int{2, 4, 6, 8, 10, 20, 18, 16, 14}) {
			t.Errorf("list2() = %v, want %v", list2, []int{2, 4, 6, 8, 10, 20, 18, 16, 14})
		}
		list3 := newArrayList[int](DefaultListConfig)
		err = it.ForEachRemaining(func(e int) error {
			if e >= 10 {
				return endError
			}
			list3.Add(e * 2)
			return nil
		})
		if !errors.Is(err, endError) {
			t.Errorf("ForEachRemaining() = %v, want %v", err, endError)
		}
		if !reflect.DeepEqual(list3.ToArray(), []int{2, 4, 6, 8, 10}) {
			t.Errorf("list3() = %v, want %v", list2, []int{2, 4, 6, 8, 10})
		}
		if c.Safe {
			it.Close()
		}
	}
}

func Test_listIteratorAt(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 10, 9, 8, 7}
	slen := len(s)
	for _, c := range configList {
		it := newArrayList(c, s...).ListIteratorAt(2)
		for i := 2; i < slen; i++ {
			if !it.HasNext() {
				t.Error("HasNext() = false, want true")
				return
			}
			if it.NextIndex() != i {
				t.Errorf("NextIndex() = %v, want %v", it.NextIndex(), i)
			}
			next, err := it.Next()
			if err != nil || s[i] != next {
				t.Errorf("Next() = %v, want %v", next, s[i])
			}
		}
		if it.HasNext() {
			t.Error("HasNext() = true, want false")
		}
		_, err := it.Next()
		if !errors.Is(err, ErrNoSuchElement) {
			t.Error("Next() = nil, want ErrNoSuchElement")
		}
		for i := slen - 1; i >= 0; i-- {
			if !it.HasPrevious() {
				t.Error("HasPrevious() = false, want true")
			}
			if it.PreviousIndex() != i {
				t.Errorf("PreviousIndex() = %v, want %v", it.PreviousIndex(), i)
			}
			prev, err := it.Previous()
			if err != nil || s[i] != prev {
				t.Errorf("Previous() = %v, want %v", prev, s[i])
			}
		}
		if it.HasPrevious() {
			t.Error("HasPrevious() = true, want false")
		}
		_, err = it.Previous()
		if !errors.Is(err, ErrNoSuchElement) {
			t.Error("Previous() = nil, want ErrNoSuchElement")
		}
		v1, _ := it.Next()
		v2, _ := it.Previous()
		if v1 != v2 {
			t.Errorf("Next() = %v, Previous() = %v, want equal", v1, v2)
		}
		_, _ = it.Next()
		v1, _ = it.Next()
		v2, _ = it.Previous()
		if v1 != v2 {
			t.Errorf("Next() = %v, Previous() = %v, want equal", v1, v2)
		}
		if v1 != 2 {
			t.Errorf("Next() = %v, Previous() = %v, want equal", v1, v2)
		}
		if c.Safe {
			it.Close()
		}
	}
}
