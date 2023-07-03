package slices_test

import (
	"strings"
	"testing"

	"github.com/monadicstack/abide/internal/slices"
	"github.com/stretchr/testify/suite"
)

func TestSlicesSuite(t *testing.T) {
	suite.Run(t, new(SlicesSuite))
}

type SlicesSuite struct {
	suite.Suite
}

func (suite *SlicesSuite) TestMap() {
	r := suite.Require()

	var nilSlice []string
	var nilFunc func(string) string

	r.Nil(slices.Map(nilSlice, func(value string) string { return value }))
	r.Nil(slices.Map([]string{"a", "b"}, nilFunc))

	// Identity function behaves as expected
	actualStrings := slices.Map([]string{"A", "B", "C"}, func(val string) string {
		return val
	})
	r.Equal([]string{"A", "B", "C"}, actualStrings)

	// Can transform values to another format of the same type
	actualStrings = slices.Map([]string{"A", "B", "C"}, strings.ToLower)
	r.Equal([]string{"a", "b", "c"}, actualStrings)

	// Can transform input slice to a slice of a completely different type.
	actualInts := slices.Map([]string{"Hello", "Goodbye", "Foo"}, func(text string) int {
		return len(text)
	})
	r.Equal([]int{5, 7, 3}, actualInts)
}
