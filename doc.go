// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package ints provides a fast [ints.Set]-type for non-negative small
// integers and a light weight [ints.UDecimal] type with its associated
// [ints.UContext] for non negative integer based decimal arithmetics.
// Skimming through [ints.Set]'s API should suffice to learn what you
// can do with provided Set-type.
//
// To use the non-negative integer based decimal arithmetics a
// [ints.UContext] instance is needed to create [ints.UDecimal] values
// and to operate on them.  [ints.UDec] provides the ready to use
// default context.
//
//	d1, err := ints.UDec.From.Str("19.5") // string to UDecimal
//	if err != nil {
//	    panic(err)
//	}
//
//	d2, err := ints.UDec.From.Float(22.5) // float to UDecimal
//	if err != nil {
//	    panic(err)
//	}
//
//	fmt.Println("result:", ints.UDec.MAdd(d1, d2).Str(ints.UDec))
//	// prints "result: 42,00"
//
// Note that [ints.UDecimal] values are not aware of their number of
// fractionals, i.e. positions after the decimal separator.  This
// information is kept in the context.  Hence to get a string
// representation of a [ints.UDecimal] value a context must be provided.
//
// There are also 'Must'-versions if you feel save:
//
//	result := ints.UDec.MAdd(
//	    ints.UDec.From.MStr("2.5"),
//	    ints.UDec.From.MFloat(2.5),
//	)
//
// and the possibility to calculate directly with floats in a decimal
// manner.
//
//	sum := ints.UDec.Float.MAdd(2.5, 2.5)
//
// I.e. provided floats are converted to [ints.UDecimal]-values before
// the operation is performed and the operation returns an default
// context's [ints.UDecimal].
//
// A Context has two sets of flags the arithmetic flags and format flags.
// The former control the [ints.UDecimal] conversion and their
// arithmetics.  The later control the string representation of a
// [ints.UDecimal] value.  Once a Context instance is created its
// arithmetic flags are immutable.  The [ints.UDec] arithmetic flags
// default to [DOT_SEPARATOR] | [SIX_FRACTIONALS].  I.e. the string to
// [UDecimal] conversion expects a string with a dot decimal separator
// and the last six positions of a converted [UDecimal] represent its
// fractionals.  To create a context with different arithmetic flags the
// [ints.UDec]'s [UContext.New] method must be used
//
//	dec := ints.UDec.New(COMMA_SEPARATOR|FOUR_FRACTIONALS, DEFAULTS)
//	d := dec.From.MStr("2,5")
//
// The second argument sets the format flags which in above case are
// simply a copy of [ints.UDec]'s format flags.
//
// NOTE while I try to write idiomatic go code I value a convenient to
// use (practical) API higher than the "make the zero type usable"
// idiom.  I couldn't find a way to make the zero [UContext] usable
// while providing the API it has.  Hence you should never instantiate a
// [UContext], [UConvert], [UDecimal] or [UFloats] type directly.  It is
// all done for you by a single ints.UDec.New call.  The mentioned types
// are defined as public anyway to provide their methods' documentations
// in the go documentation server.
//
// Format flags may be changed at any time
//
//	dec.SetFmt(DOT_SEPARATOR)
//	fmt.Print(d.Str(dec)) // prints "2.50"
//
// While [UDecimal.Str] truncates superfluous fractionals [UDecimal.Rnd]
// returns a "round to even" string representation of a [ints.UDecimal]
// value.
package ints
