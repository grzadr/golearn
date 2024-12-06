package parsing

import (
	"fmt"
	"iter"
	"testing"
)

func registerErrors(errs []error, message string, t *testing.T) {
	if len(errs) == 0 {
		return
	}

	t.Error(message)
	for _, err := range errs {
		t.Error(err)
	}
}

// Pair represents a tuple of two values
type NextPair[T1, T2 any] struct {
	First  T1
	Second T2
}

func IterZip[T1, T2 any](s1 []T1, s2 []T2) iter.Seq2[int, NextPair[T1, T2]] {
	return func(yield func(int, NextPair[T1, T2]) bool) {
		n := len(s1)
		if n2 := len(s2); n2 < n {
			n = n2
		}

		for i := 0; i < n; i++ {
			if !yield(i, NextPair[T1, T2]{s1[i], s2[i]}) {
				return
			}
		}
	}
}

func helperCompareNextPair[T1, T2 comparable](s1 []T1, s2 []T2, expected []NextPair[T1, T2]) []error {
	expected_len := len(expected)

	errs := make([]error, 0, expected_len+1)

	last_visited := -1

	for i, pair := range IterZip(s1, s2) {
		last_visited = i
		if i >= expected_len {
			errs = append(errs, fmt.Errorf("Index %d exceeds expected length %d", i, expected_len))
			return errs
		}
		if ref := expected[i]; pair != ref {
			errs = append(errs, fmt.Errorf("Expected %v, but got %v", ref, pair))
		}
	}

	if last_visited < expected_len-1 {
		errs = append(errs, fmt.Errorf("IterZip visited %d instead of %d", last_visited+1, expected_len))
	}

	return errs
}

func TestIterZip(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []string{"A", "B", "C"}

	expected := []NextPair[int, string]{
		{1, "A"},
		{2, "B"},
		{3, "C"},
	}

	registerErrors(helperCompareNextPair(s1, s2, expected), "Detected errors in IterZip", t)
}

func TestIterZipDiffLength(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	s2 := []string{"A", "B", "C"}

	expected := []NextPair[int, string]{
		{1, "A"},
		{2, "B"},
		{3, "C"},
	}

	registerErrors(helperCompareNextPair(s1, s2, expected), "Detected errors in IterZip with different lengths", t)
}
