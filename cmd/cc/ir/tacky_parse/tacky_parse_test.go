package tackyparse

import "testing"

func TestReturn(t *testing.T) {
	check(t, "null", "return ;")
	check(t, "integer", "return 0;")
	check(t, "identifier", "return tmp.0;")
}
func TestUnary(t *testing.T) {
	check(t, "negate", "unary - 1 tmp.0;")
	check(t, "complement", "unary ~ 1 tmp.0;")
	check(t, "logical not", "unary ! 1 tmp.0;")
}
func TestBinary(t *testing.T) {
	check(t, "add", "binary + 1 2 tmp.0;")
	check(t, "sub", "binary - 1 2 tmp.0;")
	check(t, "mul", "binary * 1 2 tmp.0;")
	check(t, "div", "binary / 1 2 tmp.0;")
	check(t, "rem", "binary % 1 2 tmp.0;")
	check(t, "bitwise OR", "binary | 1 2 tmp.0;")
	check(t, "bitwise AND", "binary & 1 2 tmp.0;")
	check(t, "bitwise XOR", "binary ^ 1 2 tmp.0;")
	check(t, "left shift", "binary << 1 2 tmp.0;")
	check(t, "right shift", "binary >> 1 2 tmp.0;")
	check(t, "less than", "binary < 1 2 tmp.0;")
	check(t, "greater than", "binary > 1 2 tmp.0;")
	check(t, "less than or equal", "binary <= 1 2 tmp.0;")
	check(t, "greater than or equal", "binary >= 1 2 tmp.0;")
	check(t, "not equal", "binary != 1 2 tmp.0;")
	check(t, "equal", "binary == 1 2 tmp.0;")
	check(t, "logical and", "binary && 1 2 tmp.0;")
	check(t, "logical or", "binary || 1 2 tmp.0;")
}
func TestMisc(t *testing.T) {
	check(t, "label", "label lb.0;")
	check(t, "jump", "jump lb.0;")
	check(t, "copy", "copy 0 tmp.0;")
	check(t, "jump_if_zero", "jump_if_zero 0 lb.0;")
	check(t, "jump_if_zero", "jump_if_not_zero 0 lb.0;")
}

func TestParseAll(t *testing.T) {
	t.Helper()

	src := "unary - 2 tmp.0; return tmp.0;"
	p := NewParser(New(src))
	tacky := p.ParseAll()

	if len(tacky) != 2 {
		t.Errorf("expect 2 instructions, got %d", len(tacky))
	}
	if actual := tacky[0].String(); actual != "unary - 2 tmp.0;" {
		t.Errorf("expect \"unary - 2 tmp.0;\", got \"%s\"", actual)
	}
	if actual := tacky[1].String(); actual != "return tmp.0;" {
		t.Errorf("expect \"return tmp.0\", got \"%s\"", actual)
	}
}

func check(t *testing.T, name string, src string) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		p := NewParser(New(src))
		tacky := p.Parse()
		actual := tacky.String()

		if len(src) != len(actual) {
			t.Errorf("output and input does not have the same length\n")
		} else if src != actual {
			t.Errorf("output and input are not identical")
		}
	})
}
