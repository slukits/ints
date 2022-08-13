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
	t.Eq(int8(6), Dec.fractionals)
}

func (s *Default) Has_max_value_length_minus_six_integrals(t *T) {
	t.Eq(int8(len(strconv.FormatUint(uint64(Dec.Max), 10)))-6,
		Dec.integrals)
}

func (s *Default) Max_value_is_18446744073708999999(t *T) {
	t.Eq(Decimal(18446744073708999999), Dec.Max)
}

func (s *Default) Max_integrals_is_18446744073708(t *T) {
	t.Eq(uint64(18446744073708), Dec.maxInts)
}

func (s *Default) Max_fractionals_is_999999(t *T) {
	t.Eq(uint64(999999), Dec.maxFrcs)
}

func (s *Default) Has_dot_separator(t *T) {
	t.Eq('.', Dec.separator)
}

func TestDefault(t *testing.T) {
	t.Parallel()
	Run(&Default{}, t)
}

type context struct{ Suite }

func (s *context) SetUp(t *T) { t.Parallel() }

func (s *context) Has_set_flags(t *T) {
	dec := Dec.New(COMMA_SEPARATOR|ONE_FRACTIONAL,
		COMMA_SEPARATOR|ONE_FRACTIONAL)
	t.Eq(',', dec.separator)
	t.Eq(1, dec.flags.Fractionals())
	t.Eq(',', dec.fmtSeparator)
	t.Eq(1, dec.flags.FmtFractionals())
	dec.SetFmt(DOT_SEPARATOR)
	t.Eq('.', dec.fmtSeparator)
}

func (s *context) Fractionals_default_to_six(t *T) {
	cntx := Dec.New(DEFAULTS, DEFAULTS)
	t.Eq(6, cntx.flags.Fractionals())
}

func (s *context) Provides_default_copy_if_new_from_nil_or_zero(t *T) {
	var fx Context
	got := fx.New(DEFAULTS, DEFAULTS)
	t.Eq(Dec.separator, got.separator)
	t.Eq(Dec.fmtSeparator, got.fmtSeparator)
	t.Eq(Dec.flags.Fractionals(), got.flags.Fractionals())
	t.Eq(Dec.flags.FmtFractionals(), got.flags.FmtFractionals())
	got = (*Context)(nil).New(DEFAULTS, DEFAULTS)
	t.Eq(Dec.separator, got.separator)
	t.Eq(Dec.fmtSeparator, got.fmtSeparator)
	t.Eq(Dec.flags.Fractionals(), got.flags.Fractionals())
	t.Eq(Dec.flags.FmtFractionals(), got.flags.FmtFractionals())
}

func TestContext(t *testing.T) {
	t.Parallel()
	Run(&context{}, t)
}

type Operation struct{ Suite }

func (s *Operation) SetUp(t *T) { t.Parallel() }

func (s *Operation) Overflows_if_addition_exceeds_max(t *T) {
	_, err := Dec.Float.Add(float64(Dec.maxInts), 1)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { Dec.Float.MAdd(float64(Dec.maxInts), 1) })
}

func (s *Operation) Overflows_if_difference_below_zero(t *T) {
	_, err := Dec.Float.Sub(1, 1.1)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { Dec.Float.MSub(1, 1.1) })
}

func (s *Operation) Overflows_if_product_exceeds_max(t *T) {
	_, err := Dec.Float.Mult(float64(Dec.maxInts), 2.1)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { Dec.Float.MMult(float64(Dec.maxInts), 2.1) })
	_, err = Dec.Float.Mult(float64(Dec.maxInts), 1.1)
	t.ErrIs(err, ErrOverflow)
	t.True(math.MaxUint64 > uint64(Dec.Max))
	_, err = Dec.Float.Mult(float64(math.MaxUint64), 1.1)
	t.ErrIs(err, ErrOverflow)
	halfMax, err := Dec.Div(Decimal(Dec.Max), 2000000)
	t.FatalOn(err)
	_, err = Dec.Float.Mult(float64(halfMax), 2.1)
	t.ErrIs(err, ErrOverflow)
}

func (s *Operation) Overflows_if_div_numerator_exceeds_max(t *T) {
	_, err := Dec.Float.Div(float64(Dec.Max), 1.25)
	t.ErrIs(err, ErrOverflow)
	t.Panics(func() { Dec.Float.MDiv(float64(Dec.Max), 1.25) })
	_, err = Dec.Float.Div(float64(Dec.maxInts), 0.0005)
	t.ErrIs(err, ErrOverflow)
}

func (s *Operation) Adding_a_and_b_returns_their_sum(t *T) {
	t.Eq(Decimal(3580100), Dec.Float.MAdd(1.2345, 2.3456))
}

func (s *Operation) Subtracting_b_from_a_returns_their_difference(t *T) {
	t.Eq(Decimal(1111100), Dec.Float.MSub(2.3456, 1.2345))
}

func (s *Operation) Multiplying_a_and_b_returns_their_product(t *T) {
	t.Eq(Decimal(2895643), Dec.Float.MMult(2.3456, 1.2345))
}

func (s *Operation) Multiply_truncates_its_product(t *T) {
	t.Eq(Decimal(1523990), Dec.Float.MMult(1.2345, 1.2345)) // truncates 25
}

func (s *Operation) Dividing_a_and_b_returns_their_quotient(t *T) {
	t.Eq(Decimal(1896000), Dec.Float.MDiv(2.37, 1.25))
	t.Eq(Dec.From.MFloat(0.2), Dec.Float.MMult(
		0.5, float64(Dec.Float.MDiv(0.2, 0.5))/float64(Dec.pow)))
	t.Eq(Decimal(9223372036854499999),
		Dec.MDiv(Dec.Max, Dec.From.MFloat(2)))
}

func (s *Operation) Division_truncates_its_quotient(t *T) {
	t.Eq(Decimal(526282), Dec.Float.MDiv(1.2345, 2.3457))
}

func (s *Operation) Dividing_by_zero_fails(t *T) {
	_, err := Dec.Div(Dec.From.MFloat(1.2345), Decimal(0))
	t.ErrIs(err, ErrDividedByZero)
}

func (s *Operation) Dividing_of_zero_is_zero(t *T) {
	t.Eq(Decimal(0), Dec.MDiv(0, Dec.From.MFloat(1.2345)))
}

func (s *Operation) Multiplying_by_zero_is_zero(t *T) {
	t.Eq(Decimal(0), Dec.MMult(0, Dec.From.MFloat(1.2345)))
	t.Eq(Decimal(0), Dec.MMult(Dec.From.MFloat(1.2345), 0))
}

func (s *Operation) Zero_is_additive_neutral_element(t *T) {
	a := Dec.From.MFloat(1.2345)
	t.Eq(a, Dec.MAdd(0, a))
	t.Eq(a, Dec.MAdd(a, 0))
}

func (s *Operation) One_is_multiplicative_neutral_element(t *T) {
	x := Dec.From.MFloat(1.2345)
	one := Dec.From.MFloat(1)
	t.Eq(x, Dec.MMult(one, x))
	t.Eq(x, Dec.MMult(x, one))
}

func TestOperation(t *testing.T) {
	t.Parallel()
	Run(&Operation{}, t)
}
