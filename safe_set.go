// Copyright 2023 Wang Bohan <wangbohan2000@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package goset

import "sync"

type ThreadSafeSet struct {
	sync.RWMutex
	unsafeSet ThreadUnsafeSet
}

func newThreadSafeSet() ThreadSafeSet {
	return ThreadSafeSet{unsafeSet: newThreadUnsafeSet()}
}

// Add adds an element to the set. Returns whether
// the item was added.
func (set *ThreadSafeSet) Add(val interface{}) bool {
	set.Lock()
	ret := set.unsafeSet.Add(val)
	set.Unlock()
	return ret
}

// Cardinality Returns the number of elements in the set.
func (set *ThreadSafeSet) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.unsafeSet.dat)
}

// Size Returns the number of elements in the set.
func (set *ThreadSafeSet) Size() int {
	set.RLock()
	defer set.RUnlock()
	return set.unsafeSet.Size()
}

// Clear removes all elements from the set, leaving
// the empty set.
func (set *ThreadSafeSet) Clear() {
	set.Lock()
	set.unsafeSet = newThreadUnsafeSet()
	set.Unlock()
}

// Clone returns a deep-clone of the set using the same
// implementation, duplicating all keys.
func (set *ThreadSafeSet) Clone() Set {
	set.RLock()
	unsafeClone := set.unsafeSet.Clone().(*ThreadUnsafeSet)
	ret := &ThreadSafeSet{unsafeSet: *unsafeClone}
	set.RUnlock()
	return ret
}

// Contains returns whether the given items
// are all in the set.
func (set *ThreadSafeSet) Contains(val ...interface{}) bool {
	set.RLock()
	ret := set.unsafeSet.Contains(val...)
	set.RUnlock()
	return ret
}

