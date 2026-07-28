package checkn

import "fmt"

type Store map[string]Symbol
type Symbol struct {
	name             string
	fromCurrentBlock bool
}

var make_temporary = namer()

func namer() func(string) string {
	counter := 0
	return func(var_name string) string {
		name := fmt.Sprintf("%s.%d", var_name, counter)
		counter++
		return name
	}
}
