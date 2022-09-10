// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

// UFloats provides a shortcut for decimal operations on floats by
// automatically converting provided floats to decimal values.  This
// type's intent was to free the tests of unnecessary clutter.
type UFloats struct {
	cntx *UContext
}

func (flt *UFloats) convert(a, b float64) (UDecimal, UDecimal, error) {
	aDec, err := flt.cntx.From.Float(a)
	if err != nil {
		return 0, 0, err
	}
	bDec, err := flt.cntx.From.Float(b)
	if err != nil {
		return 0, 0, err
	}
	return aDec, bDec, nil
}

// Add converts a and b to [UDecimal]-values and returns their sum.  It fails
// if the conversion of a or b fails, or if their addition fails.
func (flt *UFloats) Add(a, b float64) (UDecimal, error) {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		return 0, err
	}
	return flt.cntx.Add(aDec, bDec)
}

// MAdd is the 'Must'-variant of [UFloat.Add] which panics if
// corresponding Add-call fails.
func (flt *UFloats) MAdd(a, b float64) UDecimal {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		panic(err)
	}
	return flt.cntx.MAdd(aDec, bDec)
}

// Sub converts a and b to [UDecimal]-values and returns their
// difference.  It fails if the conversion of a or b fails, or if their
// subtraction fails.
func (flt *UFloats) Sub(a, b float64) (UDecimal, error) {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		return 0, err
	}
	return flt.cntx.Sub(aDec, bDec)
}

// MSub is the 'Must'-variant of [UFloats.Sub] which panics if
// corresponding Sub-call fails.
func (flt *UFloats) MSub(a, b float64) UDecimal {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		panic(err)
	}
	return flt.cntx.MSub(aDec, bDec)
}

// Mult converts a and b to [UDecimal]-values and returns their product.
// It fails if the conversion of a or b fails, or if their
// multiplication fails.
func (flt *UFloats) Mult(a, b float64) (UDecimal, error) {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		return 0, err
	}
	return flt.cntx.Mult(aDec, bDec)
}

// MMult is the 'Must'-variant of [UFloats.Mult] which panics if
// corresponding Mult-call fails.
func (flt *UFloats) MMult(a, b float64) UDecimal {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		panic(err)
	}
	return flt.cntx.MMult(aDec, bDec)
}

// Div converts a and b to [UDecimal]-values and returns their quotient.  It
// fails if the conversion of a or b fails, or if their division fails.
func (flt *UFloats) Div(a, b float64) (UDecimal, error) {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		return 0, err
	}
	return flt.cntx.Div(aDec, bDec)
}

// MDiv is the 'Must'-variant of [UFloats.Div] which panics if
// corresponding Div-call fails.
func (flt *UFloats) MDiv(a, b float64) UDecimal {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		panic(err)
	}
	return flt.cntx.MDiv(aDec, bDec)
}
