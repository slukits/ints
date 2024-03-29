// Copyright 2022 Stephan.Lukits@gmail.com. All rights reserved.
// Use of this source code is governed by the MIT license.

package ints

import (
	"fmt"
	"math"
	"strings"
)

// UDecimal value represents together with a [UContext]-instance a
// decimal value whose last n positions are interpreted as the decimal's
// fractionals and the remaining positions are interpreted as the
// decimal's integral part.  n is defined by the context [n]_FRACTIONALS
// arithmetic-flag.  NOTE a UContext defines the interpretation of an
// UDecimal's underlying integer.  Changing the context for a value
// doesn't change the value but its interpretation under the new context
// which is usually not what you want.  I.e. using an UDecimal value v'
// with a different context c but the creating context c' can be only
// done safely if c and c' have the same fractionals arithmetic flag or
// v' is converted to a c value by v := c.From.Cntx(v', c').
type UDecimal uint64

// Len returns a decimal's number of digits.
func (d UDecimal) Len() int {
	l := 0
	for d != 0 {
		l++
		d /= 10
	}
	return l
}

// Fractionals returns a decimal's fractional part relative to provided
// contexts ..._FRACTIONALS arithmetic-flag (see [ints.Flags]).
func (d UDecimal) Fractionals(c *UContext) UDecimal {
	return d % c.pow
}

// Fractionals returns a decimal's integral part relative to provided
// contexts ..._FRACTIONALS arithmetic-flag.
func (d UDecimal) Integrals(c *UContext) UDecimal {
	return d / c.pow
}

// Str returns a string representation of given value relative to the
// settings of given context.  If ..._FRACTIONALS of given context's
// format flags is smaller than its corresponding  arithmetic flag the
// fractionals are accordingly truncated; is it bigger zeros are
// accordingly padded.  See [UDecimal.Rnd] for a rounded string
// representation
func (d UDecimal) Str(c *UContext) string {
	ii, ff, prePad := d.Integrals(c), d.Fractionals(c), ""
	cf, sf := c.flags.Fractionals(), c.flags.FmtFractionals()
	ln := ff.Len()
	if ln < cf { // need to pad zeros
		prePad = strings.Repeat("0", cf-ln)
	}
	if sf >= cf {
		sufPad := strings.Repeat("0", sf-cf)
		return fmt.Sprintf("%d%c%s%d%s", ii, c.fmtSeparator, prePad, ff, sufPad)
	}
	for i := 0; i < cf-sf; i++ {
		ff /= 10
	}
	if ff == 0 {
		return fmt.Sprintf("%d%c%s", ii, c.fmtSeparator,
			strings.Repeat("0", sf))
	}
	return fmt.Sprintf("%d%c%s%d", ii, c.fmtSeparator, prePad, ff)
}

// Rnd returns a rounded to even string representation of given value
// iff the ..._FRACTIONALS of given Context's format flags is smaller
// than its corresponding arithmetic flag.
//
// "round to even" is implemented as follows: let v be the value to
// round, d the first digit to round away and v' the decimal until d
// (e.g. v = 2.4358, format-flag is [TWO_FRACTIONALS]: v' = 2.43, d = 5)
// then the rounded value of v is for:
//   - d < 5: v'
//   - d > 5: v'+ 0.x1 with x is a sequence of 0s one less than
//     format-fractionals
//   - d == 5 and there are non-null positions after d: v'+0.x1
//   - d == 5 and there are no non-null positions after d:
//     v' is returned if v' is even; otherwise v'+0.x1 is returned
func (d UDecimal) Rnd(c *UContext) string {
	nf := c.flags.FmtFractionals()
	off := c.flags.Fractionals() - nf
	if off <= 0 {
		return d.Str(c)
	}
	pow, rest := UDecimal(math.Pow10(off)), UDecimal(0)
	rnd, d := d%pow, d/pow
	if off > 1 {
		rest = rnd % (pow / 10)
		rnd /= pow / 10
	}
	if rnd > 5 || (rnd == 5 && (rest > 0 || d%2 == 1)) {
		d += 1
	}
	pow = UDecimal(math.Pow10(nf))
	ff, _nf := d%pow, 0
	if ff == 0 {
		return fmt.Sprintf("%d%c%s",
			d/pow, c.separator, strings.Repeat("0", nf))
	}
	for ff != 0 {
		_nf++
		ff /= 10
	}
	if _nf < nf {
		return fmt.Sprintf("%d%c%s%d",
			d/pow, c.separator, strings.Repeat("0", nf-_nf), d%pow)
	}
	return fmt.Sprintf("%d%c%d", d/pow, c.separator, d%pow)
}
