// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

import (
	"strconv"
	"strings"
)

const wordLength = 32 << (^uint(0) >> 63)

// Set provides a fast implementation for integer sets of small
// non-negative integers.
type Set struct {
	words       []uint
	cardinality int
}

// Len returns the set's cardinality.
func (s *Set) Len() int { return s.cardinality }

// IsEmpty returns true if the set's cardinality is zero.
func (s *Set) IsEmpty() bool { return s.cardinality == 0 }

// FromSlice constructs a int-set from a given slice.
func FromSlice(elms []int) *Set {
	return (&Set{}).Add(elms...)
}

// ToSlice converts the (ordered) integers of receiving set to a slice.
func (s *Set) ToSlice() (elms []int) {

	appendElm := func(elm int) { elms = append(elms, elm) }
	s.For(appendElm)

	return
}

// Eq returns true if receiving set has the same elements as given other
// set.
func (s *Set) Eq(other *Set) bool {
	return s.Len() == other.Len() && s.HasSub(other)
}

// HasSub returns true if receiving set has other given set as subset.
func (s *Set) HasSub(other *Set) (is bool) {

	if len(other.words) == 0 {
		return true
	}
	if s.Len() < other.Len() {
		return false
	}

	defer func() {
		if recover() != nil {
			is = false
		}
	}()

	is = true
	other.For(func(elm int) {
		if !s.has(elm) {
			panic("no subset")
		}
	})

	return is
}

// Has returns true if given integers are in receiving set; false
// otherwise
func (s *Set) Has(elm int, elms ...int) bool {
	if len(elms) == 0 {
		return s.has(elm)
	}
	if !s.has(elm) {
		return false
	}
	for _, elm := range elms {
		if !s.has(elm) {
			return false
		}
	}
	return true
}

func (s *Set) has(elm int) bool {
	word, bit := elm/wordLength, uint(elm%wordLength)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *Set) add(elm int) {
	if elm < 0 || s.has(elm) {
		return
	}
	s.cardinality++
	word, bit := elm/wordLength, uint(elm%wordLength)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// Add adds given integers to receiving set.
func (s *Set) Add(elms ...int) *Set {
	for _, elm := range elms {
		s.add(elm)
	}
	return s
}

// For calls back for each element e providing e.
func (s *Set) For(elm func(int)) {
	for idx, word := range s.words {
		if word == 0 {
			continue
		}
		for bit := 0; bit <= wordLength; bit++ {
			if word&(1<<bit) != 0 {
				elm(bit + idx*wordLength)
			}
		}
	}
}

func (s *Set) del(elm int) {
	if !s.has(elm) {
		return
	}
	word, bit := elm/wordLength, uint(elm%wordLength)
	s.words[word] &^= 1 << bit
	s.cardinality--
}

// Del removes given elements from receiving set.
func (s *Set) Del(elm int, elms ...int) *Set {
	if len(elms) == 0 {
		s.del(elm)
		return s
	}

	for _, elm := range elms {
		s.del(elm)
	}

	return s
}

// String returns a set's string representation {e1, e2, e3, ..., eN} with eI
// in |N.
func (s *Set) String() string {
	var elms []string
	s.For(func(elm int) { elms = append(elms, strconv.Itoa(elm)) })
	return "{" + strings.Join(elms, ", ") + "}"
}
