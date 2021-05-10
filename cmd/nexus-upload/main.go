package main

import (
	"flag"

	"github.com/opendevstack/pipeline/pkg/nexus"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	nexusURL := flag.String("nexus-url", "", "URL of Nexus instance")
	nexusUser := flag.String("nexus-user", "", "User of Nexus instance")
	nexusPassword := flag.String("nexus-password", "", "Password of Nexus user")
	repository := flag.String("repository", "", "Nexus repository")
	group := flag.String("group", "", "Repository group")
	file := flag.String("file", "", "Filename to upload (absolute")

	flag.Parse()

	err := nexus.Upload(
		*nexusURL,
		*nexusUser,
		*nexusPassword,
		*repository,
		*group,
		*file,
	)
	check(err)
}
