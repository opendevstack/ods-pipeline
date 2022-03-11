package main

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDirectoryCleaningSparesCache(t *testing.T) {

	tests := map[string]struct {
		fileSystem       fstest.MapFS
		expectedRemovals []string
	}{
		"testCacheSparedEmpty": { // likely to never happen
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
				".ODS-Cache/deps/npm/hithere_1.0/package.json": {},
				"src/app.js":                                   {},
				"package.json":                                 {},
				".env":                                         {},
			},
			[]string{
				".ODS-Cache", // no plan to support case insensitive checking
				".env",
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

	timeNow := time.Now()

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
				".ods-cache/deps/dep1.txt",
			},
		},
		"testCacheCleanNotRemovingOutsideFiles": {
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
				".ods-cache/deps/dep1.txt",
			},
		},
		"testCacheCleanRemovesPreviousBuildTasksWithoutTimestamp": {
			fstest.MapFS{
				".ods-cache/deps/go/gd2.txt":        {},
				".ods-cache/build-task/sha-0/a.txt": {},
				".ods-cache/build-task/sha-0/b.txt": {},
				".ods-cache/build-task/sha-1/c.txt": {},
				".ods-cache/build-task/sha-1/.ods-last-used-stamp": {
					ModTime: timeNow,
				},
				".ods-cache/build-task/go/sha-0/a.txt": {},
				".ods-cache/build-task/go/sha-0/b.txt": {},
				".ods-cache/build-task/go/sha-1/c.txt": {},
				".ods-cache/build-task/go/sha-1/.ods-last-used-stamp": {
					ModTime: timeNow,
				},
			},
			[]string{
				".ods-cache/build-task/go/sha-0",
				".ods-cache/build-task/sha-0/a.txt",
				".ods-cache/build-task/sha-0/b.txt",
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

func TestCleanNotRecentlyUsed(t *testing.T) {

	timeNow := time.Now()
	timeMonthAgo := timeNow.AddDate(0, -1, 0)
	time8DaysAgo := timeNow.AddDate(0, 0, -8)
	time7DaysAgo := timeNow.AddDate(0, 0, -7)
	time6DaysAgo := timeNow.AddDate(0, 0, -6)
	tests := map[string]struct {
		fileSystem       fstest.MapFS
		parentDir        string
		keepTimestamp    time.Time
		expectedCount    int
		expectedRemovals []string
	}{
		"testCleanNRUWhenEmptyFS": {
			fstest.MapFS{},
			"built-task",
			timeNow,
			0,
			[]string{},
		},
		"testCleanNRUWhenNotHavingMarkers": {
			fstest.MapFS{
				".a":                                   {},
				"deps/dep1.txt":                        {},
				"deps/go/gd2.txt":                      {},
				"build-task/go-arch0/sha-0/a.txt":      {},
				"build-task/go-arch0/sha-0/b.txt":      {},
				"build-task/python/sha-1/c.txt":        {},
				"build-task/unexpected/unexpected.txt": {},
				"other/something/c.txt":                {},
			},
			"build-task",
			timeNow,
			6,
			[]string{
				"build-task/go-arch0/sha-0",
				"build-task/python/sha-1",
				"build-task/unexpected/unexpected.txt",
			},
		},
		"testCleanNRUWhenWithMarkers": {
			fstest.MapFS{
				"build-task/go-arch0/sha-0/.ods-last-used-stamp": {
					ModTime: time.Time(time6DaysAgo),
				},
				"build-task/go-arch0/sha-0/a.txt": {},
				"build-task/go-arch0/sha-0/b.txt": {},
				"build-task/python/sha-1/.ods-last-used-stamp": {
					ModTime: time.Time(time8DaysAgo),
				},
				"build-task/python/sha-1/a.txt": {},
				"build-task/python/sha-1/b.txt": {},
				"build-task/python/sha-2/.ods-last-used-stamp": {
					ModTime: time.Time(timeMonthAgo),
				},
				"build-task/python/sha-2/c.txt": {},
			},
			"build-task",
			time7DaysAgo,
			5,
			[]string{
				"build-task/python/sha-1",
				"build-task/python/sha-2",
			},
		},
		"testCleanNRUWhenWithUpperlevelMarkers": {
			fstest.MapFS{
				"build-task/go-arch0/.ods-last-used-stamp": {
					ModTime: time.Time(time8DaysAgo),
				},
				"build-task/go-arch0/unexpected/unexpected.txt": {},
				"build-task/go-arch0/sha-0/a.txt":               {},
				"build-task/go-arch0/sha-0/b.txt":               {},
				"build-task/python/sha-2/.ods-last-used-stamp": {
					ModTime: time.Time(timeMonthAgo),
				},
				"build-task/python/sha-2/c.txt": {},
			},
			"build-task",
			time7DaysAgo,
			3,
			[]string{
				"build-task/go-arch0",
				"build-task/python/sha-2",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			removed := []string{}
			// fsSub, err := tc.fileSystem.Sub(".")
			// if err != nil {
			// 	t.Fatal(err)
			// }
			count, err := cleanupNotRecentlyUsed(
				tc.fileSystem,
				tc.parentDir,
				tc.keepTimestamp,
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
			if count != tc.expectedCount {
				t.Fatalf("expected count %d got %d)", tc.expectedCount, count)
			}
		})
	}
}
