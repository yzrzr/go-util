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
	"sync"
)

// newSafeListIterator 创建安全的迭代器
// 一旦调用了迭代器的 Remove 方法来修改 List 集合，m 锁将会升级为写锁
// 安全的迭代器使用完之后必须主动调用 Close 进行关闭，释放锁。
func newSafeListIterator[E comparable](list ListIterator[E], m *sync.RWMutex) ListIterator[E] {
	m.RLock()
	return &safeListIterator[E]{
		ListIterator: list,
		m:            m,
	}
}

type safeListIterator[E comparable] struct {
	ListIterator[E]
	*sync.RWMutex
	m   *sync.RWMutex
	mup bool
}

func (l *safeListIterator[E]) HasNext() bool {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()
	return l.ListIterator.HasNext()
}

func (l *safeListIterator[E]) Next() (e E, err error) {
	l.RWMutex.Lock()
	defer l.RWMutex.Unlock()
	return l.ListIterator.Next()
}

func (l *safeListIterator[E]) Remove() error {
	l.RWMutex.Lock()
	defer l.RWMutex.Unlock()
	if !l.mup {
		l.m.RUnlock()
		l.m.Lock()
		l.mup = true
	}
	return l.ListIterator.Remove()
}

func (l *safeListIterator[E]) ForEachRemaining(action Consumer[E]) error {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()
	return l.ListIterator.ForEachRemaining(action)
}

func (l *safeListIterator[E]) HasPrevious() bool {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()
	return l.ListIterator.HasPrevious()
}

func (l *safeListIterator[E]) Previous() (e E, err error) {
	l.RWMutex.Lock()
	defer l.RWMutex.Unlock()
	return l.ListIterator.Previous()
}

func (l *safeListIterator[E]) NextIndex() int {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()
	return l.ListIterator.NextIndex()
}

func (l *safeListIterator[E]) PreviousIndex() int {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()
	return l.ListIterator.PreviousIndex()
}

func (l *safeListIterator[E]) Close() {
	l.RWMutex.Lock()
	defer l.RWMutex.Unlock()
	l.ListIterator.Close()
	if l.mup {
		l.m.RLock()
	} else {
		l.m.RUnlock()
	}
}
