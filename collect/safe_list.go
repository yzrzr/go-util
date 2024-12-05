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
	"sync"
)

type safeList[E any] struct {
	List[E]
	*sync.RWMutex
}

func NewSafeList[E any](list List[E]) List[E] {
	return &safeList[E]{
		List:    list,
		RWMutex: &sync.RWMutex{},
	}
}

func (a *safeList[E]) Size() int {
	a.RLock()
	defer a.RUnlock()
	return a.List.Size()
}

func (a *safeList[E]) IsEmpty() bool {
	a.RLock()
	defer a.RUnlock()
	return a.List.IsEmpty()
}

func (a *safeList[E]) Contains(e E) bool {
	a.RLock()
	defer a.RUnlock()
	return a.List.Contains(e)
}

func (a *safeList[E]) Iterator() Iterator[E] {
	return a.ListIterator()
}

func (a *safeList[E]) ToArray() []E {
	a.RLock()
	defer a.RUnlock()
	return a.List.ToArray()
}

func (a *safeList[E]) Add(e E) bool {
	a.Lock()
	defer a.Unlock()
	return a.List.Add(e)
}

func (a *safeList[E]) Remove(e E) bool {
	a.Lock()
	defer a.Unlock()
	return a.List.Remove(e)
}

func (a *safeList[E]) RemoveN(e E, n int) int {
	a.Lock()
	defer a.Unlock()
	return a.List.RemoveN(e, n)
}

func (a *safeList[E]) ContainsAll(c Collection[E]) bool {
	a.RLock()
	defer a.RUnlock()
	return a.List.ContainsAll(c)
}

func (a *safeList[E]) AddAll(c Collection[E]) {
	a.Lock()
	defer a.Unlock()
	a.List.AddAll(c)
}

func (a *safeList[E]) RemoveAll(c Collection[E]) int {
	a.Lock()
	defer a.Unlock()
	return a.List.RemoveAll(c)
}

func (a *safeList[E]) RemoveIf(filter Predicate[E]) int {
	a.Lock()
	defer a.Unlock()
	return a.List.RemoveIf(filter)
}

func (a *safeList[E]) RemoveIfN(filter Predicate[E], n int) int {
	a.Lock()
	defer a.Unlock()
	return a.List.RemoveIfN(filter, n)
}

func (a *safeList[E]) RetainAll(c Collection[E]) int {
	a.Lock()
	defer a.Unlock()
	return a.List.RetainAll(c)
}

func (a *safeList[E]) Clear() {
	a.Lock()
	defer a.Unlock()
	a.List.Clear()
}

func (a *safeList[E]) Equals(c Collection[E]) bool {
	a.RLock()
	defer a.RUnlock()
	return a.List.Equals(c)
}

func (a *safeList[E]) ForEach(f Consumer[E]) error {
	a.RLock()
	defer a.RUnlock()
	return a.List.ForEach(f)
}

func (a *safeList[E]) ReplaceAll(operator UnaryOperator[E]) {
	a.Lock()
	defer a.Unlock()
	a.List.ReplaceAll(operator)
}

func (a *safeList[E]) Sort(less SortLess[E]) {
	a.Lock()
	defer a.Unlock()
	a.List.Sort(less)
}

func (a *safeList[E]) Get(index int) (E, error) {
	a.RLock()
	defer a.RUnlock()
	return a.List.Get(index)
}

func (a *safeList[E]) Set(index int, e E) (E, error) {
	a.Lock()
	defer a.Unlock()
	return a.List.Set(index, e)
}

func (a *safeList[E]) AddAt(index int, e E) error {
	a.Lock()
	defer a.Unlock()
	return a.List.AddAt(index, e)
}

func (a *safeList[E]) RemoveAt(index int) (E, error) {
	a.Lock()
	defer a.Unlock()
	return a.List.RemoveAt(index)
}

func (a *safeList[E]) IndexOf(e E) int {
	a.RLock()
	defer a.RUnlock()
	return a.List.IndexOf(e)
}

func (a *safeList[E]) LastIndexOf(e E) int {
	a.RLock()
	defer a.RUnlock()
	return a.List.LastIndexOf(e)
}

func (a *safeList[E]) ListIterator() ListIterator[E] {
	a.RLock()
	defer a.RUnlock()
	return newSafeListIterator(a.List.ListIterator(), a.RWMutex)
}

func (a *safeList[E]) ListIteratorAt(index int) ListIterator[E] {
	a.RLock()
	defer a.RUnlock()
	return newSafeListIterator(a.List.ListIteratorAt(index), a.RWMutex)
}

func (a *safeList[E]) SubList(fromIndex, toIndex int) List[E] {
	a.RLock()
	defer a.RUnlock()
	return a.List.SubList(fromIndex, toIndex)
}

func (a *safeList[E]) String() string {
	a.RLock()
	defer a.RUnlock()
	return fmt.Sprintf("%+v", a.List)
}
