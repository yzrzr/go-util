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
	"errors"
	"math"
)

var (
	ErrNoSuchElement = errors.New("no such element")
	ErrIllegalState  = errors.New("illegal state")
	ErrIteratorClose = errors.New("iterator is close")
)

type ListIterator[E any] interface {
	Iterator[E]

	// HasPrevious 如果有上一个元素，则返回true
	HasPrevious() bool

	// Previous 返回迭代中的上一个元素
	Previous() (E, error)

	// NextIndex 返回迭代中的下一个元素的索引
	NextIndex() int

	// PreviousIndex 返回迭代中的上一个元素的索引
	PreviousIndex() int
}

func newListIterator[E any](list List[E], start int) ListIterator[E] {
	return &listIterator[E]{
		lastRet: -1,
		cursor:  start,
		list:    list,
	}
}

type listIterator[E any] struct {
	// cursor 当前指针指向位置，下一次调用 Next() 方法返回 cursor 位置的值，每次调用 Next() 方法后，cursor 的值都会加一
	// lastRet 上一次调用 Next() 方法是的 cursor 值
	cursor, lastRet int

	isClose bool
	list    List[E]
}

func (l *listIterator[E]) HasNext() bool {
	if l.isClose {
		return false
	}
	return l.cursor < l.list.Size()
}

func (l *listIterator[E]) Next() (e E, err error) {
	if l.isClose {
		err = ErrIteratorClose
		return
	}
	i := l.cursor
	if i >= l.list.Size() {
		err = ErrNoSuchElement
		return
	}
	e, err = l.list.Get(i)
	if err != nil {
		return
	}
	l.cursor = i + 1
	l.lastRet = i
	return
}

func (l *listIterator[E]) Remove() error {
	if l.isClose {
		return ErrIteratorClose
	}
	if l.lastRet < 0 {
		return ErrIllegalState
	}
	_, err := l.list.RemoveAt(l.lastRet)
	if err != nil {
		return err
	}
	l.cursor = l.lastRet
	l.lastRet = -1
	return nil
}

func (l *listIterator[E]) ForEachRemaining(action Consumer[E]) error {
	if l.isClose {
		return ErrIteratorClose
	}
	size := l.list.Size()
	var err error
	var e E
	for i := l.cursor; i < size; i++ {
		e, err = l.list.Get(i)
		if err != nil {
			return err
		}
		err = action(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *listIterator[E]) HasPrevious() bool {
	return l.cursor != 0 && l.isClose == false
}

func (l *listIterator[E]) Previous() (e E, err error) {
	if l.isClose {
		err = ErrIteratorClose
		return
	}
	i := l.lastRet
	if i < 0 {
		err = ErrNoSuchElement
		return
	}
	e, err = l.list.Get(l.lastRet)
	if err != nil {
		return
	}
	l.cursor = i
	l.lastRet = i - 1
	return
}

func (l *listIterator[E]) NextIndex() int {
	return l.cursor
}

func (l *listIterator[E]) PreviousIndex() int {
	return l.lastRet
}

func (l *listIterator[E]) Close() {
	l.isClose = true
	l.cursor = math.MaxInt
	l.lastRet = -1
}

func newLinkedListIterator[E any](list *linkedList[E], start int) ListIterator[E] {
	e, _ := list.getElement(start)
	return &linkedListIterator[E]{
		cursor:    e,
		nextIndex: start,
		list:      list,
	}
}

type linkedListIterator[E any] struct {
	// cursor 当前指针指向元素，下一次调用 Next() 方法返回 cursor 的值
	// lastRet 上一次调用 Next() 方法时的 cursor 的值
	cursor, lastRet *list.Element
	list            *linkedList[E]
	nextIndex       int
	isClose         bool
}

func (l *linkedListIterator[E]) HasNext() bool {
	return l.cursor != nil && l.isClose == false
}

func (l *linkedListIterator[E]) Next() (e E, err error) {
	if l.isClose {
		err = ErrIteratorClose
		return
	}
	cur := l.cursor
	if cur == nil {
		err = ErrNoSuchElement
		return
	}
	l.cursor = l.cursor.Next()
	l.lastRet = cur
	l.nextIndex++
	return cur.Value.(E), nil
}

func (l *linkedListIterator[E]) Remove() error {
	if l.isClose {
		return ErrIteratorClose
	}
	if l.lastRet == nil {
		return ErrIllegalState
	}
	l.list.removeElement(l.lastRet)
	l.lastRet = nil
	return nil
}

func (l *linkedListIterator[E]) ForEachRemaining(action Consumer[E]) error {
	if l.isClose {
		return ErrIteratorClose
	}
	var err error
	cur := l.cursor
	for cur != nil {
		err = action(cur.Value.(E))
		if err != nil {
			return err
		}
		cur = cur.Next()
	}
	return nil
}

func (l *linkedListIterator[E]) HasPrevious() bool {
	return l.nextIndex > 0 && l.isClose == false
}

func (l *linkedListIterator[E]) Previous() (e E, err error) {
	if l.isClose {
		err = ErrIteratorClose
		return
	}
	if !l.HasPrevious() {
		err = ErrNoSuchElement
		return
	}
	element := l.lastRet
	l.cursor = l.lastRet
	l.lastRet = l.cursor.Prev()
	l.nextIndex--
	return element.Value.(E), nil
}

func (l *linkedListIterator[E]) NextIndex() int {
	return l.nextIndex
}

func (l *linkedListIterator[E]) PreviousIndex() int {
	return l.nextIndex - 1
}

func (l *linkedListIterator[E]) Close() {
	l.isClose = true
	l.cursor = nil
	l.lastRet = nil
}
