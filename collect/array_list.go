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
	"fmt"
	"github.com/yzrzr/go-util/constraints"
	"sort"
	"strings"
)

// NewArrayList Abstract Factory
func NewArrayList[E any](initialCapacity int, comparator constraints.EqualComparator[E]) List[E] {
	return &arrayList[E]{
		elementData: make([]E, initialCapacity),
		capacity:    initialCapacity,
		size:        0,
		comparator:  comparator,
	}
}

type arrayList[E any] struct {
	elementData []E
	size        int
	capacity    int
	zeroVal     E
	comparator  constraints.EqualComparator[E]
}

func (a *arrayList[E]) Size() int {
	return a.size
}

func (a *arrayList[E]) IsEmpty() bool {
	return a.size == 0
}

func (a *arrayList[E]) Contains(e E) bool {
	return a.IndexOf(e) >= 0
}

func (a *arrayList[E]) Iterator() Iterator[E] {
	return a.ListIterator()
}

func (a *arrayList[E]) ToArray() []E {
	if a.size == 0 {
		return nil
	}
	tmp := make([]E, a.size)
	copy(tmp, a.elementData)
	return tmp
}

func (a *arrayList[E]) Add(e E) bool {
	a.grow(a.size + 1)
	a.elementData[a.size] = e
	a.size++
	return true
}

func (a *arrayList[E]) Remove(e E) bool {
	return a.RemoveIfN(func(o E) bool {
		return a.comparator.Equal(e, o)
	}, 1) == 1
}

func (a *arrayList[E]) RemoveN(e E, n int) int {
	return a.RemoveIfN(func(o E) bool {
		return a.comparator.Equal(e, o)
	}, n)
}

func (a *arrayList[E]) ContainsAll(c Collection[E]) bool {
	itr := c.Iterator()
	for itr.HasNext() {
		if e, err := itr.Next(); err != nil || !a.Contains(e) {
			return false
		}
	}
	return true
}

func (a *arrayList[E]) AddAll(c Collection[E]) {
	arr := c.ToArray()
	num := len(arr)
	if num == 0 {
		return
	}
	a.grow(num + a.size)
	copy(a.elementData[a.size:], arr)
	a.size += num
}

func (a *arrayList[E]) RemoveAll(c Collection[E]) int {
	return a.RemoveIfN(func(e E) bool {
		return c.Contains(e)
	}, -1)
}

func (a *arrayList[E]) RemoveIf(filter Predicate[E]) int {
	return a.RemoveIfN(filter, -1)
}

func (a *arrayList[E]) RemoveIfN(filter Predicate[E], n int) int {
	if n == 0 {
		return 0
	} else if n == -1 {
		n = a.size
	}
	var cnt int
	indexSlice := make([]int, 0, 2)
	for i := 0; i < a.size; i++ {
		if filter(a.elementData[i]) {
			indexSlice = append(indexSlice, i)
			cnt++
			if cnt == n {
				break
			}
		}
	}
	if cnt == 0 {
		return 0
	}
	var s, l, r, size int
	s = indexSlice[0]
	for i := 0; i < cnt; {
		l = indexSlice[i] + 1
		for i < cnt-1 && indexSlice[i+1] == l {
			i++
			l++
		}
		if i == cnt-1 {
			r = a.size
		} else {
			r = indexSlice[i+1]
		}
		size = r - l
		copy(a.elementData[s:s+size], a.elementData[l:r])
		s = s + size
		i++
	}
	for i := s; i < a.size; i++ {
		a.elementData[i] = a.zeroVal
	}
	a.size -= cnt
	return cnt
}

func (a *arrayList[E]) RetainAll(c Collection[E]) int {
	return a.RemoveIfN(func(e E) bool {
		return !c.Contains(e)
	}, -1)
}

func (a *arrayList[E]) Clear() {
	for i := 0; i < a.size; i++ {
		a.elementData[i] = a.zeroVal
	}
	a.size = 0
}

func (a *arrayList[E]) Equals(c Collection[E]) bool {
	return equals[E](a, c)
}

