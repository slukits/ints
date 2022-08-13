// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type Conversion struct{ Suite }

func (s *Conversion) SetUp(t *T) { t.Parallel() }

func (s *Conversion) Overflows_if_integral_part_is_to_big(t *T) {
	_, err := Dec.From.Float(float64(Dec.maxInts + 1))
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { Dec.From.MFloat(float64(Dec.maxInts + 1)) })
	_, err = Dec.From.Str(fmt.Sprintf("%d", Dec.maxInts+1))
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { Dec.From.MStr(fmt.Sprintf("%d", Dec.maxInts+1)) })
	_, err = Dec.From.Ints(Dec.maxInts+1, 0, 0)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { Dec.From.MInts(Dec.maxInts+1, 0, 0) })
	_, err = Dec.From.Ints(Dec.maxInts, 0, int(Dec.fractionals)+1)
	t.ErrIs(err, ErrOverflow)
	_, err = Dec.From.Ints(Dec.maxInts, 0, int(Dec.fractionals)+1)
	t.ErrIs(err, ErrOverflow)
	_, err = Dec.Float.Sub(1, float64(Dec.maxInts+1))
	t.ErrIs(err, ErrOverflow)
}

func (s *Conversion) Truncates_superfluous_fractionals(t *T) {
	dec := Dec.New(FOUR_FRACTIONALS, DEFAULTS)
	t.Eq(Decimal(12345), dec.From.MFloat(1.23456))
	t.Eq(Decimal(12345), dec.From.MStr("1.23456"))
	t.Eq(Decimal(12345), dec.From.MInts(1, 23456, 0))
}

func (s *Conversion) Pads_missing_fractionals(t *T) {
	t.Eq(Decimal(4020000), Dec.From.MFloat(4.02))
	t.Eq(Decimal(4200000), Dec.From.MStr("4.2"))
	t.Eq(Decimal(4200000), Dec.From.MInts(4, 2, 0))
}

func (s *Conversion) Preserves_leading_fractional_zeros(t *T) {
	t.Eq(Decimal(4020000), Dec.From.MFloat(4.0200))
	t.Eq(Decimal(4002000), Dec.From.MStr("4.002"))
	t.Eq(Decimal(4000200), Dec.From.MInts(4, 2, 3))
}

func (s *Conversion) From_string_fails_if_uint_parsing_fails(t *T) {
	_, err := Dec.From.Str("abc")
	t.ErrMatched(err, "invalid syntax")
	_, err = Dec.From.Str("abc.42")
	t.ErrMatched(err, "invalid syntax")
	_, err = Dec.From.Str("42.abc")
	t.ErrMatched(err, "invalid syntax")
}

func (s *Conversion) Recognizes_variations_of_zero_strings(t *T) {
	for _, z := range []string{"", ".", "0.", ".0", "0.0"} {
		a, err := Dec.From.Str(z)
		t.FatalOn(err)
		t.Eq(Decimal(0), a)
	}
	a, err := Dec.From.Ints(0, 0, 0)
	t.FatalOn(err)
	t.Eq(Decimal(0), a)
	for _, z := range []float64{0, 0., .0, 0.0} {
		a, err := Dec.From.Float(z)
		t.FatalOn(err)
		t.Eq(Decimal(0), a)
	}
}

func TestConversion(t *testing.T) {
	t.Parallel()
	Run(&Conversion{}, t)
}
