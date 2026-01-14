// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package map_utils

import (
	"cmp"
	"iter"
)

func RemapFuncSeq[K1 comparable, V1 any, K2 comparable, V2 any](m iter.Seq2[K1, V1], f func(key K1, val V1) (K2, V2, error)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range m {
			k2, v2, err := f(k, v)
			if err != nil {
				panic(err)
			}

			if !yield(k2, v2) {
				return
			}
		}
	}
}

func WeightFuncSeq[K comparable, V any, S cmp.Ordered](m iter.Seq2[K, V], f func(key K, val V) S) iter.Seq[S] {
	return func(yield func(S) bool) {
		for k, v := range m {
			s := f(k, v)

			if !yield(s) {
				return
			}
		}
	}
}

func SliceFuncSeq[K comparable, V any, R any](m iter.Seq2[K, V], f func(key K, val V) (*R, error)) iter.Seq[R] {
	return func(yield func(R) bool) {
		var nilPtr *R

		for k, v := range m {
			val, err := f(k, v)
			if err != nil {
				panic(err)
			}

			if val != nilPtr {
				if !yield(*val) {
					return
				}
			}
		}
	}
}
