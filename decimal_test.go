// Copyright 2022 Stephan.Lukits@gmail.com. All rights reserved.
// Use of this source code is governed by the MIT license.

package ints

import (
	"testing"

	. "github.com/slukits/gounit"
)

type decimal struct{ Suite }

func (s *decimal) SetUp(t *T) { t.Parallel() }

func (s *decimal) Provides_integral_part(t *T) {
	t.Eq(Decimal(42), Dec.From.MFloat(42.42).Integrals(Dec))
}

func (s *decimal) Provides_fractional_part(t *T) {
	t.Eq(Decimal(420000), Dec.From.MFloat(42.42).Fractionals(Dec))
	t.Eq(Decimal(42), Dec.From.MFloat(42.000042).Fractionals(Dec))
}

func (s *decimal) String_uses_set_decimal_mark(t *T) {
	t.Eq("3,20", Dec.From.MFloat(3.2).Str(Dec))
	dec := Dec.New(DEFAULTS, DOT_SEPARATOR)
	t.Eq("3.20", dec.From.MFloat(3.2).Str(dec))
}

func (s *decimal) String_truncates_fractionals_as_set(t *T) {
	dec := Dec.New(DEFAULTS, TWO_FRACTIONALS|DOT_SEPARATOR)
	t.Eq("3.20", dec.From.MFloat(3.2).Str(dec))
	t.Eq("3.00", dec.From.MFloat(3.002).Str(dec))
}

func (s *decimal) String_pads_prefixing_fractional_zeros(t *T) {
	dec := Dec.New(DEFAULTS, TWO_FRACTIONALS|DOT_SEPARATOR)
	t.Eq("3.02", dec.From.MFloat(3.02).Str(dec))
}

func (s *decimal) String_pads_suffixing_fractional_zeros(t *T) {
	dec := Dec.New(DEFAULTS, EIGHT_FRACTIONALS|DOT_SEPARATOR)
	t.Eq("3.20000000", dec.From.MFloat(3.2).Str(dec))
}

func (s *decimal) Rounds_to_even(t *T) {
	dec := Dec.New(DEFAULTS, EIGHT_FRACTIONALS|DOT_SEPARATOR)
	t.Eq("3.19400000", dec.From.MFloat(3.194).Rnd(dec))
	dec.flags.SetFmt(TWO_FRACTIONALS)
	t.Eq("3.19", dec.From.MFloat(3.194).Rnd(dec))
	t.Eq("3.19", dec.From.MFloat(3.193).Rnd(dec))
	t.Eq("3.19", dec.From.MFloat(3.192).Rnd(dec))
	t.Eq("3.19", dec.From.MFloat(3.191).Rnd(dec))
	t.Eq("3.19", dec.From.MFloat(3.190).Rnd(dec))
	t.Eq("3.20", dec.From.MFloat(3.196).Rnd(dec))
	t.Eq("3.20", dec.From.MFloat(3.197).Rnd(dec))
	t.Eq("3.20", dec.From.MFloat(3.198).Rnd(dec))
	t.Eq("3.20", dec.From.MFloat(3.199).Rnd(dec))
	t.Eq("3.20", dec.From.MFloat(3.195).Rnd(dec))
	t.Eq("3.20", dec.From.MFloat(3.1950).Rnd(dec))
	t.Eq("3.20", dec.From.MFloat(3.205).Rnd(dec))
	t.Eq("3.20", dec.From.MFloat(3.205000).Rnd(dec))
	t.Eq("3.21", dec.From.MFloat(3.2051).Rnd(dec))
	t.Eq("3.21", dec.From.MFloat(3.205001).Rnd(dec))
	t.Eq("3.03", dec.From.MFloat(3.025001).Rnd(dec))
	t.Eq("3.00", dec.From.MFloat(3.005000).Rnd(dec))
	t.Eq("4.00", dec.From.MFloat(3.999999).Rnd(dec))
	t.Eq("10.00", dec.From.MFloat(9.999999).Rnd(dec))
	t.Eq("0.00", dec.From.MFloat(0.005).Rnd(dec))
}

func (s *decimal) String_matches_converted_literal(t *T) {
	dec := Dec.New(DEFAULTS, SIX_FRACTIONALS|DOT_SEPARATOR)
	for _, s := range []string{
		"0.123456", "0.999999", "0.000001", "0.005000",
	} {
		t.Eq(s, dec.From.MStr(s).Str(dec))
	}
	t.Eq("0.123456", dec.From.MFloat(0.123456).Str(dec))
	t.Eq("0.999999", dec.From.MFloat(0.999999).Str(dec))
	t.Eq("0.000001", dec.From.MFloat(0.000001).Str(dec))
	t.Eq("0.005000", dec.From.MFloat(0.005000).Str(dec))
}

func TestDecimal(t *testing.T) {
	t.Parallel()
	Run(&decimal{}, t)
}
