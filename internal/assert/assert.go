package assert

import "testing"

func Equal[T comparable](t *testing.T, actual T, expect T) {
	t.Helper()

	if actual != expect {
		t.Errorf("expect %v (%T), got %v (%T)\n",
			expect, expect, actual, actual)
	}
}
