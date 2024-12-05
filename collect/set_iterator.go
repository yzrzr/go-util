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

import "math"

func NewSetIterator[E comparable](set Set[E]) Iterator[E] {
	return &setIterator[E]{
		lastRet: -1,
		size:    set.Size(),
		values:  set.ToArray(),
		set:     set,
	}
}

type setIterator[E comparable] struct {
	cursor, lastRet, size int

	isClose bool
	set     Set[E]
	values  []E
}

func (s *setIterator[E]) HasNext() bool {
	return s.cursor < s.size && s.isClose == false
}

func (s *setIterator[E]) Next() (e E, err error) {
	if s.isClose {
		err = ErrIteratorClose
		return
	}
	i := s.cursor
	if i >= s.size {
		err = ErrNoSuchElement
		return
	}
	s.cursor = i + 1
	s.lastRet = i
	return s.values[i], nil
}

func (s *setIterator[E]) Remove() error {
	if s.isClose {
		return ErrIteratorClose
	}
	if s.lastRet < 0 {
		return ErrIllegalState
	}
	s.set.Remove(s.values[s.lastRet])
	s.lastRet = -1
	return nil
}

func (s *setIterator[E]) ForEachRemaining(action Consumer[E]) error {
	if s.isClose {
		return ErrIteratorClose
	}
	var err error
	for i := s.cursor; i < s.size; i++ {
		err = action(s.values[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *setIterator[E]) Close() {
	s.isClose = true
	s.cursor = math.MaxInt
}