func (a *arrayList[E]) ForEach(f Consumer[E]) error {
	var err error
	for _, v := range a.elementData {
		err = f(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *arrayList[E]) ReplaceAll(operator UnaryOperator[E]) {
	if operator == nil {
		return
	}
	for i := 0; i < a.size; i++ {
		a.elementData[i] = operator(a.elementData[i])
	}
}

func (a *arrayList[E]) Sort(less SortLess[E]) {
	sort.Slice(a.elementData[:a.size], func(i, j int) bool {
		return less(a.elementData[i], a.elementData[j])
	})
}

func (a *arrayList[E]) Get(index int) (E, error) {
	if err := a.rangeCheck(index); err != nil {
		return a.zeroVal, err
	}
	return a.elementData[index], nil
}

func (a *arrayList[E]) Set(index int, e E) (E, error) {
	old, err := a.Get(index)
	if err != nil {
		return a.zeroVal, err
	}
	a.elementData[index] = e
	return old, nil
}

func (a *arrayList[E]) AddAt(index int, e E) error {
	if err := a.rangeCheck(index); err != nil {
		return err
	}
	a.grow(a.size + 1)
	copy(a.elementData[index+1:], a.elementData[index:])
	a.elementData[index] = e
	a.size++
	return nil
}

func (a *arrayList[E]) RemoveAt(index int) (E, error) {
	old, err := a.Get(index)
	if err != nil {
		return a.zeroVal, err
	}
	copy(a.elementData[index:], a.elementData[index+1:])
	a.size--
	return old, nil
}

func (a *arrayList[E]) IndexOf(e E) int {
	for i := 0; i < a.size; i++ {
		if a.comparator.Equal(e, a.elementData[i]) {
			return i
		}
	}
	return -1
}

func (a *arrayList[E]) LastIndexOf(e E) int {
	for i := a.size - 1; i >= 0; i-- {
		if a.comparator.Equal(e, a.elementData[i]) {
			return i
		}
	}
	return -1
}

func (a *arrayList[E]) ListIterator() ListIterator[E] {
	return a.ListIteratorAt(0)
}

func (a *arrayList[E]) ListIteratorAt(index int) ListIterator[E] {
	return newListIterator[E](a, index)
}

func (a *arrayList[E]) SubList(fromIndex, toIndex int) List[E] {
	data := make([]E, toIndex-fromIndex)
	copy(data, a.elementData[fromIndex:toIndex])
	return &arrayList[E]{
		elementData: data,
		capacity:    cap(data),
		size:        len(data),
		comparator:  a.comparator,
	}
}

func (a *arrayList[E]) grow(minCapacity int) {
	if a.capacity >= minCapacity {
		return
	}
	newCapacity := a.capacity
	doubleCapacity := a.capacity << 1
	if minCapacity > doubleCapacity {
		newCapacity = minCapacity
	} else {
		if a.capacity < 256 {
			newCapacity = doubleCapacity
		} else {
			for 0 < newCapacity && newCapacity < minCapacity {
				newCapacity += (newCapacity + 3*256) / 4
			}
			if newCapacity <= 0 {
				newCapacity = minCapacity
			}
		}
	}
	tmp := make([]E, newCapacity)
	copy(tmp, a.elementData)
	a.elementData = tmp
	a.capacity = newCapacity
}

func (a *arrayList[E]) rangeCheck(index int) error {
	if index < 0 || index >= a.size {
		return fmt.Errorf("index out of range [%d] with length %d", index, a.size)
	}
	return nil
}

func (a *arrayList[E]) String() string {
	build := strings.Builder{}
	build.WriteByte('[')
	for i := 0; i < a.size; i++ {
		build.WriteString(fmt.Sprintf("%v", a.elementData[i]))
		if i < a.size-1 {
			build.WriteByte(' ')
		}
	}
	build.WriteByte(']')
	return build.String()
}

func (a *arrayList[E]) GetEqualComparator() constraints.EqualComparator[E] {
	return a.comparator
}

//---- sort.Interface -----

type sortList[E any] struct {
	data []E
	less SortLess[E]
	swap func(i, j int)
}

func (a sortList[E]) Len() int {
	return len(a.data)
}

func (a sortList[E]) Less(i, j int) bool {
	return a.less(a.data[i], a.data[j])
}

func (a sortList[E]) Swap(i, j int) {
	if a.swap != nil {
		a.swap(i, j)
	} else {
		a.data[i], a.data[j] = a.data[j], a.data[i]
	}
}
