package e2e

import (
	"flag"
)

var privateCertFlag = flag.Bool("private-cert", false, "Whether to run tests using a private cert")
