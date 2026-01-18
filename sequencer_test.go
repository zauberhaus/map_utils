// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package map_utils_test

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/map_utils"
)

func TestRemapFuncSeq(t *testing.T) {
	t.Run("successful remap", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20}
		seq := maps.All(m)

		remapped := map_utils.RemapFuncSeq(seq, func(k, v int) (string, string, error) {
			return fmt.Sprintf("k%d", k), fmt.Sprintf("v%d", v), nil
		})

		result := maps.Collect(remapped)
		expected := map[string]string{"k1": "v10", "k2": "v20"}
		assert.Equal(t, expected, result)
	})

	t.Run("panic on error", func(t *testing.T) {
		m := map[int]int{1: 10}
		seq := maps.All(m)

		remapped := map_utils.RemapFuncSeq(seq, func(k, v int) (string, string, error) {
			return "", "", errors.New("remap error")
		})

		assert.PanicsWithError(t, "remap error", func() {
			for range remapped {
			}
		})
	})

	t.Run("empty sequence", func(t *testing.T) {
		m := map[int]int{}
		seq := maps.All(m)

		remapped := map_utils.RemapFuncSeq(seq, func(k, v int) (string, string, error) {
			return "key", "val", nil
		})

		result := maps.Collect(remapped)
		assert.Empty(t, result)
	})

	t.Run("early termination", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20, 3: 30}
		seq := maps.All(m)

		count := 0
		remapped := map_utils.RemapFuncSeq(seq, func(k, v int) (int, int, error) {
			count++
			return k, v, nil
		})

		remapped(func(k, v int) bool {
			return false
		})

		assert.Equal(t, 1, count)
	})
}

func TestWeightFuncSeq(t *testing.T) {
	t.Run("calculate weights", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		seq := maps.All(m)

		weighted := map_utils.WeightFuncSeq(seq, func(k string, v int) int {
			return v * 10
		})

		result := slices.Collect(weighted)
		sort.Ints(result)
		assert.Equal(t, []int{10, 20, 30}, result)
	})

	t.Run("empty sequence", func(t *testing.T) {
		m := map[string]int{}
		seq := maps.All(m)

		weighted := map_utils.WeightFuncSeq(seq, func(k string, v int) int {
			return v
		})

		result := slices.Collect(weighted)
		assert.Empty(t, result)
	})

	t.Run("early termination", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20, 3: 30}
		seq := maps.All(m)

		count := 0
		weighted := map_utils.WeightFuncSeq(seq, func(k, v int) int {
			count++
			return v
		})

		weighted(func(w int) bool {
			return false
		})

		assert.Equal(t, 1, count)
	})
}

func TestSliceFuncSeq(t *testing.T) {
	t.Run("filter and map", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		seq := maps.All(m)

		sliceSeq := map_utils.SliceFuncSeq(seq, func(k string, v int) (*int, error) {
			if v%2 == 0 {
				return nil, nil
			}
			return &v, nil
		})

		result := slices.Collect(sliceSeq)
		sort.Ints(result)
		assert.Equal(t, []int{1, 3}, result)
	})

	t.Run("panic on error", func(t *testing.T) {
		m := map[string]int{"a": 1}
		seq := maps.All(m)

		sliceSeq := map_utils.SliceFuncSeq(seq, func(k string, v int) (*int, error) {
			return nil, errors.New("oops")
		})

		assert.PanicsWithError(t, "oops", func() {
			tmp := slices.Collect(sliceSeq)
			assert.Empty(t, tmp)
		})
	})

	t.Run("early termination", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3}
		seq := maps.All(m)

		count := 0
		sliceSeq := map_utils.SliceFuncSeq(seq, func(k, v int) (*int, error) {
			count++
			return &v, nil
		})

		sliceSeq(func(v int) bool {
			return false
		})

		assert.Equal(t, 1, count)
	})
}

func TestFlattenSeq(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		m := map[string]int{"a": 1, "c": 3, "b": 2}
		seq := maps.All(m)

		weighted := map_utils.FlattenSeq(seq)

		result := slices.Collect(weighted)

		assert.Len(t, result, 6)

		pairs := [][2]any{}
		for i := 0; i < len(result); i += 2 {
			pairs = append(pairs, [2]any{result[i], result[i+1]})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return fmt.Sprintf("%v", pairs[i][0]) < fmt.Sprintf("%v", pairs[j][0])
		})

		expected := [][2]any{
			{"a", 1},
			{"b", 2},
			{"c", 3},
		}

		assert.Equal(t, expected, pairs)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]int{}
		seq := maps.All(m)

		weighted := map_utils.FlattenSeq(seq)

		result := slices.Collect(weighted)
		assert.Empty(t, result)
	})

	t.Run("mixed types", func(t *testing.T) {
		m := map[string]any{"a": 1, "b": "hello", "c": true}
		seq := maps.All(m)

		weighted := map_utils.FlattenSeq(seq)

		result := slices.Collect(weighted)

		assert.Len(t, result, 6)

		pairs := [][2]any{}
		for i := 0; i < len(result); i += 2 {
			pairs = append(pairs, [2]any{result[i], result[i+1]})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return fmt.Sprintf("%v", pairs[i][0]) < fmt.Sprintf("%v", pairs[j][0])
		})

		expected := [][2]any{
			{"a", 1},
			{"b", "hello"},
			{"c", true},
		}

		assert.Equal(t, expected, pairs)
	})

	t.Run("early termination", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20, 3: 30}
		seq := maps.All(m)

		count := 0
		flattened := map_utils.FlattenSeq(seq)

		flattened(func(v any) bool {
			count++
			return count < 3
		})

		assert.Equal(t, 3, count)
	})
}
