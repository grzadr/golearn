package parsing

import (
	"io/fs"
	"iter"
	"path"
	"strings"
)

type WalkEntry struct {
	Path      string // path relative to root
	Name      string // name
	Ext       string // extension for files
	IsDir     bool   // info about being a directory
	IsRegular bool   // info about being a regular file
}

func (e *WalkEntry) isFile() bool {
	return e.IsRegular
}

func (e *WalkEntry) isFileWithExt(ext string) bool {
	return e.isFile() && e.Ext == ext
}

func (e *WalkEntry) isJSONFile() bool {
	return e.isFileWithExt((".json"))
}

func newWalkEntry(entry fs.DirEntry, root string) WalkEntry {
	ext := path.Ext(entry.Name())
	return WalkEntry{
		Path:      strings.TrimPrefix(path.Join(root, entry.Name()), "./"),
		Name:      strings.TrimSuffix(entry.Name(), ext),
		Ext:       path.Ext(entry.Name()),
		IsDir:     entry.IsDir(),
		IsRegular: entry.Type().IsRegular(),
	}
}

func walkFS(fsys fs.FS, root string) iter.Seq2[WalkEntry, error] {
	return func(yield func(WalkEntry, error) bool) {
		// Read directory entries
		entries, err := fs.ReadDir(fsys, root)
		if err != nil {
			yield(WalkEntry{}, err)
			return
		}

		// Iterate over entries
		for _, entry := range entries {
			// Create WalkEntry
			walkEntry := newWalkEntry(entry, root)

			// Yield the current entry
			if !yield(walkEntry, nil) {
				return
			}

			// If it's a directory, recursively walk it
			if entry.IsDir() {
				subPath := path.Join(root, entry.Name())
				subIter := walkFS(fsys, subPath)

				// Create wrapper to handle the subdirectory iteration
				subIter(func(subEntry WalkEntry, err error) bool {
					return yield(subEntry, err)
				})
			}
		}
	}
}
