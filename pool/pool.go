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

package pool

import (
	"sync"
)

func NewPool[T any](f func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() any {
				return f()
			},
		},
	}
}

type Pool[T any] struct {
	sync.Pool
}

func (t *Pool[T]) Get() T {
	return t.Pool.Get()
}

func (t *Pool[T]) Put(x T) {
	t.Pool.Put(x)
}
