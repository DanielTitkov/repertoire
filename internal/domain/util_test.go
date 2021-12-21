package domain

import (
	"reflect"
	"testing"
)

func TestReverstFloat64(t *testing.T) {
	s := []Term{
		{Title: "Dad"},
		{Title: "Goose"},
		{Title: "Mom"},
		{Title: "Cat"},
	}
	exp := []Term{
		{Title: "Cat"},
		{Title: "Mom"},
		{Title: "Goose"},
		{Title: "Dad"},
	}

	res := ReverseSliceTerms(s)
	if !reflect.DeepEqual(res, exp) {
		t.Errorf("expected to get %v, got %v", exp, res)
	}
}
