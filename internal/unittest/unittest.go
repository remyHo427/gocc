package unittest

import (
	"testing"

	"cc260717/internal/assert"
)

func MakeTester[T comparable](t *testing.T, f func() T, seq []T) func() {
	return func() {
		t.Helper()

		i := 0
		for _, expect := range seq {
			assert.Equal(t, f(), expect)
			i++
		}

		if ttlen := len(seq); ttlen < i {
			t.Errorf("expect %d value(s), only got %d\n", ttlen, i)
		}
	}
}

type NamedTestPair struct {
	Name   string
	Src    string
	Expect string
}

func MakePairTester(t *testing.T, f func(string) string, tt []NamedTestPair) func() {
	return func() {
		t.Helper()

		for _, test := range tt {
			t.Run(test.Name, func(t *testing.T) {
				assert.Equal(t, f(test.Src), test.Expect)
			})
		}
	}
}

type NamedTestObj struct {
	Name   string
	Src    string
	Expect []string
}

func MakeObjTester(t *testing.T, f func(string) []string, tt []NamedTestObj) func() {
	return func() {
		t.Helper()

		for _, test := range tt {
			t.Run(test.Name, func(t *testing.T) {
				actuals := f(test.Src)

				lactual := len(actuals)
				lexpect := len(test.Expect)
				if lactual != lexpect {
					t.Errorf("expect %d items, got %d\n", lactual, lexpect)
				}

				for i, actual := range actuals {
					assert.Equal(t, actual, test.Expect[i])
				}
			})
		}
	}
}
