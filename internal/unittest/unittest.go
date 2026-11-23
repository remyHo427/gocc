package unittest

import (
	"testing"

	"cc260717/internal/assert"
)

func Test[T comparable](t *testing.T, f func() T, seq []T) {
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

type NamedTestPair struct {
	Name   string
	Src    string
	Expect string
}

func TestPair(t *testing.T, f func(string) string, tt []NamedTestPair) {
	t.Helper()

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, f(test.Src), test.Expect)
		})
	}
}

type NamedTestObj struct {
	Name   string
	Src    string
	Expect []string
}

func TestObj(t *testing.T, f func(string) []string, tt []NamedTestObj) {
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
