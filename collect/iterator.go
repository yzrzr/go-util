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

// Iterator 集合迭代器接口
type Iterator[E any] interface {
	// HasNext 如果有更多的元素，则返回true
	HasNext() bool

	// Next 返回迭代中的下一个元素
	Next() (E, error)

	// Remove 从底层集合中删除此迭代器返回的最后一个元素。 此方法只能在调用一次 Next() 后用一次
	Remove() error

	// ForEachRemaining 对每个剩余元素执行给定的操作，直到所有元素都被处理或返回错误
	ForEachRemaining(action Consumer[E]) error

	// Close 关闭迭代器
	Close()
}
