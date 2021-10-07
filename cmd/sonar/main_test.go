package main

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
)

type fakeClient struct {
	scanPerformed        bool
	passQualityGate      bool
	qualityGateRetrieved bool
	reportGenerated      bool
}

func (c *fakeClient) Scan(sonarProject, branch, commit string, pr *sonar.PullRequest) (string, error) {
	c.scanPerformed = true
	return "", nil
}

func (c *fakeClient) QualityGateGet(p sonar.QualityGateGetParams) (*sonar.QualityGate, error) {
	c.qualityGateRetrieved = true
	status := sonar.QualityGateStatusError
	if c.passQualityGate {
		status = sonar.QualityGateStatusOk
	}
	return &sonar.QualityGate{ProjectStatus: sonar.QualityGateProjectStatus{Status: status}}, nil
}

func (c *fakeClient) GenerateReports(sonarProject, author, branch, rootPath, artifactPrefix string) (string, error) {
	c.reportGenerated = true
	return "", nil
}

func (c *fakeClient) ExtractComputeEngineTaskID(filename string) (string, error) {
	return "abc123", nil
}

func (c *fakeClient) ComputeEngineTaskGet(params sonar.ComputeEngineTaskGetParams) (*sonar.ComputeEngineTask, error) {
	return &sonar.ComputeEngineTask{Status: sonar.TaskStatusSuccess}, nil
}

func TestSonarScan(t *testing.T) {
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}

	tests := map[string]struct {
		// which SQ edition is in use
		optSonarEdition string
		// whether quality gate is required to pass
		optQualityGate bool

		// PR key
		ctxtPrKey string
		// PR base
		ctxtPrBase string

		// whether the quality gate in SQ passes (faked)
		passQualityGate bool

		// whether scan should have been performed
		wantScanPerformed bool
		// whether report should have been generated
		wantReportGenerated bool
		// whether quality gate should have been retrieved
		wantQualityGateRetrieved bool
		// whether scanning should fail - if not empty, the actual error message
		// will be checked to contain wantErr.
		wantErr string
	}{
		"developer edition generates report when no PR is present": {
			optSonarEdition:          "developer",
			optQualityGate:           true,
			ctxtPrKey:                "",
			ctxtPrBase:               "",
			passQualityGate:          true,
			wantScanPerformed:        true,
			wantReportGenerated:      true,
			wantQualityGateRetrieved: true,
		},
		"developer edition does not generate report when PR is present": {
			optSonarEdition:          "developer",
			optQualityGate:           true,
			ctxtPrKey:                "3",
			ctxtPrBase:               "master",
			passQualityGate:          true,
			wantScanPerformed:        true,
			wantReportGenerated:      false,
			wantQualityGateRetrieved: true,
		},
		"community edition generates report": {
			optSonarEdition:          "community",
			optQualityGate:           true,
			ctxtPrKey:                "",
			ctxtPrBase:               "",
			passQualityGate:          true,
			wantScanPerformed:        true,
			wantReportGenerated:      true,
			wantQualityGateRetrieved: true,
		},
		"does not check quality gate if disabled": {
			optSonarEdition:          "community",
			optQualityGate:           false,
			ctxtPrKey:                "",
			ctxtPrBase:               "",
			passQualityGate:          true,
			wantScanPerformed:        true,
			wantReportGenerated:      true,
			wantQualityGateRetrieved: false,
		},
		"fails if quality gate does not pass": {
			optSonarEdition:          "community",
			optQualityGate:           true,
			ctxtPrKey:                "",
			ctxtPrBase:               "",
			passQualityGate:          false,
			wantScanPerformed:        true,
			wantReportGenerated:      true,
			wantQualityGateRetrieved: true,
			wantErr:                  "quality gate status is 'ERROR', not 'OK'",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := options{
				sonarEdition: tc.optSonarEdition,
				qualityGate:  tc.optQualityGate,
			}
			ctxt := &pipelinectxt.ODSContext{PullRequestKey: tc.ctxtPrKey, PullRequestBase: tc.ctxtPrBase}
			sonarClient := &fakeClient{passQualityGate: tc.passQualityGate}
			err := sonarScan(logger, opts, ctxt, sonarClient)
			if err != nil {
				if tc.wantErr == "" || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("want err to contain: %s, got err: %s", tc.wantErr, err)
				}
			}
			if sonarClient.scanPerformed != tc.wantScanPerformed {
				t.Fatalf("want scan performed to be %v, got %v", tc.wantScanPerformed, sonarClient.scanPerformed)
			}
			if sonarClient.reportGenerated != tc.wantReportGenerated {
				t.Fatalf("want report generated to be %v, got %v", tc.wantReportGenerated, sonarClient.reportGenerated)
			}
			if sonarClient.qualityGateRetrieved != tc.wantQualityGateRetrieved {
				t.Fatalf("want quality gate retrieved to be %v, got %v", tc.wantQualityGateRetrieved, sonarClient.qualityGateRetrieved)
			}
		})
	}

}
