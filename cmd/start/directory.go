// Selective deletion in a directory tree.
//
// The implementation is a wrapper around [io.fs.WalkDir] which
// operates on the filesystem abstraction [io.fs.FS].
// The reason this is used as they facilitate unit testing.
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

type SkipFileRemovalFlags int

const (
	// Used for directories only, where it indicates that the directory is not touched and the children are skipped.
	leaveAlone SkipFileRemovalFlags = 1 << iota
	// If Flag is set for a file the the associated file is deleted.
	// If Flag is set for a directory, the directory is removed recursively and further walking is skipped.
	remove = 1 << iota
	// Used for directories only, where it indicates to continue with the directories children.
	walkChildren = 1 << iota
)

type FileSkipFunc func(path string, d fs.DirEntry) SkipFileRemovalFlags

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
func deleteDirRecursiveWithSkip(root fs.FS, fnSkip FileSkipFunc, fnRemove RemoveFunc) error {
	return fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("could not read files in %s: %w", path, err)
		}
		if d.IsDir() {
			if path == "." {
				return nil
			}
			skipFlags := fnSkip(path, d)
			if skipFlags&leaveAlone != 0 {
				return fs.SkipDir
			}
			if skipFlags&remove != 0 {
				err = fnRemove(path, true)
			}
			if err != nil {
				return err
			}
			if skipFlags&walkChildren != 0 {
				return nil
			}
			return fs.SkipDir

		}
		if d.Type().IsRegular() {
			skipFlags := fnSkip(path, d)
			if skipFlags&remove == 0 {
				return nil
			}
			err = fnRemove(path, false)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
