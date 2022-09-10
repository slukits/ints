// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

import (
	"math"
	"strconv"
	"strings"
)

// UConvert provides the [UDecimal]-value creating conversion functions of
// a [UContext]'s From-property.
//
// NOTE the zero value is NOT ready to use; you shouldn't need to create
// a UConvert instance directly.  The main motivation for this type is to
// keep [UContext]'s API slim.  It was made public anyway to have its
// methods documentations shown in the go documentation server.
type UConvert struct {
	cntx *UContext
}

// Str converts given string to a [UDecimal]-value.  If given string
// doesn't contain the decimal separator s + ".0" is assumed, i.e. given
// string is converted to an int which is padded by fractional zeros.
// Otherwise the string is split at the decimal mark and the first
// substring is interpreted as the integrals while the second is
// interpreted as the fractionals.  Str fails if resulting Decimal is
// greater than associated [UContext]'s Max-property while superfluous
// fractionals are truncated.  Str also fails if
// integrals- or fractionals-parsing fails.
func (c UConvert) Str(s string) (_ UDecimal, err error) {
	if s == "" {
		return 0, nil
	}
	if !strings.ContainsRune(s, c.cntx.separator) {
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return c.combine(i, 0, 0)
	}
	ifStr := strings.SplitN(s, string(c.cntx.separator), 2)
	var iUint, fUint uint64
	if ifStr[0] == "" {
		iUint = 0
	} else {
		iUint, err = strconv.ParseUint(ifStr[0], 10, 64)
		if err != nil {
			return 0, err
		}
	}
	if ifStr[1] == "" {
		fUint = 0
	} else {
		fUint, err = strconv.ParseUint(ifStr[1], 10, 64)
		if err != nil {
			return 0, err
		}
	}
	return c.combine(iUint, fUint, len(ifStr[1]))
}

// MStr is the "must"-variant of [UConvert.Str] which panics if
// corresponding Str-call fails.
func (c UConvert) MStr(s string) UDecimal {
	v, err := c.Str(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Float converts given float to a [UDecimal]-value by extracting its
// integrals and fractionals.  It overflows if resulting Decimal is
// greater than associated [UContext]'s Max-property while superfluous
// fractionals are truncated.
func (c UConvert) Float(f float64) (UDecimal, error) {
	unit := math.Pow10(int(c.cntx.fractionals + 1))
	fInt, fFrc := math.Modf(f)
	fFrc = math.Round(fFrc*unit) / 10
	return c.combine(
		uint64(fInt),
		uint64(fFrc),
		int(c.cntx.fractionals),
	)
}

// MFloat is the "Must"-variant of [UConvert.Float] which panics if
// corresponding Float-call fails.
func (c UConvert) MFloat(f float64) UDecimal {
	v, err := c.Float(f)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints converts given integrals, fractionals and leading zeros lz of
// the fractionals to a [UDecimal]-value.  It overflows if resulting
// Decimal is greater than associated [UContext]'s Max-property while
// superfluous fractionals are truncated.
func (c UConvert) Ints(
	integrals, fractionals uint64, lz int,
) (UDecimal, error) {
	if integrals == 0 && fractionals == 0 {
		return 0, nil
	}
	len := c.len(fractionals)
	if lz >= int(c.cntx.fractionals) {
		return 0, ErrOverflow
	}
	return c.combine(integrals, fractionals, len+lz)
}

// MInts is the "Must"-variant of [UConvert.Ints] which panics if
// corresponding Ints-call fails.
func (c UConvert) MInts(integrals, fractionals uint64, lz int) UDecimal {
	v, err := c.Ints(integrals, fractionals, lz)
	if err != nil {
		panic(err)
	}
	return v
}

// Cntx converts given decimal of given context to a decimal of
// receiving context.  Cntx fails if the conversion overflows.  If
// receiving context has less fractionals than given then the
// superfluous fractions are rounded evenly off (see UDecimal.Rnd).
func (c UConvert) Cntx(d UDecimal, cx *UContext) (UDecimal, error) {
	if c.cntx.flags.Fractionals() < cx.flags.Fractionals() {
		return c.rnd(d, cx)
	}
	pow := UDecimal(math.Pow10(
		c.cntx.flags.Fractionals() - cx.flags.Fractionals()))
	cnv := d * pow
	if cnv/pow != d {
		return 0, ErrOverflow
	}
	return cnv, nil
}

// MCntx is the "Must"-variant of [UConvert.Cntx] which panics if
// corresponding Cntx-call fails.
func (c UConvert) MCntx(d UDecimal, cx *UContext) UDecimal {
	d, err := c.Cntx(d, cx)
	if err != nil {
		panic(err)
	}
	return d
}

// rnd rounds given decimal evenly down to an decimal of given
// converters context with lesser fractionals.  NOTE there is no test
// that given converter's context has less fractionals than given
// context.  Since Cntx is the only caller of rnd it is not really
// possible to test this case yet hence not test.
func (c UConvert) rnd(d UDecimal, cx *UContext) (UDecimal, error) {
	// to must be smaller than cx's fractionals
	to := c.cntx.flags.Fractionals()
	off := cx.flags.Fractionals() - to
	pow, rest := UDecimal(math.Pow10(off)), UDecimal(0)
	rnd, d := d%pow, d/pow
	if off > 1 {
		rest = rnd % (pow / 10)
		rnd /= pow / 10
	}
	if rnd > 5 || (rnd == 5 && (rest > 0 || d%2 == 1)) {
		d += 1
	}
	return d, nil
}

func (c UConvert) len(a uint64) int {
	len := 0
	for a != 0 {
		len++
		a /= 10
	}
	return len
}

func (c UConvert) cut(fractionals uint64, lenWithZeros int) int {
	len := c.len(fractionals)
	if lenWithZeros > len {
		len = lenWithZeros
	}
	return len - int(c.cntx.fractionals)
}

func (c UConvert) combine(
	integrals, fractionals uint64, lenWithZeros int,
) (UDecimal, error) {
	if integrals > c.cntx.maxInts {
		return 0, ErrOverflow
	}
	if integrals == 0 && fractionals == 0 {
		return 0, nil
	}
	if fractionals == 0 {
		return UDecimal(integrals) * c.cntx.pow, nil
	}
	if lenWithZeros < int(c.cntx.fractionals) {
		fractionals *= uint64(
			math.Pow10(int(c.cntx.fractionals) - lenWithZeros))
	}
	if fractionals > c.cntx.maxFrcs { // truncate superfluous fractionals
		fractionals = fractionals / uint64(math.Pow10(c.cut(
			fractionals, lenWithZeros)))
	}
	if integrals == 0 {
		return UDecimal(fractionals), nil
	}
	return UDecimal(integrals*uint64(c.cntx.pow) + fractionals), nil
}
