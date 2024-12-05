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
	"github.com/yzrzr/go-util/constraints"
)

func equals[E any](a, b Collection[E]) bool {
	if a == b {
		return true
	}
	if a.Size() != b.Size() {
		return false
	}
	itr1 := a.Iterator()
	itr2 := b.Iterator()
	comparator := a.GetEqualComparator()
	for itr1.HasNext() && itr2.HasNext() {
		e1, err1 := itr1.Next()
		e2, err2 := itr2.Next()
		if err1 != nil || err2 != nil || !comparator.Equal(e1, e2) {
			return false
		}
	}
	return !(itr1.HasNext() || itr2.HasNext())
}

// SortLessOrdered 基础类型升序排序方法
// 参数 asc 表示是否为升序
// List.Sort(SortLessOrdered(true))
func SortLessOrdered[E constraints.Ordered](aes bool) SortLess[E] {
	if aes {
		return func(e1, e2 E) bool {
			return e1 < e2
		}
	}
	return func(e1, e2 E) bool {
		return e1 > e2
	}
}

// SortLessComparable 基础类型降序排序方法
// 参数 asc 表示是否为升序
// List.Sort(SortLessComparable(true))
func SortLessComparable[E constraints.Comparable[E]](aes bool) SortLess[E] {
	if aes {
		return func(e1, e2 E) bool {
			return e1.Compare(e2) == -1
		}
	}
	return func(e1, e2 E) bool {
		return e1.Compare(e2) == 1
	}
}
