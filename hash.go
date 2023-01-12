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
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Hashable interface {
	Hash() string
}

func isHashableObj(obj interface{}) bool {
	switch obj.(type) {
	case Hashable, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128, string:
		return true
	default:
		return false
	}
}

func isNativeHashableObj(obj any) bool {
	if !isHashableObj(obj) {
		return false
	}
	switch obj.(type) {
	case Hashable:
		return false
	default:
		return true
	}
}
func calcHash(obj any) (string, error) {
	if !isHashableObj(obj) {
		return "", errors.New("obj is not a hashable object, can't calculate its hash")
	}
	if !isNativeHashableObj(obj) {
		return obj.(Hashable).Hash(), nil
	}
	switch o := obj.(type) {
	case string:
		return o, nil
	case int:
		return strconv.Itoa(o), nil
	case int8:
		return strconv.FormatInt(int64(o), 10), nil
	case int16:
		return strconv.FormatInt(int64(o), 10), nil
	case int32:
		return strconv.FormatInt(int64(o), 10), nil
	case int64:
		return strconv.FormatInt(o, 10), nil
	case uint:
		return strconv.FormatUint(uint64(o), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(o), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(o), 10), nil
	case uint64:
		return strconv.FormatUint(o, 10), nil
	case uintptr:
		return fmt.Sprintf("%v", o), nil
	case float32:
		return strconv.FormatFloat(float64(o), 'f', 8, 32), nil
	case float64:
		return strconv.FormatFloat(o, 'f', 8, 64), nil
	case complex64:
		return strconv.FormatComplex(complex128(o), 'b', 8, 64), nil
	case complex128:
		return strconv.FormatComplex(o, 'b', 8, 128), nil
	default:
		return "", fmt.Errorf("%s is not a hashable native object, but the ret of isNativeHashableObj seems be true", reflect.TypeOf(obj).String())
	}
}
