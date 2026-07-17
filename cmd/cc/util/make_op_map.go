package util

func MakeOpmap[I ~int, O ~int](m map[I]O) func(I) O {
	return func(input I) O {
		if op, ok := m[input]; ok {
			return op
		} else {
			Exit_with_printf("failed to map %v (of type %T)\n",
				input, input)
			return op
		}
	}
}
