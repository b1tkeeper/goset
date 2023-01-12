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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type ThreadUnsafeSet struct {
	dat map[string]any // Store {$hash: $value} of elem
	typ reflect.Type   // Set's data type
}

func newThreadUnsafeSet() ThreadUnsafeSet {
	return ThreadUnsafeSet{dat: map[string]any{}, typ: nil}
}

func (set *ThreadUnsafeSet) Add(val interface{}) bool {
	typ := reflect.ValueOf(val).Type()
	if set.typ != nil && set.typ != typ {
		panic(
			fmt.Errorf(
				"type conflict when you add a new element to set (type of set elem: %s, type of new elem %s)",
				set.typ, typ,
			))
	}
	if set.typ == nil {
		set.typ = typ
	}
	hash, err := calcHash(val)
	if err != nil {
		panic(err)
	}
	set.dat[hash] = val
	return true
}

func (set *ThreadUnsafeSet) Cardinality() int {
	return len(set.dat)
}

func (set *ThreadUnsafeSet) Size() int {
	return set.Cardinality()
}

func (set *ThreadUnsafeSet) Clear() {
	*set = newThreadUnsafeSet()
}

func (set *ThreadUnsafeSet) Clone() Set {
	cloned := newThreadUnsafeSet()
	cloned.typ = set.typ
	cloned.dat = make(map[string]interface{}, set.Size())
	for _, elem := range set.dat {
		cloned.Add(elem)
	}
	return &cloned
}

func (set *ThreadUnsafeSet) Contains(val ...interface{}) bool {
	for _, v := range val {
		hash, err := calcHash(v)
		if err != nil {
			return false
		}
		if _, ok := set.dat[hash]; !ok {
			return false
		}
	}
	return true
}

func (set *ThreadUnsafeSet) Difference(other Set) Set {
	o := other.(*ThreadUnsafeSet)
	diff := newThreadUnsafeSet()
	for _, obj := range set.dat {
		if !o.Contains(obj) {
			diff.Add(obj)
		}
	}
	return &diff
}

func (set *ThreadUnsafeSet) Equal(other Set) bool {
	if set.Size() != other.Size() {
		return true
	}
	o := other.(*ThreadUnsafeSet)
	for _, obj := range set.dat {
		if !o.Contains(obj) {
			return false
		}
	}
	return true
}

func (set *ThreadUnsafeSet) Intersect(other Set) Set {
	o := other.(*ThreadUnsafeSet)
	intersection := newThreadUnsafeSet()

	if set.Size() < o.Size() {
		for _, obj := range set.dat {
			if o.Contains(obj) {
				intersection.Add(obj)
			}
		}
	} else {
		for _, obj := range o.dat {
			if set.Contains(obj) {
				intersection.Add(obj)
			}
		}
	}
	return &intersection
}

func (set *ThreadUnsafeSet) IsProperSubset(other Set) bool {
	return set.Size() < other.Size() && set.IsSubset(other)
}

func (set *ThreadUnsafeSet) IsProperSuperset(other Set) bool {
	return set.Size() > other.Size() && set.IsSuperset(other)
}

func (set *ThreadUnsafeSet) IsSubset(other Set) bool {
	if set.Size() > other.Size() {
		return false
	}
	o := other.(*ThreadUnsafeSet)
	for _, obj := range set.dat {
		if !o.Contains(obj) {
			return false
		}
	}
	return true
}

func (set *ThreadUnsafeSet) IsSuperset(other Set) bool {
	return other.IsSubset(set)
}

func (set *ThreadUnsafeSet) Each(f func(elem interface{}) bool) {
	for _, obj := range set.dat {
		if f(obj) {
			break
		}
	}
}

func (set *ThreadUnsafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for _, obj := range set.dat {
			ch <- obj
		}
		close(ch)
	}()
	return ch
}

func (set *ThreadUnsafeSet) Iterator() *Iterator {
	iterator, ch, stopCh := newIterator()

	go func() {
	L:
		for _, obj := range set.dat {
			select {
			case <-stopCh:
				break L
			case ch <- obj:
			}
		}
		close(ch)
	}()
	return iterator
}

func (set *ThreadUnsafeSet) Remove(i interface{}) {
	hash, err := calcHash(i)
	if err != nil {
		panic(err)
	}
	delete(set.dat, hash)
}

func (set *ThreadUnsafeSet) String() string {
	var builder strings.Builder
	builder.WriteString("goset.ThreadUnsafeSet{ ")
	atLeastOnce := false
	for _, obj := range set.dat {
		builder.WriteString(fmt.Sprintf("%v, ", obj))
		atLeastOnce = true
	}
	ret := builder.String()
	if atLeastOnce {
		ret = ret[:len(ret)-2]
	}
	return ret + " }"
}

func (set *ThreadUnsafeSet) SymmetricDifference(other Set) Set {
	o := other.(*ThreadUnsafeSet)
	diff := newThreadUnsafeSet()
	for _, obj := range set.dat {
		if !o.Contains(obj) {
			diff.Add(obj)
		}
	}
	for _, obj := range o.dat {
		if !set.Contains(obj) {
			diff.Add(obj)
		}
	}
	return &diff
}

func (set *ThreadUnsafeSet) Union(other Set) Set {
	o := other.(*ThreadUnsafeSet)
	union := newThreadUnsafeSet()
	for _, obj := range set.dat {
		union.Add(obj)
	}
	for _, obj := range o.dat {
		union.Add(obj)
	}
	return &union
}

func (set *ThreadUnsafeSet) Pop() (interface{}, bool) {
	for hash, obj := range set.dat {
		delete(set.dat, hash)
		return obj, true
	}
	return nil, false
}

func (set *ThreadUnsafeSet) ToSlice() []interface{} {
	objs := make([]interface{}, 0, set.Size())
	for _, obj := range set.dat {
		objs = append(objs, obj)
	}
	return objs
}

func (set *ThreadUnsafeSet) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, set.Size())

	for _, obj := range set.dat {
		b, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

func (set *ThreadUnsafeSet) UnmarshalJSON(b []byte) error {
	var i []interface{}

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&i)
	if err != nil {
		return err
	}
	for _, v := range i {
		set.Add(v)
	}
	return nil
}