// Difference returns the difference between this set
// and other. The returned set will contain
// all elements of this set that are not also
// elements of other.
//
// Note that the argument to Difference
// must be of the same type as the receiver
// of the method. Otherwise, Difference will
// panic.
func (set *ThreadSafeSet) Difference(other Set) Set {
	o := other.(*ThreadSafeSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.unsafeSet.Difference(&o.unsafeSet).(*ThreadUnsafeSet)
	ret := &ThreadSafeSet{unsafeSet: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

// Equal determines if two sets are equal to each
// other. If they have the same cardinality
// and contain the same elements, they are
// considered equal. The order in which
// the elements were added is irrelevant.
//
// Note that the argument to Equal must be
// of the same type as the receiver of the
// method. Otherwise, Equal will panic.
func (set *ThreadSafeSet) Equal(other Set) bool {
	o := other.(*ThreadSafeSet)

	set.RLock()
	o.RLock()

	ret := set.unsafeSet.Equal(&o.unsafeSet)
	set.RUnlock()
	o.RUnlock()
	return ret
}

// Intersect returns a new set containing only the elements
// that exist only in both sets.
//
// Note that the argument to Intersect
// must be of the same type as the receiver
// of the method. Otherwise, Intersect will
// panic.
func (set *ThreadSafeSet) Intersect(other Set) Set {
	o := other.(*ThreadSafeSet)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.unsafeSet.Intersect(&o.unsafeSet).(*ThreadUnsafeSet)
	ret := &ThreadSafeSet{unsafeSet: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

// IsProperSubset determines if every element in this set is in
// the other set but the two sets are not equal.
//
// Note that the argument to IsProperSubset
// must be of the same type as the receiver
// of the method. Otherwise, IsProperSubset
// will panic.
func (set *ThreadSafeSet) IsProperSubset(other Set) bool {
	o := other.(*ThreadSafeSet)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.unsafeSet.IsProperSubset(&o.unsafeSet)
}

// IsProperSuperset determines if every element in the other set
// is in this set but the two sets are not
// equal.
//
// Note that the argument to IsSuperset
// must be of the same type as the receiver
// of the method. Otherwise, IsSuperset will
// panic.
func (set *ThreadSafeSet) IsProperSuperset(other Set) bool {
	return other.IsProperSubset(set)
}

// IsSubset determines if every element in this set is in
// the other set.
//
// Note that the argument to IsSubset
// must be of the same type as the receiver
// of the method. Otherwise, IsSubset will
// panic.
func (set *ThreadSafeSet) IsSubset(other Set) bool {
	o := other.(*ThreadSafeSet)

	set.RLock()
	o.RLock()
	ret := set.unsafeSet.IsSubset(&o.unsafeSet)
	set.RUnlock()
	o.RUnlock()
	return ret
}

// IsSuperset determines if every element in the other set
// is in this set.
//
// Note that the argument to IsSuperset
// must be of the same type as the receiver
// of the method. Otherwise, IsSuperset will
// panic.
func (set *ThreadSafeSet) IsSuperset(other Set) bool {
	return other.IsSubset(set)
}

// Each iterates over elements and executes the passed func against each element.
// If passed func returns true, stop iteration at the time.
func (set *ThreadSafeSet) Each(cb func(elem interface{}) bool) {
	set.RLock()
	for _, obj := range set.unsafeSet.dat {
		if cb(obj) {
			break
		}
	}
	set.RUnlock()
}

// Iter returns a channel of elements that you can
// range over.
func (set *ThreadSafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		set.RLock()

		for _, obj := range set.unsafeSet.dat {
			ch <- obj
		}
		close(ch)
		set.RUnlock()
	}()

	return ch
}

// Iterator returns an Iterator object that you can
// use to range over the set.
func (set *ThreadSafeSet) Iterator() *Iterator {
	iterator, ch, stopCh := newIterator()

	go func() {
		set.RLock()
	L:
		for _, obj := range set.unsafeSet.dat {
			select {
			case <-stopCh:
				break L
			case ch <- obj:
			}
		}
		close(ch)
		set.RUnlock()
	}()

	return iterator
}

// Remove remove a single element from the set.
func (set *ThreadSafeSet) Remove(i interface{}) {
	hash, err := calcHash(i)
	if err != nil {
		panic(err)
	}
	set.Lock()
	delete(set.unsafeSet.dat, hash)
	set.Unlock()
}

// String provides a convenient string representation
// of the current state of the set.
func (set *ThreadSafeSet) String() string {
	set.RLock()
	ret := set.unsafeSet.String()
	set.RUnlock()
	return ret
}

// SymmetricDifference returns a new set with all elements which are
// in either this set or the other set but not in both.
//
// Note that the argument to SymmetricDifference
// must be of the same type as the receiver
// of the method. Otherwise, SymmetricDifference
// will panic.
func (set *ThreadSafeSet) SymmetricDifference(other Set) Set {
	o := other.(*ThreadSafeSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.unsafeSet.Difference(&o.unsafeSet).(*ThreadUnsafeSet)
	ret := &ThreadSafeSet{unsafeSet: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

// Union returns a new set with all elements in both sets.
//
// Note that the argument to Union must be of the
// same type as the receiver of the method.
// Otherwise, IsSuperset will panic.
func (set *ThreadSafeSet) Union(other Set) Set {
	o := other.(*ThreadSafeSet)
	set.RLock()
	o.RLock()
	unsafeUnion := set.unsafeSet.Union(&o.unsafeSet).(*ThreadUnsafeSet)
	ret := &ThreadSafeSet{unsafeSet: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

// Pop removes and returns an arbitrary item from the set.
func (set *ThreadSafeSet) Pop() (interface{}, bool) {
	set.Lock()
	defer set.Unlock()
	return set.unsafeSet.Pop()
}

// ToSlice returns the members of the set as a slice.
func (set *ThreadSafeSet) ToSlice() []interface{} {
	objs := make([]interface{}, 0, set.Size())
	set.RLock()
	for _, obj := range set.unsafeSet.dat {
		objs = append(objs, obj)
	}
	set.RUnlock()
	return objs
}

// MarshalJSON will marshal the set into a JSON-based representation.
func (set *ThreadSafeSet) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.unsafeSet.MarshalJSON()
	set.RUnlock()

	return b, err
}

// UnmarshalJSON will unmarshal a JSON-based byte slice into a full Set datastructure.
// For this to work, set subtypes must implemented the Marshal/Unmarshal interface.
func (set *ThreadSafeSet) UnmarshalJSON(b []byte) error {
	set.RLock()
	err := set.unsafeSet.UnmarshalJSON(b)
	set.RUnlock()

	return err
}
