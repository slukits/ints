// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

import (
	"errors"
	"math"
	"math/big"
	"strconv"
)

// ErrOverflow is returned by any arithmetic operation whose result (or
// intermediate result) would be bigger than Context.Max or negative.
var ErrOverflow = errors.New("ints: dec: overflow")

// ErrDividedByZero is returned by an Div-operation whose second
// argument is zero.
var ErrDividedByZero = errors.New("intdec: div: divided by zero")

// A Context represents an environment to do arithmetic with
// uint64-based decimals.  A Context's zero value is NOT ready to use.
// Create a new context with needed flags by calling [ints.Dec]'s
// [Context.New] method.  A context's arithmetic flags may not be
// changed once the context is created.  There is a separate set of
// format flags which control the string representation of a [Decimal],
// i.e. its decimal separator and number of fractional positions.
// Format flags may be changed at any time.
type Context struct {
	flags                  *flags
	separator              rune
	fmtSeparator           rune
	integrals, fractionals int8
	pow                    Decimal
	maxInts, maxFrcs       uint64

	// From provides conversion methods returning Decimal-values.
	From *Convert

	// Float provides a ready to use "floats-decimal-calculator".
	Float *Floats

	// Max is the maximal representable Decimal value having a context's
	// set arithmetic fractionals.
	Max Decimal
}

// Dec returns the default context whose arithmetic flags default to
// [DOT_SEPARATOR] | [SIX_FRACTIONALS], i.e. the string [Decimal]
// conversion expects a string with a dot decimal separator and the
// returned [Decimal]'s last six positions are interpreted as its
// fractionals.  The format flags default to a [COMMA_SEPARATOR] |
// [TWO_FRACTIONALS], i.e. a [Decimal.Str] representation's fractionals
// are truncated at the second position and a comma is used as decimal
// separator.  Respectively [Decimal.Rnd] "rounds to even" to the second
// position.  In the unlikely case that there are more format
// fractionals as arithmetic fractionals the string representation is
// accordingly padded with zeros.
//
// Note you cannot change a contexts arithmetics flags.  Create a new
// context with different arithmetics flags by using Dec's [Context.New]
// method.
var Dec = func() *Context {
	cntx := &Context{}
	cntx.flags = newFlags(cntx)
	initialize(cntx)
	return cntx
}()

// New creates a new Context-instance with provided arithmetic and
// format flags.  Using the DEFAULTS flag for arithmetic or format Flags
// copies the respective flag set of given context.  In general if
// a fractionals or a separator flag is omitted the respective flag of given
// context is used.  Is more than one fractional or separator flag given
// only one of them is used and it is undefined which one.
func (c *Context) New(art, fmt Flags) *Context {
	if c == nil || c.flags == nil {
		return Dec.New(art, fmt)
	}
	cntx := &Context{flags: c.flags.copy(art, fmt)}
	initialize(cntx)
	return cntx
}

func initialize(c *Context) {
	c.separator = flagsToSeparator[c.flags.art&ffSeparators]
	c.fmtSeparator = flagsToSeparator[c.flags.fmt&ffSeparators]
	c.integrals, c.fractionals, c.Max, c.maxInts, c.maxFrcs =
		c.fractionalProperties(flagsToFractionals[c.flags.art&ffFractionals])
	c.pow = Decimal(math.Pow10(int(c.fractionals)))
	c.flags.cntx = c
	c.From = &Convert{cntx: c}
	c.Float = &Floats{cntx: c}
}

func (c *Context) fractionalProperties(n int) (
	ii, ff int8, max Decimal, imx, fmx uint64,
) {
	if ii, ok := inits[n]; ok {
		return ii.ii, ii.ff, ii.max, ii.imx, ii.fmx
	}
	vv := nInit{ff: int8(n)}
	expFrc := uint64(math.Pow10(int(vv.ff)))
	var maxDifFrc uint64 = math.MaxUint64 - math.MaxUint64%expFrc
	if maxDifFrc%(expFrc*10) == 0 {
		panic("intdec: context: properties: unexpected zero-position")
	}
	vv.max = Decimal(maxDifFrc - 1)
	vv.imx = (maxDifFrc / expFrc) - 1
	vv.fmx = uint64(vv.max) % expFrc
	vv.ii = int8(len(strconv.Itoa(int(vv.imx))))
	inits[n] = vv
	return vv.ii, vv.ff, vv.max, vv.imx, vv.fmx
}

