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
	"reflect"
)

// List 有序集合接口
type List[E any] interface {
	Collection[E]

	// ReplaceAll 将该列表的每个元素替换为 operator 运算符应用于该元素的结果
	ReplaceAll(operator UnaryOperator[E])

	// Sort 对集合元素进行排序。
	Sort(less SortLess[E])

	// Get 返回此列表中指定位置的元素。索引不在有效范围会返回越界错误
	Get(index int) (E, error)

	// Set 用指定的元素替换此列表中指定位置的元素。索引不在有效范围会返回越界错误
	// 返回旧值
	Set(index int, e E) (E, error)

	// AddAt 将指定的元素插入此列表中的指定位置。 将当前位于该位置的元素（如果有）及后面的元素元素向后移动。
	AddAt(index int, e E) error

	// RemoveAt 删除该列表中指定位置的元素。 将后续元素向左移动。 返回被删除的元素。索引不在有效范围会返回越界错误
	RemoveAt(index int) (E, error)

	// IndexOf 返回此列表中指定元素的第一次出现的索引，如果此列表不包含元素，则返回-1。
	IndexOf(e E) int

	// LastIndexOf 返回此列表中指定元素的最后一次出现的索引，如果此列表不包含元素，则返回-1。
	LastIndexOf(e E) int

	// ListIterator 返回列表迭代器
	ListIterator() ListIterator[E]

	// ListIteratorAt 返回列表迭代器,指定开始迭代位置
	ListIteratorAt(index int) ListIterator[E]

	// SubList 返回列表中指定的fromIndex （含）和toIndex之间的部分
	SubList(fromIndex, toIndex int) List[E]

	// RemoveN 从集合中移除指定的元素
	// 参数 n 表示移除个数, -1表示全部移除
	// 成功移除返回true, 不存在指定元素返回false
	RemoveN(e E, n int) int

	// RemoveIfN 将集合中的元素应用filter函数，如果函数返回true元素将被移除
	// 参数 n 表示移除个数, -1表示全部移除
	// 返回移除的元素个数
	RemoveIfN(filter Predicate[E], n int) int
}

// List 底层实现结构选择
const (
	_ = iota
	DataStructSlice
	DataStructLinked
)

type ListConfig struct {
	// InitialCapacity 初始容量，底层结构为DataStructSlice时有效，默认16
	InitialCapacity int
	// Safe 是否需要并发安全，需要在多个goroutine中并发操作时需要设置为true，默认false
	Safe bool
	// DataStruct 底层实现结构，默认为切片实现
	DataStruct int
	// EqualComparator 元素相等比较函数
	// 如何元素是可比较的如：int、string等，可以不设置，默认使用 == 比较
	// 如果元素不可比较，需要设置改字段，如果没有设置当调用需要比较元素的方法时，比较结果永远返回false
	EqualComparator constraints.EqualComparator[any]
}

// DefaultListConfig 默认配置
var DefaultListConfig = ListConfig{
	InitialCapacity: 16,
	Safe:            false,
	DataStruct:      DataStructSlice,
}

type ListOption func(config *ListConfig)

// NewList 根据配置创建一个 List
func NewList[E any](config ListConfig, options ...ListOption) List[E] {
	for _, option := range options {
		option(&config)
	}
	if config.InitialCapacity < 1 {
		config.InitialCapacity = 16
	}
	if config.EqualComparator == nil {
		config.EqualComparator = DefaultEqualFunc()
	}
	var list List[E]
	if config.DataStruct == DataStructLinked {
		list = NewLinkedList[E](AnyEqualComparableFunc[E](func(v1, v2 E) bool {
			return config.EqualComparator.Equal(v1, v2)
		}))
	} else {
		list = NewArrayList[E](config.InitialCapacity, AnyEqualComparableFunc[E](func(v1, v2 E) bool {
			return config.EqualComparator.Equal(v1, v2)
		}))
	}
	if config.Safe {
		list = NewSafeList[E](list)
	}
	return list
}

func WithEqualFunc[E any](f func(v1, v2 E) bool) ListOption {
	return func(config *ListConfig) {
		config.EqualComparator = AnyEqualComparableFunc[any](func(v1, v2 any) bool {
			return f(v1.(E), v2.(E))
		})
	}
}

func DefaultEqualFunc() constraints.EqualComparator[any] {
	return AnyEqualComparableFunc[any](func(v1, v2 any) bool {
		switch v1.(type) {
		case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint, float32, float64, string:
			return v1 == v2
		}
		val := reflect.ValueOf(v1)
		if val.Type().Comparable() {
			return v1 == v2
		}
		return reflect.DeepEqual(v1, v2)
	})
}
