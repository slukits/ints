// Copyright 2022 Stephan.Lukits@gmail.com. All rights reserved.
// Use of this source code is governed by the MIT license.

package ints_test

import (
	"fmt"

	"github.com/slukits/ints"
)

func Example() {
	d1, err := ints.UDec.From.Str("19.5") // string to UDecimal
	if err != nil {
		panic(err)
	}

	d2, err := ints.UDec.From.Float(22.5) // float to UDecimal
	if err != nil {
		panic(err)
	}

	fmt.Println("result:", ints.UDec.MAdd(d1, d2).Str(ints.UDec))
	// Output: result: 42,00
}

func ExampleUFloats() {
	fmt.Println("result:",
		ints.UDec.Float.MAdd(13.4384, 28.5616).Str(ints.UDec))
	// Output: result: 42,00
}

func Example_floatDecimals() {
	fmt.Println("result:",
		ints.UDec.Float.MAdd(13.4384, 28.5616).Str(ints.UDec))
	// Output: result: 42,00
}
