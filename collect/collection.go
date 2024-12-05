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

import "github.com/yzrzr/go-util/constraints"

// Collection 集合的根接口
type Collection[E any] interface {
	// Size 返回此集合中的元素数
	Size() int

	// IsEmpty 如果此集合不包含元素，则返回 true
	IsEmpty() bool

	// Contains 如果此集合包含指定的元素，则返回true
	Contains(e E) bool

	// Iterator 返回此集合中的元素的迭代器
	Iterator() Iterator[E]

	// ToArray 一个包含此集合中所有元素的数组
	ToArray() []E

	// Add 将元素加入到集合中。
	// 如果元素是在本次调用中加入集合返回true，如果此集合不允许有重复元素，并且已包含指定的元素，则返回false。
	Add(e E) bool

	// Remove 从集合中移除一个指定的元素，只移除第一个指定的元素
	// 成功移除返回true, 不存在指定元素返回false
	Remove(e E) bool

	// ContainsAll 如果此集合包含指定集合中的所有元素，则返回true，否则返回false
	ContainsAll(c Collection[E]) bool

	// AddAll 将指定集合中的所有元素添加到此集合。
	AddAll(c Collection[E])

	// RemoveAll 删除指定集合中包含的所有此集合的元素。 此调用返回后，此集合将不包含与指定集合相同的元素。
	// 返回删除元素个数
	RemoveAll(c Collection[E]) int

	// RemoveIf 将集合中的元素应用filter函数，如果函数返回true元素将被移除
	// 返回移除的元素个数
	RemoveIf(filter Predicate[E]) int

	// RetainAll 仅保留此集合中包含在指定集合中的元素。
	// 返回移除的元素个数
	RetainAll(c Collection[E]) int

	// Clear 删除集合所有元素。
	Clear()

	// Equals 判断两个集合中元素是否相等
	Equals(c Collection[E]) bool

	// ForEach 迭代集合中的元素，直到所有元素都被处理或返回错误
	ForEach(f Consumer[E]) error

	// GetEqualComparator 返回元素比较器
	GetEqualComparator() constraints.EqualComparator[E]
}

type AnyEqualComparableFunc[E any] func(v1, v2 E) bool

func (f AnyEqualComparableFunc[E]) Equal(v1, v2 E) bool {
	return f(v1, v2)
}
