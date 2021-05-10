package main

import (
	"flag"
	"fmt"

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
	flag.Parse()

	assets, err := nexus.Download(
		*nexusURL,
		*nexusUser,
		*nexusPassword,
		*repository,
		*group,
	)
	check(err)

	for _, a := range assets {
		fmt.Printf("%s ", a)
	}
}
