package status

import (
	"fmt"
)

func ExampleFake() {
	fake := Fake()

	fmt.Println(fake)
	// Output:
	// interval: 1s
	// setting:
	//     sampling: 1h0m0s
	//     timeout: 1m0s
	// services:
	//     - title: example
	//       subject: The example service
	//       link: https://example.com
	//       tags:
	//         - Group 1
	//         - Group 2
	//       mode: "on"
	//     - title: example
	//       subject: The example service
	//       link: https://example.com/off
	//       tags:
	//         - Group 1
	//         - Group 2
	//       mode: "off"
	//     - title: example
	//       subject: The example service
	//       link: https://example.com/incident
	//       tags:
	//         - Group 1
	//       mode: incident
	// summary:
	//     Group 1:
	//         - title: example
	//           subject: The example service
	//           link: https://example.com
	//           tags:
	//             - Group 1
	//             - Group 2
	//           mode: "on"
	//         - title: example
	//           subject: The example service
	//           link: https://example.com/off
	//           tags:
	//             - Group 1
	//             - Group 2
	//           mode: "off"
	//         - title: example
	//           subject: The example service
	//           link: https://example.com/incident
	//           tags:
	//             - Group 1
	//           mode: incident
	//     Group 2:
	//         - title: example
	//           subject: The example service
	//           link: https://example.com
	//           tags:
	//             - Group 1
	//             - Group 2
	//           mode: "on"
	//         - title: example
	//           subject: The example service
	//           link: https://example.com/off
	//           tags:
	//             - Group 1
	//             - Group 2
	//           mode: "off"
}
