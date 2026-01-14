# Map Utils

`map_utils` is a Go library providing a collection of utility functions for working with maps. It leverages Go generics and Go 1.23+ iterators (`iter` package) to offer a fluent and type-safe experience for common map operations like filtering, mapping, reducing, and converting.

## Installation

```bash
go get github.com/zauberhaus/map_utils
```

## Features

*   **Filtering & Selection**: `Select`, `Delete`, `CountFunc`, `ExistsFunc`
*   **Existence Checks**: `ContainsKey`, `Contains`
*   **Transformation**: `Remap`, `Convert`
*   **Aggregation**: `Summarize`
*   **Access**: `First`, `Last`, `At` (access by index based on sorted keys)
*   **Conversion**: `Slice` (to slice), `Join` (to string)
*   **Iterators**: `RemapFuncSeq`, `WeightFuncSeq`, `SliceFuncSeq`

## Usage

### Filtering

```go
import "github.com/zauberhaus/map_utils"

m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
// Select even values
result := map_utils.Select(m, func(key int, val int) bool {
    return val%2 == 0
})
// result: {2: 2, 4: 4}
```

### Remapping

Transform keys and values.

```go
m := map[int]int{1: 10, 2: 20}
result := map_utils.Remap(m, func(k, v int) (string, string, error) {
    return fmt.Sprintf("k%d", k), fmt.Sprintf("v%d", v), nil
})
// result: {"k1": "v10", "k2": "v20"}
```

### Converting to Slice

Convert a map to a slice, optionally filtering or transforming elements.

```go
m := map[string]int{"a": 1, "b": 2}
s := map_utils.Slice(m, func(k string, v int) (*int, error) {
    return &v, nil
})

for _, v := range s {
    fmt.Println(v)
}
```

### Ordered Access

Access map elements by index (keys are sorted implicitly).

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
val, err := map_utils.At(m, 1) // Access element at index 1 (sorted by key)
// val: 2
```

## License

Copyright 2026 Zauberhaus

Licensed to Zauberhaus under one or more agreements.
Zauberhaus licenses this file to you under the Apache 2.0 License.
See the LICENSE file in the project root for more information.