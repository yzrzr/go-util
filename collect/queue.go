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

type Queue[E comparable] interface {
	Collection[E]

	// Put 如果队列容量足够，立即将元素插入队列中，返回true。容量不足会插入失败返回false
	Put(e E) bool

	// Take 获取并移除该队列的头部，
	// 如果队列为空，第二个返回值为 false
	Take() (E, bool)

	// Peek 检索但不删除该队列的头部
	// 如果队列为空，第二个返回值为 false
	Peek() (E, bool)
}
