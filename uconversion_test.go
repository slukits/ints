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
	_, err := UDec.From.Float(float64(UDec.maxInts + 1))
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { UDec.From.MFloat(float64(UDec.maxInts + 1)) })
	_, err = UDec.From.Str(fmt.Sprintf("%d", UDec.maxInts+1))
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { UDec.From.MStr(fmt.Sprintf("%d", UDec.maxInts+1)) })
	_, err = UDec.From.Ints(UDec.maxInts+1, 0, 0)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { UDec.From.MInts(UDec.maxInts+1, 0, 0) })
	_, err = UDec.From.Ints(UDec.maxInts, 0, int(UDec.fractionals)+1)
	t.ErrIs(err, ErrOverflow)
	_, err = UDec.From.Ints(UDec.maxInts, 0, int(UDec.fractionals)+1)
	t.ErrIs(err, ErrOverflow)
	_, err = UDec.Float.Sub(1, float64(UDec.maxInts+1))
	t.ErrIs(err, ErrOverflow)
}

func (s *Conversion) Truncates_superfluous_fractionals(t *T) {
	dec := UDec.New(FOUR_FRACTIONALS, DEFAULTS)
	t.Eq(UDecimal(12345), dec.From.MFloat(1.23456))
	t.Eq(UDecimal(12345), dec.From.MStr("1.23456"))
	t.Eq(UDecimal(12345), dec.From.MInts(1, 23456, 0))
}

func (s *Conversion) Pads_missing_fractionals(t *T) {
	t.Eq(UDecimal(4020000), UDec.From.MFloat(4.02))
	t.Eq(UDecimal(4200000), UDec.From.MStr("4.2"))
	t.Eq(UDecimal(4200000), UDec.From.MInts(4, 2, 0))
}

func (s *Conversion) Preserves_leading_fractional_zeros(t *T) {
	t.Eq(UDecimal(4020000), UDec.From.MFloat(4.0200))
	t.Eq(UDecimal(4002000), UDec.From.MStr("4.002"))
	t.Eq(UDecimal(4000200), UDec.From.MInts(4, 2, 3))
}

func (s *Conversion) From_string_fails_if_uint_parsing_fails(t *T) {
	_, err := UDec.From.Str("abc")
	t.ErrMatched(err, "invalid syntax")
	_, err = UDec.From.Str("abc.42")
	t.ErrMatched(err, "invalid syntax")
	_, err = UDec.From.Str("42.abc")
	t.ErrMatched(err, "invalid syntax")
}

func (s *Conversion) Recognizes_variations_of_zero_strings(t *T) {
	for _, z := range []string{"", ".", "0.", ".0", "0.0"} {
		a, err := UDec.From.Str(z)
		t.FatalOn(err)
		t.Eq(UDecimal(0), a)
	}
	a, err := UDec.From.Ints(0, 0, 0)
	t.FatalOn(err)
	t.Eq(UDecimal(0), a)
	for _, z := range []float64{0, 0., .0, 0.0} {
		a, err := UDec.From.Float(z)
		t.FatalOn(err)
		t.Eq(UDecimal(0), a)
	}
}

func (s *Conversion) Converts_to_a_different_context(t *T) {
	dec4 := UDec.New(FOUR_FRACTIONALS, DEFAULTS)
	dec2 := UDec.New(TWO_FRACTIONALS, DEFAULTS)
	d41, d42 := dec4.From.MStr("2.3450"), dec4.From.MStr("2.3550")
	d21, d22 := dec2.From.MCntx(d41, dec4), dec2.From.MCntx(d42, dec4)
	t.Eq(UDecimal(234), d21) // round to even
	t.Eq(UDecimal(236), d22) // conversion
	d41, d42 = dec4.From.MCntx(d21, dec2), dec4.From.MCntx(d22, dec2)
	t.Eq(UDecimal(23400), d41)
	t.Eq(UDecimal(23600), d42)
}

func (s *Conversion) Fails_if_context_value_conversion_overflows(t *T) {
	dec4 := UDec.New(FOUR_FRACTIONALS, DEFAULTS)
	dec2 := UDec.New(TWO_FRACTIONALS, DEFAULTS)
	_, err := dec4.From.Cntx(dec2.Max, dec2)
	t.ErrIs(err, ErrOverflow)
}

func (s *Conversion) Panics_if_must_context_conversion_overflows(t *T) {
	dec4 := UDec.New(FOUR_FRACTIONALS, DEFAULTS)
	dec2 := UDec.New(TWO_FRACTIONALS, DEFAULTS)
	t.Panics(func() { dec4.From.MCntx(dec2.Max, dec2) })
}

func TestConversion(t *testing.T) {
	t.Parallel()
	Run(&Conversion{}, t)
}
