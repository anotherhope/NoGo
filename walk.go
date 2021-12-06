package nogo

import (
	"errors"
	"io/fs"
	"path/filepath"
)

// WalkFunc can be used in any Walk function to automatically ignore ignored files.
// It is similar to ForWalkDir but with it you can write a WalkFunc for any other (than fs.WalkDir) Walk function:
//
// Example for afero:
//  err = afero.Walk(baseFS, ".", func(path string, info fs.FileInfo, err error) error {
//		if ok, err := n.WalkFunc(afero.NewIOFS(baseFS), ".gitignore", path, info.IsDir(), err); !ok {
//			return err
//		}
//
//		fmt.Println(path, info.Name())
//		return nil
//	})
func (n *NoGo) WalkFunc(fsys fs.FS, ignoreFileName string, path string, isDir bool, err error) (bool, error) {
	if err != nil {
		return false, err
	}

	if path != "." {
		if match, _ := n.MatchWithoutParents(path, isDir); match {
			if isDir {
				return false, fs.SkipDir
			}
			return false, nil
		}
	}

	if isDir {
		// Load a maybe existing ignore file if it is not itself ignored.
		possibleIgnoreFile := filepath.Join(path, ignoreFileName)
		if match, _ := n.MatchWithoutParents(possibleIgnoreFile, false); !match {
			err := n.AddFile(fsys, filepath.Join(path, ignoreFileName))
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return false, err
			}
		}
	}

	return true, nil
}

// ForWalkDir can be used to set all parameters of fs.WalkDir.
// It only calls the passed WalkDirFunc for files and directories
// which are not ignored.
//
// If you need something similar for any other Walk function (e.g. afero.Walk)
// You can use WalkFunc for that.
//
// Example:
//  n := nogo.New(nogo.DotGitRule)
//  err = fs.WalkDir(n.ForWalkDir(walkFS, ".", ".gitignore", func(path string, d fs.DirEntry, err error) error {
//		if err != nil {
//			return err
//		}
//		fmt.Println(path, d.Name())
//		return nil
//	}))
func (n *NoGo) ForWalkDir(fsys fs.FS, root string, ignoreFilename string, fn fs.WalkDirFunc) (fs.FS, string, fs.WalkDirFunc) {
	return fsys, root, func(path string, d fs.DirEntry, err error) error {
		ok, err := n.WalkFunc(fsys, ignoreFilename, path, d.IsDir(), err)
		if err != nil {
			return err
		}

		if ok {
			return fn(path, d, err)
		}

		return nil
	}
}
