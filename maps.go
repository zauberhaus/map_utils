// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package map_utils

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"maps"

	"github.com/zauberhaus/slice_utils"
)

func Select[K comparable, V any](m map[K]V, f func(key K, val V) bool) map[K]V {
	result := map[K]V{}

	for k, v := range m {
		if f(k, v) {
			result[k] = v
		}
	}

	return result
}

func CountFunc[K comparable, V any](m map[K]V, f func(key K, val V) bool) int {
	cnt := 0

	for k, v := range m {
		if f(k, v) {
			cnt++
		}
	}

	return cnt
}

func ExistsFunc[K comparable, V any](m map[K]V, f func(key K, val V) bool) bool {
	for k, v := range m {
		if f(k, v) {
			return false
		}
	}

	return true
}

func ContainsKey[K comparable, V any](m map[K]V, keys ...K) bool {
	for _, k := range keys {
		if _, ok := m[k]; ok {
			return true
		}
	}
	return false
}

func Contains[K comparable, V any](m map[K]V, f func(key K, val V) bool) bool {
	for k, v := range m {
		if f(k, v) {
			return true
		}
	}
	return false
}

func Delete[K comparable, V any](m map[K]V, f func(k K, v V) bool) int {
	var keys []K

	for k, v := range m {
		if f(k, v) {
			keys = append(keys, k)
		}
	}

	for _, k := range keys {
		delete(m, k)
	}

	return len(keys)
}

func First[K cmp.Ordered, V any](m map[K]V) (V, error) {
	if len(m) == 0 {
		return *new(V), nil
	}

	return At(m, 0)
}

func Last[K cmp.Ordered, V any](m map[K]V) (V, error) {
	if len(m) == 0 {
		return *new(V), nil
	}

	index := len(m) - 1
	return At(m, index)
}

func At[K cmp.Ordered, V any](m map[K]V, index int) (V, error) {
	if len(m) == 0 {
		return *new(V), fmt.Errorf("utils.At: map is empty")
	}

	if index < 0 || index >= len(m) {
		return *new(V), fmt.Errorf("utils.At: index oyt of bounds")
	}

	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)

	return m[keys[index]], nil
}

func Convert[K comparable, V1 any, V2 any](m map[K]V1, f func(key K, val V1) (V2, error)) map[K]V2 {
	tmp := func(key K, val V1) (K, V2, error) {
		val2, err := f(key, val)
		if err != nil {
			return key, *new(V2), err
		}

		return key, val2, nil
	}

	return maps.Collect(RemapFuncSeq(maps.All(m), tmp))
}

func Remap[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, f func(key K1, val V1) (K2, V2, error)) map[K2]V2 {
	return maps.Collect(RemapFuncSeq(maps.All(m), f))
}

func Slice[Map ~map[K]V, K comparable, V any, S any](m Map, f func(key K, val V) (*S, error)) []S {
	return slices.Collect(SliceFuncSeq(maps.All(m), f))
}

func Summarize[K cmp.Ordered, V any, M ~map[K]V, S cmp.Ordered](m M, f func(key K, val V) S) S {
	return slice_utils.SumSeq(WeightFuncSeq(maps.All(m), f))
}

func Join[K cmp.Ordered, V any](m map[K]V, sep string) string {
	entries := []string{}

	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)

	for _, k := range keys {
		v := m[k]
		entries = append(entries, fmt.Sprintf("%v=%v", k, v))
	}

	return strings.Join(entries, sep)
}

func Flatten[K cmp.Ordered, V any](m map[K]V) []any {
	return slices.Collect(FlattenSeq(maps.All(m)))
}
