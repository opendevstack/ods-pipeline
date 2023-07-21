package main

import (
	"fmt"
	"time"

	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
)

// createBitbucketInsightReport attaches a code insight report to the Git commit
// being built in Bitbucket. The code insight report points to the Aqua security scan.
func createBitbucketInsightReport(opts options, aquaScanUrl string, success bool, ctxt *pipelinectxt.ODSContext) error {
	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	}
	bitbucketClient, err := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: opts.bitbucketAccessToken,
		BaseURL:  opts.bitbucketURL,
		Logger:   logger,
	})
	if err != nil {
		return fmt.Errorf("bitbucket client: %w", err)
	}
	reportKey := "org.opendevstack.aquasec"
	scanResult := bitbucket.InsightReportFail
	if success {
		scanResult = bitbucket.InsightReportPass
	}
	_, err = bitbucketClient.InsightReportCreate(
		ctxt.Project,
		ctxt.Repository,
		ctxt.GitCommitSHA,
		reportKey,
		bitbucket.InsightReportCreatePayload{
			Title:       "Aqua Security",
			Reporter:    "OpenDevStack",
			CreatedDate: time.Now().Unix(),
			Details:     "Please visit the following link to review the Aqua Security scan report:",
			Result:      scanResult,
			Data: []bitbucket.InsightReportData{
				{
					Title: "Report",
					Type:  "LINK",
					Value: map[string]string{
						"linktext": "Result in Aqua",
						"href":     aquaScanUrl,
					},
				},
			},
		},
	)
	return err
}
