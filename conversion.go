// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

import (
	"math"
	"strconv"
	"strings"
)

// Convert provides the [Decimal]-value creating conversion functions of
// a [Context]'s From-property.
//
// NOTE the zero value is NOT ready to use; you shouldn't need to create
// a Convert instance directly.  The main motivation for this type is to
// keep [Context]'s API slim.  It was made public anyway to have its
// methods documentations shown in the go documentation server.
type Convert struct {
	cntx *Context
}

// Str converts given string to a [Decimal]-value.  If given string
// doesn't contain the decimal separator s + ".0" is assumed, i.e. given
// string is converted to an int which is padded by fractional zeros.
// Otherwise the string is split at the decimal mark and the first
// substring is interpreted as the integrals while the second is
// interpreted as the fractionals.  Str fails if resulting Decimal is
// greater than associated [Context]'s Max-property while superfluous
// fractionals are truncated.  Str also fails if
// integrals- or fractionals-parsing fails.
func (c Convert) Str(s string) (_ Decimal, err error) {
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

// MStr is the "must"-variant of [Convert.Str] which panics if
// corresponding Str-call fails.
func (c Convert) MStr(s string) Decimal {
	v, err := c.Str(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Float converts given float to a [Decimal]-value by extracting its
// integrals and fractionals.  It overflows if resulting Decimal is
// greater than associated [Context]'s Max-property while superfluous
// fractionals are truncated.
func (c Convert) Float(f float64) (Decimal, error) {
	unit := math.Pow10(int(c.cntx.fractionals + 1))
	fInt, fFrc := math.Modf(f)
	fFrc = math.Round(fFrc*unit) / 10
	return c.combine(
		uint64(fInt),
		uint64(fFrc),
		int(c.cntx.fractionals),
	)
}

// MFloat is the "Must"-variant of [Convert.Float] which panics if
// corresponding Float-call fails.
func (c Convert) MFloat(f float64) Decimal {
	v, err := c.Float(f)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints converts given integrals, fractionals and leading zeros lz of
// the fractionals to a [Decimal]-value.  It overflows if resulting
// Decimal is greater than associated [Context]'s Max-property while
// superfluous fractionals are truncated.
func (c Convert) Ints(
	integrals, fractionals uint64, lz int,
) (Decimal, error) {
	if integrals == 0 && fractionals == 0 {
		return 0, nil
	}
	len := c.len(fractionals)
	if lz >= int(c.cntx.fractionals) {
		return 0, ErrOverflow
	}
	return c.combine(integrals, fractionals, len+lz)
}

// MInts is the "Must"-variant of [Convert.Ints] which panics if
// corresponding Ints-call fails.
func (c Convert) MInts(integrals, fractionals uint64, lz int) Decimal {
	v, err := c.Ints(integrals, fractionals, lz)
	if err != nil {
		panic(err)
	}
	return v
}

func (c Convert) len(a uint64) int {
	len := 0
	for a != 0 {
		len++
		a /= 10
	}
	return len
}

func (c Convert) cut(fractionals uint64, lenWithZeros int) int {
	len := c.len(fractionals)
	if lenWithZeros > len {
		len = lenWithZeros
	}
	return len - int(c.cntx.fractionals)
}

func (c Convert) combine(
	integrals, fractionals uint64, lenWithZeros int,
) (Decimal, error) {
	if integrals > c.cntx.maxInts {
		return 0, ErrOverflow
	}
	if integrals == 0 && fractionals == 0 {
		return 0, nil
	}
	if fractionals == 0 {
		return Decimal(integrals) * c.cntx.pow, nil
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
		return Decimal(fractionals), nil
	}
	return Decimal(integrals*uint64(c.cntx.pow) + fractionals), nil
}
