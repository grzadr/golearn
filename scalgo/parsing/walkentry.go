package parsing

import (
	"io/fs"
	"iter"
	"path"
)

type WalkEntry struct {
    Path     string // path relative to root
    Name     string // name
    Ext      string // extension for files
    IsDir    bool   // info about being a directory
    IsRegular bool  // info about being a regular file
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
            walkEntry := WalkEntry{
                Path:     path.Join(root, entry.Name()),
                Name:     entry.Name(),
                Ext:      path.Ext(entry.Name()),
                IsDir:    entry.IsDir(),
                IsRegular: entry.Type().IsRegular(),
            }

            // If root is ".", remove the leading "./"
            if root == "." {
                walkEntry.Path = walkEntry.Name
            }

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
