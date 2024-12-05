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
	"strings"
)

type hashSet[E comparable] struct {
	data map[E]struct{}
}

func (h *hashSet[E]) Size() int {
	return len(h.data)
}

func (h *hashSet[E]) IsEmpty() bool {
	return h.Size() == 0
}

func (h *hashSet[E]) Contains(e E) bool {
	_, ok := h.data[e]
	return ok
}

func (h *hashSet[E]) Iterator() Iterator[E] {
	return NewSetIterator[E](h)
}

func (h *hashSet[E]) ToArray() []E {
	arr := make([]E, 0, h.Size())
	for v := range h.data {
		arr = append(arr, v)
	}
	return arr
}

func (h *hashSet[E]) Add(e E) bool {
	if _, ok := h.data[e]; ok {
		return false
	}
	h.data[e] = struct{}{}
	return true
}

func (h *hashSet[E]) Remove(e E) bool {
	if _, ok := h.data[e]; ok {
		delete(h.data, e)
		return true
	}
	return false
}

func (h *hashSet[E]) ContainsAll(c Collection[E]) bool {
	itr := c.Iterator()
	for itr.HasNext() {
		if e, err := itr.Next(); err != nil || !h.Contains(e) {
			return false
		}
	}
	return true
}

func (h *hashSet[E]) AddAll(c Collection[E]) {
	_ = c.ForEach(func(e E) error {
		h.data[e] = struct{}{}
		return nil
	})
}

func (h *hashSet[E]) RemoveAll(c Collection[E]) int {
	return h.RemoveIf(func(e E) bool {
		return c.Contains(e)
	})
}

func (h *hashSet[E]) RemoveIf(filter Predicate[E]) int {
	var cnt int
	for k := range h.data {
		if filter(k) {
			delete(h.data, k)
			cnt++
		}
	}
	return cnt
}

func (h *hashSet[E]) RetainAll(c Collection[E]) int {
	return h.RemoveIf(func(e E) bool {
		return !c.Contains(e)
	})
}

func (h *hashSet[E]) Clear() {
	h.data = make(map[E]struct{})
}

func (h *hashSet[E]) Equals(c Collection[E]) bool {
	if h == c {
		return true
	}
	if h.Size() != c.Size() {
		return false
	}
	return h.ContainsAll(c)
}

func (h *hashSet[E]) ForEach(f Consumer[E]) error {
	var err error
	for k := range h.data {
		if err = f(k); err != nil {
			return err
		}
	}
	return nil
}

func (h *hashSet[E]) GetEqualComparator() constraints.EqualComparator[E] {
	return nil
}

func (h *hashSet[E]) String() string {
	build := strings.Builder{}
	build.WriteByte('[')
	l := h.Size()
	i := 0
	for v := range h.data {
		build.WriteString(fmt.Sprintf("%v", v))
		if i < l-1 {
			build.WriteByte(' ')
		}
		i++
	}
	build.WriteByte(']')
	return build.String()
}
