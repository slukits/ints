// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

// Floats provides a shortcut for decimal operations on floats by
// automatically converting provided floats to decimal values.  This
// type's intent was to free the tests of unnecessary clutter.
type Floats struct {
	cntx *Context
}

func (flt *Floats) convert(a, b float64) (Decimal, Decimal, error) {
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

// Add converts a and b to [Decimal]-values and returns their sum.  It fails
// if the conversion of a or b fails, or if their addition fails.
func (flt *Floats) Add(a, b float64) (Decimal, error) {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		return 0, err
	}
	return flt.cntx.Add(aDec, bDec)
}

// MAdd is the 'Must'-variant of [Float.Add] which panics if
// corresponding Add-call fails.
func (flt *Floats) MAdd(a, b float64) Decimal {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		panic(err)
	}
	return flt.cntx.MAdd(aDec, bDec)
}

// Sub converts a and b to [Decimal]-values and returns their
// difference.  It fails if the conversion of a or b fails, or if their
// subtraction fails.
func (flt *Floats) Sub(a, b float64) (Decimal, error) {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		return 0, err
	}
	return flt.cntx.Sub(aDec, bDec)
}

// MSub is the 'Must'-variant of [Floats.Sub] which panics if
// corresponding Sub-call fails.
func (flt *Floats) MSub(a, b float64) Decimal {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		panic(err)
	}
	return flt.cntx.MSub(aDec, bDec)
}

// Mult converts a and b to [Decimal]-values and returns their product.
// It fails if the conversion of a or b fails, or if their
// multiplication fails.
func (flt *Floats) Mult(a, b float64) (Decimal, error) {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		return 0, err
	}
	return flt.cntx.Mult(aDec, bDec)
}

// MMult is the 'Must'-variant of [Floats.Mult] which panics if
// corresponding Mult-call fails.
func (flt *Floats) MMult(a, b float64) Decimal {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		panic(err)
	}
	return flt.cntx.MMult(aDec, bDec)
}

// Div converts a and b to [Decimal]-values and returns their quotient.  It
// fails if the conversion of a or b fails, or if their division fails.
func (flt *Floats) Div(a, b float64) (Decimal, error) {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		return 0, err
	}
	return flt.cntx.Div(aDec, bDec)
}

// MDiv is the 'Must'-variant of [Floats.Div] which panics if
// corresponding Div-call fails.
func (flt *Floats) MDiv(a, b float64) Decimal {
	aDec, bDec, err := flt.convert(a, b)
	if err != nil {
		panic(err)
	}
	return flt.cntx.MDiv(aDec, bDec)
}