// nInit represents the number of integrals, of fractionals, the maximum
// unit64 value, the maximum integral number and the maximum fractionals
// uint64.
type nInit struct {
	ii, ff   int8
	max      Decimal
	imx, fmx uint64
}

// inits stores initial values calculated depending on the number of
// fractionals; i.e. they are calculated only once since they are
// constant.
var inits = map[int]nInit{}

// SetFmt sets given context's format flags.  Is more than one fractional
// or separator flag given only one of them is used and it is undefined
// which one.
func (c *Context) SetFmt(ff Flags) {
	c.flags.fmt.set(ff)
	c.fmtSeparator = flagsToSeparator[c.flags.fmt&ffSeparators]
}

// Add adds given decimals and returns their sum.  Add fails if the
// result overflows given context's *Max* property.
func (c *Context) Add(a, b Decimal) (Decimal, error) {
	if a > c.Max-b {
		return 0, ErrOverflow
	}
	return a + b, nil
}

// MAdd is the 'Must'-version of [Context.Add] which panics if
// corresponding Add-call fails.
func (c *Context) MAdd(a, b Decimal) Decimal {
	sum, err := c.Add(a, b)
	if err != nil {
		panic(err)
	}
	return sum
}

// Sub subtracts given decimal b from a and returns their difference.
// Sub fails with an overflow error if b is  greater than a.
func (c *Context) Sub(a, b Decimal) (Decimal, error) {
	if b > a {
		return 0, ErrOverflow
	}
	return a - b, nil
}

// MSub is the 'must'-version of [Context.Sub] which panics if
// corresponding Sub-call fails.
func (c *Context) MSub(a, b Decimal) Decimal {
	diff, err := c.Sub(a, b)
	if err != nil {
		panic(err)
	}
	return diff
}

// Mult multiplies given decimals and returns their product.  Mult fails
// if the product is greater than Max of given Context.
func (c *Context) Mult(a, b Decimal) (Decimal, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}
	bInt := b / c.pow
	prodABInts := a * bInt
	if prodABInts > 0 && (prodABInts > c.Max || prodABInts/bInt != a) {
		return 0, ErrOverflow
	}
	bFrc := b % c.pow
	if bFrc == 0 {
		return prodABInts, nil
	}
	prodAIntsBFrc := (a / c.pow) * bFrc
	if prodABInts > c.Max-prodAIntsBFrc {
		return 0, ErrOverflow
	}
	prodABWithoutProdFractionals := prodABInts + prodAIntsBFrc
	prodAFrcBFrc := ((a % c.pow) * bFrc) / c.pow
	if prodABWithoutProdFractionals > c.Max-prodAFrcBFrc {
		return 0, ErrOverflow
	}
	return prodABWithoutProdFractionals + prodAFrcBFrc, nil
}

// MMult is the 'Must'-variant of [Context.Mult] which panics if
// corresponding Mult-call fails.
func (c *Context) MMult(a, b Decimal) Decimal {
	prd, err := c.Mult(a, b)
	if err != nil {
		panic(err)
	}
	return prd
}

// Div divides a by b and returns resulting quotient.  Div fails if it
// overflows (i.e. a is "big" and 0 < b < 1) or if b is zero.
func (c *Context) Div(a, b Decimal) (Decimal, error) {
	if b == 0 {
		return 0, ErrDividedByZero
	}
	if a == 0 {
		return 0, nil
	}
	if b/c.pow == 0 { // b <= 0
		fkt := c.pow * 10 / b
		rsl := a * fkt
		if rsl/fkt != a {
			return 0, ErrOverflow
		}
	}
	num := a * c.pow
	if num/c.pow != a {
		return c.bigIntDiv(a, b)
	}
	return num / b, nil
}

// MDiv is the 'Must'-variant of [Context.Div] which panics if
// corresponding Div-call fails.
func (c *Context) MDiv(a, b Decimal) Decimal {
	qut, err := c.Div(a, b)
	if err != nil {
		panic(err)
	}
	return qut
}

const maxInt64 = Decimal(math.MaxInt64)

func (c *Context) bigIntDiv(a, b Decimal) (Decimal, error) {
	bigA := (&big.Int{}).SetUint64(uint64(a))
	bigA.Mul(bigA, big.NewInt(int64(c.pow)))
	bigA.Div(bigA, (&big.Int{}).SetUint64(uint64(b)))
	if !bigA.IsUint64() {
		return 0, ErrOverflow
	}
	return Decimal(bigA.Uint64()), nil
}
