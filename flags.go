// Copyright 2022 Stephan.Lukits@gmail.com. All rights reserved.
// Use of this source code is governed by the MIT license.

package ints

// Flags constants control a decimal contexts Decimal-conversion, its
// arithmetics and a [Decimal]-values string formatting.
type Flags uint64

func (ff *Flags) set(flags Flags) {
	for f := range ffFractionalsSet {
		if flags&f == f {
			*ff &^= ffFractionals
			*ff |= f
		}
	}
	for s := range ffSeparatorsSet {
		if flags&s == s {
			*ff &^= ffSeparators
			*ff |= s
		}
	}
}

// ..._SEPARATOR defines a context's Separator-property for Decimal
// conversion or string output;
// ..._FRACTIONALS defines a context's number of fractional positions
// for Decimal conversion and decimal arithmetics or string output.
const (
	COMMA_SEPARATOR Flags = 1 << iota
	DOT_SEPARATOR
	ONE_FRACTIONAL
	TWO_FRACTIONALS
	THREE_FRACTIONALS
	FOUR_FRACTIONALS
	FIVE_FRACTIONALS
	SIX_FRACTIONALS
	SEVEN_FRACTIONALS
	EIGHT_FRACTIONALS

	// DEFAULTS used at Context.New allows to indicate that the format
	// flags or arithmetic flags are copied from the used Context
	// instance.
	DEFAULTS = 0
)

const ffFractionals = ONE_FRACTIONAL | TWO_FRACTIONALS |
	THREE_FRACTIONALS | FOUR_FRACTIONALS | FIVE_FRACTIONALS |
	SIX_FRACTIONALS | SEVEN_FRACTIONALS | EIGHT_FRACTIONALS

var ffFractionalsSet = map[Flags]bool{
	ONE_FRACTIONAL:    true,
	TWO_FRACTIONALS:   true,
	THREE_FRACTIONALS: true,
	FOUR_FRACTIONALS:  true,
	FIVE_FRACTIONALS:  true,
	SIX_FRACTIONALS:   true,
	SEVEN_FRACTIONALS: true,
	EIGHT_FRACTIONALS: true,
}

var flagsToFractionals = map[Flags]int{
	0:                 0,
	ONE_FRACTIONAL:    1,
	TWO_FRACTIONALS:   2,
	THREE_FRACTIONALS: 3,
	FOUR_FRACTIONALS:  4,
	FIVE_FRACTIONALS:  5,
	SIX_FRACTIONALS:   6,
	SEVEN_FRACTIONALS: 7,
	EIGHT_FRACTIONALS: 8,
}

const ffSeparators = COMMA_SEPARATOR | DOT_SEPARATOR

var ffSeparatorsSet = map[Flags]bool{
	COMMA_SEPARATOR: true,
	DOT_SEPARATOR:   true,
}

var flagsToSeparator = map[Flags]rune{
	0:               0,
	COMMA_SEPARATOR: ',',
	DOT_SEPARATOR:   '.',
}

type flags struct {
	cntx     *Context
	art, fmt Flags
}

func newFlags(c *Context) *flags {
	ff := flags{cntx: c, art: DOT_SEPARATOR | SIX_FRACTIONALS}
	ff.fmt = COMMA_SEPARATOR | TWO_FRACTIONALS
	return &ff
}

func (ff *flags) copy(ia, fmt Flags) *flags {
	_ff := flags{art: ff.art, fmt: ff.fmt}
	if ia != DEFAULTS {
		_ff.art.set(ia)
	}
	if fmt != DEFAULTS {
		_ff.fmt.set(fmt)
	}
	return &_ff
}

// SetFmt sets the flags controlling the string formatting of a
// Decimal-value.  Note SetFmt(Default) is an no-op and if more than one
// separator constant or more than one fractional positions constant is
// given only one of them will be chosen and it is undefined which one.
func (ff *flags) SetFmt(flags Flags) *Context {
	ff.fmt.set(flags)
	return ff.cntx
}

func (ff *flags) Fractionals() int {
	return flagsToFractionals[ff.art&ffFractionals]
}

func (ff *flags) FmtFractionals() int {
	return flagsToFractionals[ff.fmt&ffFractionals]
}
