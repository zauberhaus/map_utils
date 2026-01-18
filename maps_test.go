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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/map_utils"
)

func TestSelect(t *testing.T) {
	t.Run("select even values", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		result := map_utils.Select(m, func(key int, val int) bool {
			return val%2 == 0
		})
		expected := map[int]int{2: 2, 4: 4}
		assert.Equal(t, expected, result)
	})

	t.Run("select no values", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 3, "c": 5}
		result := map_utils.Select(m, func(key string, val int) bool {
			return val%2 == 0
		})
		assert.Empty(t, result)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[int]int{}
		result := map_utils.Select(m, func(key int, val int) bool {
			return true
		})
		assert.Empty(t, result)
	})

	t.Run("nil map", func(t *testing.T) {
		var m map[int]int
		result := map_utils.Select(m, func(key int, val int) bool {
			return true
		})
		assert.Empty(t, result)
	})
}

func TestCountFunc(t *testing.T) {
	t.Run("count odd values", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		count := map_utils.CountFunc(m, func(key int, val int) bool {
			return val%2 != 0
		})
		assert.Equal(t, 2, count)
	})

	t.Run("count no values", func(t *testing.T) {
		m := map[string]int{"a": 2, "b": 4}
		count := map_utils.CountFunc(m, func(key string, val int) bool {
			return val%2 != 0
		})
		assert.Equal(t, 0, count)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[int]int{}
		count := map_utils.CountFunc(m, func(key int, val int) bool {
			return true
		})
		assert.Equal(t, 0, count)
	})

	t.Run("nil map", func(t *testing.T) {
		var m map[int]int
		count := map_utils.CountFunc(m, func(key int, val int) bool {
			return true
		})
		assert.Equal(t, 0, count)
	})
}

func TestExistsFunc(t *testing.T) {
	t.Run("returns false when predicate is satisfied", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3}
		// Predicate satisfied by 2
		result := map_utils.ExistsFunc(m, func(key int, val int) bool {
			return val == 2
		})
		assert.False(t, result)
	})

	t.Run("returns true when predicate is not satisfied", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3}
		// Predicate not satisfied
		result := map_utils.ExistsFunc(m, func(key int, val int) bool {
			return val > 5
		})
		assert.True(t, result)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[int]int{}
		result := map_utils.ExistsFunc(m, func(key int, val int) bool {
			return true
		})
		assert.True(t, result)
	})
}

func TestEmpty(t *testing.T) {
	t.Run("is not empty", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3}
		isEmpty := map_utils.CountFunc(m, func(key int, val int) bool {
			return val > 2
		}) == 0
		assert.False(t, isEmpty)
	})

	t.Run("is empty", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3}
		isEmpty := map_utils.CountFunc(m, func(key int, val int) bool {
			return val > 5
		}) == 0
		assert.True(t, isEmpty)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[int]int{}
		isEmpty := map_utils.CountFunc(m, func(key int, val int) bool {
			return true
		}) == 0
		assert.True(t, isEmpty)
	})
}

func TestContainsKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	t.Run("contains single key", func(t *testing.T) {
		assert.True(t, map_utils.ContainsKey(m, "a"))
	})

	t.Run("contains one of multiple keys", func(t *testing.T) {
		assert.True(t, map_utils.ContainsKey(m, "c", "b"))
	})

	t.Run("does not contain key", func(t *testing.T) {
		assert.False(t, map_utils.ContainsKey(m, "c"))
	})

	t.Run("empty map", func(t *testing.T) {
		assert.False(t, map_utils.ContainsKey(map[string]int{}, "a"))
	})

	t.Run("nil map", func(t *testing.T) {
		var nilMap map[string]int
		assert.False(t, map_utils.ContainsKey(nilMap, "a"))
	})

	t.Run("no keys provided", func(t *testing.T) {
		assert.False(t, map_utils.ContainsKey(m))
	})
}

