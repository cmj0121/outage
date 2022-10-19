package status

import (
	"fmt"
)

func ExampleMode() {
	var mode Mode = -1

	fmt.Println(ON)
	fmt.Println(OFF)
	fmt.Println(INCIDENT)
	fmt.Println(UNKNOWN)
	fmt.Println(mode)
	// Output:
	// on
	// off
	// incident
	// unknown
	// unknown
}
