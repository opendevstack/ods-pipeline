package main

import (
	"flag"
	"github.com/opendevstack/pipeline/internal/testoutput"
	"os"
)

type options struct {
	testResultDir string
}

func main() {
	opts := options{}
	flag.StringVar(&opts.testResultDir, "test-result-dir", os.Getenv("TEST_RESULT_DIR"), "test-result-dir")
	flag.Parse()

	os.MkdirAll(opts.testResultDir, 0777)
	testoutput.DisentangleTestOutputs(opts.testResultDir, os.Stdin, os.Stdout)
}
