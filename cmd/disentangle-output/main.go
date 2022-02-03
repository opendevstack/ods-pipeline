package main

import (
	"flag"
	"log"
	"os"

	"github.com/opendevstack/pipeline/internal/testoutput"
)

type options struct {
	testResultDir string
}

func main() {
	opts := options{}
	flag.StringVar(&opts.testResultDir, "test-result-dir", os.Getenv("TEST_RESULT_DIR"), "test-result-dir")
	flag.Parse()

	os.MkdirAll(opts.testResultDir, 0777)
	err := testoutput.DisentangleTestOutputs(opts.testResultDir, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("could not disentangle test output: %v", err)
	}
}
