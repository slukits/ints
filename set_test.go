// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ints

import (
	"testing"

	. "github.com/slukits/gounit"
)

type set struct{ Suite }

func (s *set) Initial_length_is_zero(t *T) {
	t.Eq(0, (&Set{}).Len())
}

func (s *set) Is_initially_empty(t *T) {
	t.True((&Set{}).IsEmpty())
}

func (s *set) From_empty_slice_is_empty(t *T) {
	t.True(FromSlice([]int{}).IsEmpty())
}

func (s *set) From_slice_has_elements_of_slice(t *T) {
	ss := []int{1, 2, 3}
	st := FromSlice(ss)
	for _, i := range ss {
		t.True(st.Has(i))
	}
}

func (s *set) From_slice_has_len_of_unique_slice_elements(t *T) {
	t.Eq(2, FromSlice([]int{1, 1, 2}).Len())
}

func (s *set) Adding_increases_len_by_unique_added_elements(t *T) {
	t.Eq(2, (&Set{}).Add(1, 1, 2).Len())
}

func (s *set) Has_added_elements(t *T) {
	add := []int{2, 2, 3}
	ii := (&Set{}).Add(1).Add(add...)
	for _, a := range append([]int{1}, add...) {
		t.True(ii.Has(a))
	}
	t.True(ii.Has(1, add...))
}

func (s *set) Has_not_element_not_in_set(t *T) {
	t.False((&Set{}).Has(1))
	t.False((&Set{}).Has(1, 2))
	t.False((&Set{}).Add(1).Has(1, 2))
}

func (s *set) Provides_all_its_elements(t *T) {
	ii, visited := (&Set{}).Add(22, 42, 1200), map[int]bool{}
	in := func(e int) bool {
		for _, _e := range []int{22, 42, 1200} {
			if e != _e {
				continue
			}
			return true
		}
		return false
	}
	ii.For(func(e int) {
		t.True(in(e))
		t.False(visited[e])
		visited[e] = true
	})
	t.Eq(3, len(visited))
}

func (s *set) Empty_set_is_subset_of_empty_set(t *T) {
	t.True((&Set{}).HasSub(&Set{}))
}

func (s *set) Is_no_subset_of_smaller_set(t *T) {
	t.False((&Set{}).HasSub((&Set{}).Add(1)))
}

func (s *set) Is_subset_of_set_having_all_its_elements(t *T) {
	t.True((&Set{}).Add(1, 2).HasSub((&Set{}).Add(1)))
}

func (s *set) Is_no_subset_of_set_having_not_all_its_elements(t *T) {
	t.False((&Set{}).Add(1, 2).HasSub((&Set{}).Add(5)))
}

func (s *set) Is_not_equal_to_other_set_with_different_len(t *T) {
	t.False((&Set{}).Add(1, 2).Eq((&Set{}).Add(1)))
}

func (s *set) Is_not_equal_to_other_set_which_is_no_subset(t *T) {
	t.False((&Set{}).Add(1, 2).Eq((&Set{}).Add(1, 3)))
}

func (s *set) Is_equal_to_other_set_if_subset_and_lens_equal(t *T) {
	t.True((&Set{}).Add(1, 3).Eq((&Set{}).Add(1, 3)))
}

func (s *set) Noops_if_it_doesnt_has_deleted_element(t *T) {
	st := (&Set{}).Add(5)
	t.True(st.Del(1).Eq(st))
}

func (s *set) Doesnt_have_deleted_elements(t *T) {
	st := (&Set{}).Add(3, 3403, 455505, 22, 42)
	t.False(st.Del(22).Has(22))
	t.False(st.Del(3403, 455505).Has(3403, 455505))
}

func (s *set) Provides_its_elements_as_slice(t *T) {
	st := FromSlice([]int{3, 3403, 455505, 22, 42})
	ee := st.ToSlice()
	t.Eq(st.Len(), len(ee))
	in := func(e int) bool {
		for _, _e := range ee {
			if e != _e {
				continue
			}
			return true
		}
		return false
	}
	st.For(func(e int) { t.True(in(e)) })
}

func (s *set) Provides_its_elements_as_string(t *T) {
	st, exp := FromSlice([]int{3, 22, 42}), "{3, 22, 42}"
	t.Eq(exp, st.String())
}

func TestSet(t *testing.T) {
	Run(&set{}, t)
}
