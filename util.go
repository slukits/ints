// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

import "math"

// Max returns the maximum of given ints a, b and variadic ii.
func Max(a int, b int, ii ...int) int {
	ii, max := append([]int{a, b}, ii...), math.MinInt
	for _, i := range ii {
		if i <= max {
			continue
		}
		max = i
	}
	return max
}
