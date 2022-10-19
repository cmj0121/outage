package status

import (
	"fmt"
)

func ExampleService() {
	service := Service{
		Meta: Meta{
			Title:   "example",
			Subject: "The example service",
			Link:    "https://example.com",
		},
	}

	fmt.Println(service)
	// Output:
	// title: example
	// subject: The example service
	// link: https://example.com
}