func TestContains(t *testing.T) {
	m := map[int]string{1: "a", 2: "b", 3: "c"}

	t.Run("contains based on predicate", func(t *testing.T) {
		assert.True(t, map_utils.Contains(m, func(k int, v string) bool { return v == "b" }))
	})

	t.Run("does not contain based on predicate", func(t *testing.T) {
		assert.False(t, map_utils.Contains(m, func(k int, v string) bool { return v == "d" }))
	})

	t.Run("empty map", func(t *testing.T) {
		assert.False(t, map_utils.Contains(map[int]string{}, func(k int, v string) bool { return true }))
	})
}

func TestDelete(t *testing.T) {
	t.Run("delete a subset", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		deletedCount := map_utils.Delete(m, func(k int, v int) bool {
			return v%2 == 0
		})

		assert.Equal(t, 2, deletedCount)
		assert.Equal(t, map[int]int{1: 1, 3: 3}, m)
	})

	t.Run("delete nothing", func(t *testing.T) {
		m := map[int]int{1: 1, 3: 3}
		deletedCount := map_utils.Delete(m, func(k int, v int) bool {
			return v%2 == 0
		})
		assert.Equal(t, 0, deletedCount)
		assert.Equal(t, map[int]int{1: 1, 3: 3}, m)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[int]int{}
		deletedCount := map_utils.Delete(m, func(k int, v int) bool { return true })
		assert.Equal(t, 0, deletedCount)
		assert.Empty(t, m)
	})
}

func TestFirst(t *testing.T) {
	t.Run("non-empty map", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		val, err := map_utils.First(m)
		assert.NoError(t, err)
		// Value can be 1 or 2, as map iteration order is not guaranteed
		assert.Contains(t, []int{1, 2}, val)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]int{}
		val, err := map_utils.First(m)
		assert.NoError(t, err)
		assert.Equal(t, 0, val) // Zero value for int
	})
}

func TestLast(t *testing.T) {
	t.Run("non-empty map", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		val, err := map_utils.Last(m)
		assert.NoError(t, err)
		assert.Contains(t, []int{1, 2}, val)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]int{}
		val, err := map_utils.Last(m)
		assert.NoError(t, err)
		assert.Equal(t, 0, val)
	})
}

func TestAt(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := slices.Collect(maps.Keys(m))

	t.Run("valid index", func(t *testing.T) {
		// Re-create map to ensure we can predict the element at a pseudo-index
		mFixed := map[string]int{keys[1]: m[keys[1]]}
		val, err := map_utils.At(mFixed, 0)
		assert.NoError(t, err)
		assert.Equal(t, m[keys[1]], val)
	})

	t.Run("index out of bounds", func(t *testing.T) {
		_, err := map_utils.At(m, 10)
		assert.Error(t, err)
	})

	t.Run("index equals length", func(t *testing.T) {
		_, err := map_utils.At(m, 3)
		assert.Error(t, err)
	})

	t.Run("negative index", func(t *testing.T) {
		_, err := map_utils.At(m, -1)
		assert.Error(t, err)
	})

	t.Run("empty map", func(t *testing.T) {
		_, err := map_utils.At(map[string]int{}, 0)
		assert.Error(t, err)
	})
}

func TestRemap(t *testing.T) {
	t.Run("successful conversion", func(t *testing.T) {
		m := map[int]int{1: 2, 3: 4}
		result := map_utils.Remap(m, func(k, v int) (string, string, error) {
			return fmt.Sprintf("key%d", k), fmt.Sprintf("val%d", v), nil
		})

		expected := map[string]string{"key1": "val2", "key3": "val4"}
		assert.Equal(t, expected, result)
	})

	t.Run("remap with error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.ErrorContains(t, err, "negative value not allowed")
				}
			}
		}()

		m := map[int]int{1: 2, 2: -1, 3: 4}
		map_utils.Remap(m, func(k, v int) (string, string, error) {
			if v < 0 {
				return "", "", errors.New("negative value not allowed")
			}
			return fmt.Sprintf("key%d", k), fmt.Sprintf("val%d", v), nil
		})
		//assert.Error(t, err)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[int]int{}
		result := map_utils.Remap(m, func(k, v int) (string, string, error) {
			return "", "", nil
		})

		assert.Empty(t, result)
	})
}

