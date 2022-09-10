// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

import (
	"math"
	"strconv"
	"testing"

	. "github.com/slukits/gounit"
)

type Default struct{ Suite }

func (s *Default) SetUp(t *T) { t.Parallel() }

func (s *Default) Has_six_fractionals(t *T) {
	t.Eq(int8(6), UDec.fractionals)
}

func (s *Default) Has_max_value_length_minus_six_integrals(t *T) {
	t.Eq(int8(len(strconv.FormatUint(uint64(UDec.Max), 10)))-6,
		UDec.integrals)
}

func (s *Default) Max_value_is_18446744073708999999(t *T) {
	t.Eq(UDecimal(18446744073708999999), UDec.Max)
}

func (s *Default) Max_integrals_is_18446744073708(t *T) {
	t.Eq(uint64(18446744073708), UDec.maxInts)
}

func (s *Default) Max_fractionals_is_999999(t *T) {
	t.Eq(uint64(999999), UDec.maxFrcs)
}

func (s *Default) Has_dot_separator(t *T) {
	t.Eq('.', UDec.separator)
}

func TestDefault(t *testing.T) {
	t.Parallel()
	Run(&Default{}, t)
}

type context struct{ Suite }

func (s *context) SetUp(t *T) { t.Parallel() }

func (s *context) Has_set_flags(t *T) {
	dec := UDec.New(COMMA_SEPARATOR|ONE_FRACTIONAL,
		COMMA_SEPARATOR|ONE_FRACTIONAL)
	t.Eq(',', dec.separator)
	t.Eq(1, dec.flags.Fractionals())
	t.Eq(',', dec.fmtSeparator)
	t.Eq(1, dec.flags.FmtFractionals())
	dec.SetFmt(DOT_SEPARATOR)
	t.Eq('.', dec.fmtSeparator)
}

func (s *context) Fractionals_default_to_six(t *T) {
	cntx := UDec.New(DEFAULTS, DEFAULTS)
	t.Eq(6, cntx.flags.Fractionals())
}

func (s *context) Provides_default_copy_if_new_from_nil_or_zero(t *T) {
	var fx UContext
	got := fx.New(DEFAULTS, DEFAULTS)
	t.Eq(UDec.separator, got.separator)
	t.Eq(UDec.fmtSeparator, got.fmtSeparator)
	t.Eq(UDec.flags.Fractionals(), got.flags.Fractionals())
	t.Eq(UDec.flags.FmtFractionals(), got.flags.FmtFractionals())
	got = (*UContext)(nil).New(DEFAULTS, DEFAULTS)
	t.Eq(UDec.separator, got.separator)
	t.Eq(UDec.fmtSeparator, got.fmtSeparator)
	t.Eq(UDec.flags.Fractionals(), got.flags.Fractionals())
	t.Eq(UDec.flags.FmtFractionals(), got.flags.FmtFractionals())
}

func TestContext(t *testing.T) {
	t.Parallel()
	Run(&context{}, t)
}

type Operation struct{ Suite }

func (s *Operation) SetUp(t *T) { t.Parallel() }

func (s *Operation) Overflows_if_addition_exceeds_max(t *T) {
	_, err := UDec.Float.Add(float64(UDec.maxInts), 1)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { UDec.Float.MAdd(float64(UDec.maxInts), 1) })
}

func (s *Operation) Overflows_if_difference_below_zero(t *T) {
	_, err := UDec.Float.Sub(1, 1.1)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { UDec.Float.MSub(1, 1.1) })
}

func (s *Operation) Overflows_if_product_exceeds_max(t *T) {
	_, err := UDec.Float.Mult(float64(UDec.maxInts), 2.1)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { UDec.Float.MMult(float64(UDec.maxInts), 2.1) })
	_, err = UDec.Float.Mult(float64(UDec.maxInts), 1.1)
	t.ErrIs(err, ErrOverflow)
	t.True(math.MaxUint64 > uint64(UDec.Max))
	_, err = UDec.Float.Mult(float64(math.MaxUint64), 1.1)
	t.ErrIs(err, ErrOverflow)
	halfMax, err := UDec.Div(UDecimal(UDec.Max), 2000000)
	t.FatalOn(err)
	_, err = UDec.Float.Mult(float64(halfMax), 2.1)
	t.ErrIs(err, ErrOverflow)
}

func (s *Operation) Overflows_if_div_numerator_exceeds_max(t *T) {
	_, err := UDec.Float.Div(float64(UDec.Max), 1.25)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { UDec.Float.MDiv(float64(UDec.Max), 1.25) })
	_, err = UDec.Float.Div(float64(UDec.maxInts), 0.0005)
	t.ErrIs(err, ErrOverflow)
}

func (s *Operation) Adding_a_and_b_returns_their_sum(t *T) {
	t.Eq(UDecimal(3580100), UDec.Float.MAdd(1.2345, 2.3456))
}

func (s *Operation) Subtracting_b_from_a_returns_their_difference(t *T) {
	t.Eq(UDecimal(1111100), UDec.Float.MSub(2.3456, 1.2345))
}

func (s *Operation) Multiplying_a_and_b_returns_their_product(t *T) {
	t.Eq(UDecimal(2895643), UDec.Float.MMult(2.3456, 1.2345))
}

func (s *Operation) Multiply_truncates_its_product(t *T) {
	t.Eq(UDecimal(1523990), UDec.Float.MMult(1.2345, 1.2345)) // truncates 25
}

func (s *Operation) Dividing_a_and_b_returns_their_quotient(t *T) {
	t.Eq(UDecimal(1896000), UDec.Float.MDiv(2.37, 1.25))
	t.Eq(UDec.From.MFloat(0.2), UDec.Float.MMult(
		0.5, float64(UDec.Float.MDiv(0.2, 0.5))/float64(UDec.pow)))
	t.Eq(UDecimal(9223372036854499999),
		UDec.MDiv(UDec.Max, UDec.From.MFloat(2)))
}

func (s *Operation) Division_truncates_its_quotient(t *T) {
	t.Eq(UDecimal(526282), UDec.Float.MDiv(1.2345, 2.3457))
}

func (s *Operation) Dividing_by_zero_fails(t *T) {
	_, err := UDec.Div(UDec.From.MFloat(1.2345), UDecimal(0))
	t.ErrIs(err, ErrDividedByZero)
}

func (s *Operation) Dividing_of_zero_is_zero(t *T) {
	t.Eq(UDecimal(0), UDec.MDiv(0, UDec.From.MFloat(1.2345)))
}

func (s *Operation) Multiplying_by_zero_is_zero(t *T) {
	t.Eq(UDecimal(0), UDec.MMult(0, UDec.From.MFloat(1.2345)))
	t.Eq(UDecimal(0), UDec.MMult(UDec.From.MFloat(1.2345), 0))
}

func (s *Operation) Zero_is_additive_neutral_element(t *T) {
	a := UDec.From.MFloat(1.2345)
	t.Eq(a, UDec.MAdd(0, a))
	t.Eq(a, UDec.MAdd(a, 0))
}

func (s *Operation) One_is_multiplicative_neutral_element(t *T) {
	x := UDec.From.MFloat(1.2345)
	one := UDec.From.MFloat(1)
	t.Eq(x, UDec.MMult(one, x))
	t.Eq(x, UDec.MMult(x, one))
}

func TestOperation(t *testing.T) {
	t.Parallel()
	Run(&Operation{}, t)
}
