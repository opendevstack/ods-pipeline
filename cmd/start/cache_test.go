package main

import (
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
)

func TestDirectoryCleaningSparesCache(t *testing.T) {

	tests := map[string]struct {
		fileSystem       fstest.MapFS
		expectedRemovals []string
	}{
		"testCacheASparedEmpty": { // likely to never happen
			fstest.MapFS{},
			[]string{},
		},
		"testCacheSpared": {
			fstest.MapFS{
				".ods-cache/.a":                                {},
				".ods-cache/deps/dep1.txt":                     {},
				".ods-cache/deps/go/gd1.txt":                   {},
				".ods-cache/deps/go/gd1/foo.txt":               {},
				".ods-cache/deps/go/gd2.txt":                   {},
				".ods-cache/deps/npm/hithere_1.0/package.json": {},
				"src/app.js":                                   {},
				"package.json":                                 {},
				".env":                                         {},
			},
			[]string{
				".env",
				"package.json",
				"src",
			},
		},
		"testCacheSparedCaseSensitive": {
			fstest.MapFS{
				".ods-cache/.a":                                {},
				".ods-cache/deps/dep1.txt":                     {},
				".ods-cache/deps/go/gd1.txt":                   {},
				".ods-cache/deps/go/gd1/foo.txt":               {},
				".ods-cache/deps/go/gd2.txt":                   {},
				".ods-Cache/deps/npm/hithere_1.0/package.json": {},
				"src/app.js":                                   {},
				"package.json":                                 {},
				".env":                                         {},
			},
			[]string{
				".env",
				".ods-Cache",
				"package.json",
				"src",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			removed := []string{}
			err := deleteDirectoryContentsSpareCache(
				FileSystemBase{tc.fileSystem, "."},
				func(path string, isDir bool) error {
					removed = append(removed, path)
					return nil
				})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.expectedRemovals, removed); diff != "" {
				t.Fatalf("expected (-want +got):\n%s", diff)
			}
		})
	}
}
func TestCacheCleaning(t *testing.T) {

	tests := map[string]struct {
		fileSystem       fstest.MapFS
		expectedRemovals []string
	}{
		"testCacheCleanWhenNoCache": {
			fstest.MapFS{},
			[]string{},
		},
		"testCacheClean": {
			fstest.MapFS{
				".ods-cache/.a":                                {},
				".ods-cache/deps/dep1.txt":                     {},
				".ods-cache/deps/go/gd1.txt":                   {},
				".ods-cache/deps/go/gd1/foo.txt":               {},
				".ods-cache/deps/go/gd2.txt":                   {},
				".ods-cache/deps/npm/hithere_1.0/package.json": {},
			},
			[]string{
				".ods-cache/.a",
				".ods-cache/deps/dep1.txt",
			},
		},
		"testCacheCleanNotRemovingOutside files": {
			fstest.MapFS{
				".ods-cache/.a":                  {},
				".ods-cache/deps/dep1.txt":       {},
				".ods-cache/deps/go/gd1.txt":     {},
				".ods-cache/deps/go/gd1/foo.txt": {},
				".ods-cache/deps/go/gd2.txt":     {},
				"src/app.js":                     {},
				"package.json":                   {},
				".env":                           {},
			},
			[]string{
				".ods-cache/.a",
				".ods-cache/deps/dep1.txt",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			removed := []string{}
			err := cleanCache(
				FileSystemBase{tc.fileSystem, "."},
				func(path string, isDir bool) error {
					removed = append(removed, path)
					return nil
				})
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expectedRemovals, removed); diff != "" {
				t.Fatalf("expected (-want +got):\n%s", diff)
			}
		})
	}
}
