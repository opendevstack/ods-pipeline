package main

import (
	"flag"
	"os"
	"time"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	stateFlag := flag.String("state", "", "state")
	keyFlag := flag.String("key", "", "key")
	nameFlag := flag.String("name", "", "name")
	urlFlag := flag.String("url", "", "url")
	descriptionFlag := flag.String("description", "", "description")
	commitFlag := flag.String("commit", "", "commit")
	apiTokenFlag := flag.String("api-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "api-token")
	bitbucketURLFlag := flag.String("bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	flag.Parse()

	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		Timeout:    10 * time.Second,
		APIToken:   *apiTokenFlag,
		MaxRetries: 2,
		BaseURL:    *bitbucketURLFlag,
	})
	bitbucketClient.BuildStatusPost(*commitFlag, bitbucket.BuildStatusPostPayload{
		State:       *stateFlag,
		Key:         *keyFlag,
		Name:        *nameFlag,
		URL:         *urlFlag,
		Description: *descriptionFlag,
	})
}
