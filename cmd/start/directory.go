// Selective deletion in a directory tree.
//
// The implementation is a wrapper around [io.fs.WalkDir] which
// operates on the filesystem abstraction [io.fs.FS].
// This design facilitates unit testing.
//
// However the abstraction is also limited so that it does not
// directly allow operation such as removing files.
//
// This implementation supports removal of files while selectively
// walking the directory.
//
package main

import (
	"fmt"
	"io/fs"
	"os"
)

type WalkAndRemovalFlags int

const (
	// Used for directories only, where it indicates that the directory is not touched and the children are skipped.
	skipDir WalkAndRemovalFlags = 1 << iota
	// If Flag is set for a file the the associated file is deleted.
	// If Flag is set for a directory, the directory is removed recursively and further walking is skipped.
	remove = 1 << iota
	// Used for directories only, where it indicates to continue with the directories children.
	enterDir = 1 << iota
)

type WalkAndRemovalDeciderFunc func(path string, d fs.DirEntry) WalkAndRemovalFlags

// Type RemoveFunc allows to replace the error handling, for example to add logging.
type RemoveFunc func(path string, isDir bool) error

func removeFileOrDir(path string, isDir bool) error {
	err := os.RemoveAll(path)
	if err != nil {
		if isDir {
			return fmt.Errorf("could not remove directory %s: %w", path, err)
		} else {
			return fmt.Errorf("could not remove file %s: %w", path, err)
		}
	}
	return nil
}

// Selective deletions in a directory tree.
func deleteDirRecursiveWithSkip(
	root fs.FS,
	fnWalkAndRemoveDecider WalkAndRemovalDeciderFunc,
	fnRemove RemoveFunc,
) error {
	return fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("could not read files in %s: %w", path, err)
		}
		if d.IsDir() {
			if path == "." {
				return nil
			}
			flags := fnWalkAndRemoveDecider(path, d)
			if flags&skipDir != 0 {
				return fs.SkipDir
			}
			if flags&remove != 0 {
				if err := fnRemove(path, true); err != nil {
					return err
				}
			}
			if flags&enterDir != 0 {
				return nil
			}
			return fs.SkipDir
		}
		if d.Type().IsRegular() {
			flags := fnWalkAndRemoveDecider(path, d)
			if flags&remove != 0 {
				if err := fnRemove(path, true); err != nil {
					return err
				}
			}
			return nil
		}
		return nil
	})
}
