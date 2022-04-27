package bitbucket

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/testfile"
	"github.com/opendevstack/pipeline/test/testserver"
)

func TestInsightReportCreate(t *testing.T) {
	sha := "56625c80087b034847001d22502063adae9759f2"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	payload := InsightReportCreatePayload{
		Data: []InsightReportData{
			{
				Title: "Some title",
				Value: "Some value",
				Type:  "TEXT",
			},
			{
				Title: "Build length",
				Value: 60000,
				Type:  "DURATION",
			},
			{
				Title: "Download link",
				Value: map[string]string{
					"linktext": "installer.zip",
					"href":     "https://link.to.download/file.zip",
				},
				Type: "LINK",
			},
			{
				Title: "Build started date",
				Value: 1539656375,
				Type:  "DATE",
			},
			{
				Title: "Code coverage",
				Value: 85,
				Type:  "PERCENTAGE",
			},
			{
				Title: "Some count",
				Value: 5,
				Type:  "NUMBER",
			},
		},
		Details:     "This is the details of the report, it can be a longer string describing the report",
		Title:       "report.title",
		Reporter:    "Reporter/tool that produced this report",
		CreatedDate: 1621231657051,
		Link:        "http://insight.host.com",
		LogoURL:     "http://insight.host.com/logo",
		Result:      "PASS",
	}

	srv.EnqueueResponse(
		t, "/rest/insights/1.0/projects/PRJ/repos/my-repo/commits/"+sha+"/reports/report.key",
		200, "bitbucket/insight-report-create.json",
	)
	r, err := bitbucketClient.InsightReportCreate("PRJ", "my-repo", sha, "report.key", payload)
	if err != nil {
		t.Fatal(err)
	}
	if r.Key != "report.key" {
		t.Fatalf("got %s, want %s", r.Key, "report.key")
	}
	checkLastRequest(t, srv, "bitbucket/insight-report-create.json")

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/insights/1.0/projects/PRJ/repos/my-repo/commits/"+sha+"/reports/report.key",
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.InsightReportCreate("PRJ", "my-repo", sha, "report.key", payload)
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/insights/1.0/projects/PRJ/repos/my-repo/commits/"+sha+"/reports/report.key",
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.InsightReportCreate("PRJ", "my-repo", sha, "report.key", payload)
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}

func checkLastRequest(t *testing.T, srv *testserver.TestServer, golden string) {
	req, err := srv.LastRequest()
	if err != nil {
		t.Fatal(err)
	}
	gotBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	var got interface{}
	err = json.Unmarshal(gotBody, &got)
	if err != nil {
		t.Fatal(err)
	}
	wantBody := testfile.ReadGolden(t, golden)
	var want interface{}
	err = json.Unmarshal(wantBody, &want)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("request body mismatch (-want +got):\n%s", diff)
	}
}
