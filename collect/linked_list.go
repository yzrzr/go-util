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
	"container/list"
	"fmt"
	"sort"
	"strings"
)

func NewLinkedList[E comparable]() List[E] {
	return &linkedList[E]{
		list: list.New(),
	}
}

type linkedList[E comparable] struct {
	list    *list.List
	zeroVal E
}

func (l *linkedList[E]) Size() int {
	return l.list.Len()
}

func (l *linkedList[E]) IsEmpty() bool {
	return l.list.Len() == 0
}

func (l *linkedList[E]) Contains(e E) bool {
	return l.IndexOf(e) >= 0
}

func (l *linkedList[E]) Iterator() Iterator[E] {
	return l.ListIterator()
}

func (l *linkedList[E]) ToArray() []E {
	size := l.Size()
	if size == 0 {
		return nil
	}
	res := make([]E, size)
	for e, i := l.list.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		res[i] = e.Value.(E)
	}
	return res
}

func (l *linkedList[E]) Add(e E) bool {
	l.list.PushBack(e)
	return true
}

func (l *linkedList[E]) Remove(e E) bool {
	return l.RemoveIfN(func(o E) bool {
		return o == e
	}, 1) == 1
}

func (l *linkedList[E]) RemoveN(e E, n int) int {
	return l.RemoveIfN(func(o E) bool {
		return o == e
	}, n)
}

func (l *linkedList[E]) ContainsAll(c Collection[E]) bool {
	itr := c.Iterator()
	for itr.HasNext() {
		if e, err := itr.Next(); err != nil || !l.Contains(e) {
			return false
		}
	}
	return true
}

func (l *linkedList[E]) AddAll(c Collection[E]) {
	c.ForEach(func(e E) error {
		l.list.PushBack(e)
		return nil
	})
}

func (l *linkedList[E]) RemoveAll(c Collection[E]) int {
	return l.RemoveIfN(func(e E) bool {
		return c.Contains(e)
	}, -1)
}

func (l *linkedList[E]) RemoveIf(filter Predicate[E]) int {
	return l.RemoveIfN(filter, -1)
}

func (l *linkedList[E]) RemoveIfN(filter Predicate[E], n int) int {
	if n == 0 {
		return 0
	} else if n == -1 {
		n = l.Size()
	}
	var cnt int
	for cur, next := l.list.Front(), l.list.Front(); cur != nil; cur = next {
		next = cur.Next()
		if filter(cur.Value.(E)) {
			l.removeElement(cur)
			cnt++
			if cnt == n {
				break
			}
		}
	}
	return cnt
}

func (l *linkedList[E]) RetainAll(c Collection[E]) int {
	return l.RemoveIfN(func(e E) bool {
		return !c.Contains(e)
	}, -1)
}

func (l *linkedList[E]) Clear() {
	l.RemoveIfN(func(e E) bool {
		return true
	}, -1)
}

func (l *linkedList[E]) Equals(c Collection[E]) bool {
	return equals[E](l, c)
}

func (l *linkedList[E]) ForEach(f Consumer[E]) error {
	var err error
	for cur := l.list.Front(); cur != nil; cur = cur.Next() {
		if err = f(cur.Value.(E)); err != nil {
			return err
		}
	}
	return nil
}

func (l *linkedList[E]) ReplaceAll(operator UnaryOperator[E]) {
	for cur := l.list.Front(); cur != nil; cur = cur.Next() {
		cur.Value = operator(cur.Value.(E))
	}
}

func (l *linkedList[E]) Sort(less SortLess[E]) {
	data := make([]*list.Element, l.Size())
	for cur, i := l.list.Front(), 0; cur != nil; cur, i = cur.Next(), i+1 {
		data[i] = cur
	}
	sl := sortList[*list.Element]{
		data: data,
		less: func(e1, e2 *list.Element) bool {
			return less(e1.Value.(E), e2.Value.(E))
		},
		swap: func(i, j int) {
			data[i].Value, data[j].Value = data[j].Value, data[i].Value
		},
	}
	sort.Sort(sl)
	for cur, i := l.list.Front(), 0; cur != nil; cur, i = cur.Next(), i+1 {
		cur.Value = data[i].Value
	}
}

func (l *linkedList[E]) Get(index int) (E, error) {
	element, err := l.getElement(index)
	if err != nil {
		return l.zeroVal, err
	}
	return element.Value.(E), nil
}

func (l *linkedList[E]) Set(index int, e E) (E, error) {
	element, err := l.getElement(index)
	if err != nil {
		return l.zeroVal, err
	}
	old := element.Value.(E)
	element.Value = e
	return old, nil
}

func (l *linkedList[E]) AddAt(index int, e E) error {
	element, err := l.getElement(index)
	if err != nil {
		return err
	}
	l.list.InsertBefore(e, element)
	return nil
}

func (l *linkedList[E]) RemoveAt(index int) (E, error) {
	element, err := l.getElement(index)
	if err != nil {
		return l.zeroVal, err
	}
	old := element.Value.(E)
	l.removeElement(element)
	return old, nil
}

func (l *linkedList[E]) IndexOf(e E) int {
	for element, index := l.list.Front(), 0; element != nil; element, index = element.Next(), index+1 {
		if element.Value == e {
			return index
		}
	}
	return -1
}

func (l *linkedList[E]) LastIndexOf(e E) int {
	for element, index := l.list.Back(), l.Size()-1; element != nil; element, index = element.Prev(), index-1 {
		if element.Value == e {
			return index
		}
	}
	return -1
}

func (l *linkedList[E]) ListIterator() ListIterator[E] {
	return l.ListIteratorAt(0)
}

func (l *linkedList[E]) ListIteratorAt(index int) ListIterator[E] {
	return newLinkedListIterator[E](l, index)
}

func (l *linkedList[E]) SubList(fromIndex, toIndex int) List[E] {
	newLinkedList := NewLinkedList[E]()
	itr := l.ListIteratorAt(fromIndex)
	for itr.NextIndex() < toIndex && itr.HasNext() {
		e, _ := itr.Next()
		newLinkedList.Add(e)
	}
	return newLinkedList
}

func (l *linkedList[E]) String() string {
	build := strings.Builder{}
	build.WriteByte('[')
	for cur := l.list.Front(); cur != nil; cur = cur.Next() {
		build.WriteString(fmt.Sprintf("%v", cur.Value))
		if cur.Next() != nil {
			build.WriteByte(' ')
		}
	}
	build.WriteByte(']')
	return build.String()
}

func (l *linkedList[E]) getElement(index int) (*list.Element, error) {
	if err := l.rangeCheck(index); err != nil {
		return nil, err
	}
	if index < (l.Size() >> 1) {
		cur := l.list.Front()
		for i := 0; i < index; i++ {
			cur = cur.Next()
		}
		return cur, nil
	}
	cur := l.list.Back()
	for i := l.Size() - 1; i > index; i-- {
		cur = cur.Prev()
	}
	return cur, nil
}

func (l *linkedList[E]) removeElement(e *list.Element) {
	l.list.Remove(e)
}

func (l *linkedList[E]) rangeCheck(index int) error {
	if index < 0 || index >= l.Size() {
		return fmt.Errorf("index out of range [%d] with length %d", index, l.Size())
	}
	return nil
}
