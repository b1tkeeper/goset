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

type Set interface {
	// Add adds an element to the set. Returns whether
	// the item was added.
	Add(val interface{}) bool

	// Cardinality Returns the number of elements in the set.
	Cardinality() int

	// Size Returns the number of elements in the set.
	Size() int

	// Clear removes all elements from the set, leaving
	// the empty set.
	Clear()

	// Clone returns a deep-clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Set

	// Contains returns whether the given items
	// are all in the set.
	Contains(val ...interface{}) bool

	// Difference returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other Set) Set

	// Equal determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other Set) bool

	// Intersect returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other Set) Set

	// IsProperSubset determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsProperSubset
	// will panic.
	IsProperSubset(other Set) bool

	// IsProperSuperset determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsProperSuperset(other Set) bool

	// IsSubset determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will
	// panic.
	IsSubset(other Set) bool

	// IsSuperset determines if every element in the other set
	// is in this set.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsSuperset(other Set) bool

	// Each iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(elem interface{}) bool)

	// Iter returns a channel of elements that you can
	// range over.
	Iter() <-chan interface{}

	// Iterator returns an Iterator object that you can
	// use to range over the set.
	Iterator() *Iterator

	// Remove remove a single element from the set.
	Remove(i interface{})

	// String provides a convenient string representation
	// of the current state of the set.
	String() string

	// SymmetricDifference returns a new set with all elements which are
	// in either this set or the other set but not in both.
	//
	// Note that the argument to SymmetricDifference
	// must be of the same type as the receiver
	// of the method. Otherwise, SymmetricDifference
	// will panic.
	SymmetricDifference(other Set) Set

	// Union returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the
	// same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	Union(other Set) Set

	// Pop removes and returns an arbitrary item from the set.
	Pop() (interface{}, bool)

	// ToSlice returns the members of the set as a slice.
	ToSlice() []interface{}

	// MarshalJSON will marshal the set into a JSON-based representation.
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON will unmarshal a JSON-based byte slice into a full Set datastructure.
	// For this to work, set subtypes must implemented the Marshal/Unmarshal interface.
	UnmarshalJSON(b []byte) error
}

// NewSet creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewSet(vals ...interface{}) Set {
	s := newThreadSafeSet()
	for _, item := range vals {
		s.Add(item)
	}
	return &s
}

// NewThreadUnsafeSet creates and returns a new set with the given elements.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeSet(vals ...interface{}) Set {
	s := newThreadUnsafeSet()
	for _, item := range vals {
		s.Add(item)
	}
	return &s
}
