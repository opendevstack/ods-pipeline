// Provides file cleanup supporting caching.
//
// Pipeline runs have a shared workspace bound to a PVC under the name of
// 'source'.
// All ods pipeline tasks have their working directory on this PVC per
//    workingDir: $(workspaces.source.path)
// This PVC persists between builds but may be recreated by the user as needed.
//
// If the PVC persists source code from prior builds will be in the
// working directory and must be deleted before the latest code is checked
// out.
// However to enable using the PVC for caching its contents must be spared from
// wholesale deletion.
//
// In addition to function deleteDirectoryContentsSpareCache a function
// is provided to clean the cache. This is done in two functions to avoid
// making the code too complex as new cache cleaning functionality will
// likely be added in the future, for example to implement build skipping.

package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	odsCacheDirName             = ".ods-cache"
	odsCacheDependenciesDirName = "deps"
)

type FileSystemBase struct {
	filesystem fs.FS
	base       string
}

func withBaseFileRemover(base string, removeFunc RemoveFunc) RemoveFunc {
	return func(path string, isDir bool) error {
		fspath := filepath.Join(base, path)
		return removeFunc(fspath, isDir)
	}
}

func deleteDirectoryContentsSpareCache(fsb FileSystemBase, fnRemove RemoveFunc) error {
	// Open the directory and read all its files.
	cacheDirPath := filepath.Join(".", odsCacheDirName)
	return deleteDirRecursiveWithSkip(fsb.filesystem, func(path string, d fs.DirEntry) SkipFileRemovalFlags {
		if path == cacheDirPath {
			return leaveAlone
		}
		return remove
	}, withBaseFileRemover(fsb.base, fnRemove))
}

// Cleans the cache
// At the moment only a cache for dependencies is supported and
// All other contents is removed to ensure that tasks don't use
// cache areas accidentally. This effectively reserve other
// areas for future use.
// For example if in the future build skipping is supported
// there would likely be an area where the build output is kept
// per git-sha of the working-dir. In this case a suitable cleanup
// might delete such areas after a certain time (see PR #423).
func cleanCache(fsb FileSystemBase, fnRemove RemoveFunc) error {
	_, err := fsb.filesystem.Open(odsCacheDirName)
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	fsCache, err := fs.Sub(fsb.filesystem, odsCacheDirName)
	if err != nil {
		log.Fatal(err)
	}
	cacheDependenciesPath := filepath.Join(".", odsCacheDependenciesDirName)
	// To avoid spare files inside the cache which are not supported delete
	// all other areas of the cache
	return deleteDirRecursiveWithSkip(fsCache, func(path string, d fs.DirEntry) SkipFileRemovalFlags {
		if strings.HasPrefix(path, cacheDependenciesPath) {
			// Dependencies must be inside a folder specific to a technology
			// such as for npm or go.
			// Clean all files which are not directories
			if path == cacheDependenciesPath {
				return walkChildren
			}
			listOfPath := strings.Split(path, string(os.PathSeparator)) // https://stackoverflow.com/a/33619038
			if len(listOfPath) == 2 {
				if d.IsDir() {
					return leaveAlone
				} else {
					return remove
				}
			}
			// this is not expected to be reached
			return leaveAlone
		}
		return remove // delete everything else
	}, withBaseFileRemover(filepath.Join(fsb.base, odsCacheDirName), fnRemove))
}
