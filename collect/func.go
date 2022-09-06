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

func equals[E comparable](a, b Collection[E]) bool {
	if a == b {
		return true
	}
	if a.Size() != b.Size() {
		return false
	}
	itr1 := a.Iterator()
	itr2 := b.Iterator()
	for itr1.HasNext() && itr2.HasNext() {
		e1, err := itr1.Next()
		if err != nil {
			return false
		}
		e2, err := itr2.Next()
		if err != nil {
			return false
		}
		if e1 != e2 {
			return false
		}
	}
	return !(itr1.HasNext() || itr2.HasNext())
}
