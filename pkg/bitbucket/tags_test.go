package bitbucket

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestTagCreate(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	payload := TagCreatePayload{
		Name: "release-2.0.0",
	}

	// check 201 response is accepted
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
		201, "bitbucket/tag-create.json",
	)
	tag, err := bitbucketClient.TagCreate("PRJ", "my-repo", payload)
	if err != nil {
		t.Fatal(err)
	}
	if tag.DisplayID != "release-2.0.0" {
		t.Fatalf("got %s, want %s", tag.DisplayID, "release-2.0.0")
	}

	// check 200 response is accepted as well as Bitbucket seems to return 200 actually
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
		200, "bitbucket/tag-create.json",
	)
	tag, err = bitbucketClient.TagCreate("PRJ", "my-repo", payload)
	if err != nil {
		t.Fatal(err)
	}
	if tag.DisplayID != "release-2.0.0" {
		t.Fatalf("got %s, want %s", tag.DisplayID, "release-2.0.0")
	}

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.TagCreate("PRJ", "my-repo", payload)
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.TagCreate("PRJ", "my-repo", payload)
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}

func TestTagList(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/tags",
		200, "bitbucket/tag-list.json",
	)
	l, err := bitbucketClient.TagList("myproject", "my-repo", TagListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if l.Size != 1 {
		t.Fatalf("got %d, want %d", l.Size, 1)
	}

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/tags",
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.TagList("myproject", "my-repo", TagListParams{})
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/tags",
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.TagList("myproject", "my-repo", TagListParams{})
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}

func TestTagGet(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/tags/release-2.0.0",
		200, "bitbucket/tag-get.json",
	)
	tag, err := bitbucketClient.TagGet("myproject", "my-repo", "release-2.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if tag.Type != "TAG" {
		t.Fatalf("got %s, want %s", tag.Type, "TAG")
	}

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/tags/release-2.0.0",
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.TagGet("myproject", "my-repo", "release-2.0.0")
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/tags/release-2.0.0",
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.TagGet("myproject", "my-repo", "release-2.0.0")
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}