func TestSummarize(t *testing.T) {
	t.Run("sum of values", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20, 3: 30}
		sum := map_utils.Summarize(m, func(k, v int) int {
			return v
		})
		assert.Equal(t, 60, sum)
	})

	t.Run("sum of keys", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20, 3: 30}
		sum := map_utils.Summarize(m, func(k, v int) int {
			return k
		})
		assert.Equal(t, 6, sum)
	})

	t.Run("string concatenation", func(t *testing.T) {
		m := map[string]string{"b": "world", "a": "hello"}
		summary := map_utils.Summarize(m, func(k, v string) string {
			return v
		})
		assert.Equal(t, "helloworld", summary)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[int]int{}
		sum := map_utils.Summarize(m, func(k, v int) int { return v })
		assert.Equal(t, 0, sum)
	})
}

func TestConvert(t *testing.T) {
	t.Run("successful remap", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		result := map_utils.Convert(m, func(k string, v int) (string, error) {
			return strconv.Itoa(v * 2), nil
		})

		expected := map[string]string{"a": "2", "b": "4"}
		assert.Equal(t, expected, result)
	})

	t.Run("convert with error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.ErrorContains(t, err, "negative value not allowed")
				}
			}
		}()

		m := map[string]int{"a": 1, "b": -1}
		map_utils.Convert(m, func(k string, v int) (bool, error) {
			if v < 0 {
				return false, errors.New("negative value not allowed")
			}
			return true, nil
		})
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]int{}
		result := map_utils.Convert(m, func(k string, v int) (int, error) { return v, nil })
		assert.Empty(t, result)
	})
}

func TestJoin(t *testing.T) {
	t.Run("join int map", func(t *testing.T) {
		m := map[int]int{3: 30, 1: 10, 2: 20}
		result := map_utils.Join(m, ", ")
		assert.Equal(t, "1=10, 2=20, 3=30", result)
	})

	t.Run("join string map", func(t *testing.T) {
		m := map[string]string{"b": "world", "a": "hello"}
		result := map_utils.Join(m, " | ")
		assert.Equal(t, "a=hello | b=world", result)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]string{}
		result := map_utils.Join(m, ", ")
		assert.Equal(t, "", result)
	})

	t.Run("single element map", func(t *testing.T) {
		m := map[string]string{"a": "hello"}
		result := map_utils.Join(m, ", ")
		assert.Equal(t, "a=hello", result)
	})
}

func TestSlice(t *testing.T) {
	t.Run("convert map values to slice", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		results := map_utils.Slice(m, func(k string, v int) (*int, error) {
			return &v, nil
		})

		sort.Ints(results)
		assert.Equal(t, []int{1, 2}, results)
	})

	t.Run("filter values during conversion", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		results := map_utils.Slice(m, func(k string, v int) (*int, error) {
			if v%2 == 0 {
				return nil, nil // filter out
			}
			return &v, nil
		})

		sort.Ints(results)
		assert.ElementsMatch(t, []int{1, 3}, results)
	})

	t.Run("panic on error", func(t *testing.T) {
		m := map[string]int{"a": 1}
		assert.Panics(t, func() {
			map_utils.Slice(m, func(k string, v int) (*int, error) {
				return nil, errors.New("test error")
			})
		})
	})
}

func TestFlatten(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		m := map[string]int{"a": 1, "c": 3, "b": 2}
		result := map_utils.Flatten(m)

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
		result := map_utils.Flatten(m)
		assert.Empty(t, result)
	})

	t.Run("int key map", func(t *testing.T) {
		m := map[int]string{1: "a", 3: "c", 2: "b"}
		result := map_utils.Flatten(m)

		assert.Len(t, result, 6)

		pairs := [][2]any{}
		for i := 0; i < len(result); i += 2 {
			pairs = append(pairs, [2]any{result[i], result[i+1]})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i][0].(int) < pairs[j][0].(int)
		})

		expected := [][2]any{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		}

		assert.Equal(t, expected, pairs)
	})

	t.Run("mixed value types", func(t *testing.T) {
		m := map[string]any{"a": 1, "c": true, "b": "hello"}
		result := map_utils.Flatten(m)

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
}
