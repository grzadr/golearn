package parsing

import (
	"sort"
	"testing"
	"testing/fstest"
)

func TestWalkFS_EmptyDirectory(t *testing.T) {
	fsys := fstest.MapFS{}

	var entries []WalkEntry
	for entry, err := range walkFS(fsys, ".") {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		entries = append(entries, entry)
	}

	if len(entries) != 0 {
		t.Errorf("expected empty directory, got %d entries", len(entries))
	}
}

func TestWalkFS_SingleFile(t *testing.T) {
	fsys := fstest.MapFS{
		"test.txt": &fstest.MapFile{Data: []byte("content")},
	}

	var entries []WalkEntry
	for entry, err := range walkFS(fsys, ".") {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		entries = append(entries, entry)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	expected := WalkEntry{
		Path:      "test.txt",
		Name:      "test",
		Ext:       ".txt",
		IsDir:     false,
		IsRegular: true,
	}

	if !compareWalkEntries(entries[0], expected) {
		t.Errorf("\ngot:  %+v\nwant: %+v", entries[0], expected)
	}
}

func TestWalkFS_MultipleFiles(t *testing.T) {
	fsys := fstest.MapFS{
		"file1.txt":  &fstest.MapFile{Data: []byte("content")},
		"file2.json": &fstest.MapFile{Data: []byte("content")},
	}

	var entries []WalkEntry
	for entry, err := range walkFS(fsys, ".") {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		entries = append(entries, entry)
	}

	expected := []WalkEntry{
		{Path: "file1.txt", Name: "file1", Ext: ".txt", IsRegular: true},
		{Path: "file2.json", Name: "file2", Ext: ".json", IsRegular: true},
	}

	compareEntryLists(t, entries, expected)
}

func TestWalkFS_SingleDirectory(t *testing.T) {
	fsys := fstest.MapFS{
		"dir/file.txt": &fstest.MapFile{Data: []byte("content")},
	}

	var entries []WalkEntry
	for entry, err := range walkFS(fsys, ".") {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		entries = append(entries, entry)
	}

	expected := []WalkEntry{
		{Path: "dir", Name: "dir", IsDir: true},
		{Path: "dir/file.txt", Name: "file", Ext: ".txt", IsRegular: true},
	}

	compareEntryLists(t, entries, expected)
}

func TestWalkFS_NestedDirectories(t *testing.T) {
	fsys := fstest.MapFS{
		"dir1/file1.txt":       &fstest.MapFile{Data: []byte("content")},
		"dir1/dir2/file2.txt":  &fstest.MapFile{Data: []byte("content")},
		"dir1/dir2/file3.json": &fstest.MapFile{Data: []byte("content")},
	}

	var entries []WalkEntry
	for entry, err := range walkFS(fsys, ".") {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		entries = append(entries, entry)
	}

	expected := []WalkEntry{
		{Path: "dir1", Name: "dir1", IsDir: true},
		{Path: "dir1/dir2", Name: "dir2", IsDir: true},
		{Path: "dir1/dir2/file2.txt", Name: "file2", Ext: ".txt", IsRegular: true},
		{Path: "dir1/dir2/file3.json", Name: "file3", Ext: ".json", IsRegular: true},
		{Path: "dir1/file1.txt", Name: "file1", Ext: ".txt", IsRegular: true},
	}

	compareEntryLists(t, entries, expected)
}

func TestWalkFS_StartFromSubdirectory(t *testing.T) {
	fsys := fstest.MapFS{
		"dir/file1.txt":   &fstest.MapFile{Data: []byte("content")},
		"dir/file2.json":  &fstest.MapFile{Data: []byte("content")},
		"other/file3.txt": &fstest.MapFile{Data: []byte("content")},
	}

	var entries []WalkEntry
	for entry, err := range walkFS(fsys, "dir") {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		entries = append(entries, entry)
	}

	expected := []WalkEntry{
		{Path: "dir/file1.txt", Name: "file1", Ext: ".txt", IsRegular: true},
		{Path: "dir/file2.json", Name: "file2", Ext: ".json", IsRegular: true},
	}

	compareEntryLists(t, entries, expected)
}

func TestWalkFS_NonexistentDirectory(t *testing.T) {
	fsys := fstest.MapFS{}

	it := walkFS(fsys, "nonexistent")

	for entry, err := range it {
		if err == nil {
			t.Error("expected error for nonexistent directory, got entry:", entry)
		}
		return // We only need to check the first iteration
	}
}

// Helper function to compare two WalkEntry lists
func compareEntryLists(t *testing.T, got, want []WalkEntry) {
	t.Helper()

	// Sort both slices for comparison
	sort.Slice(got, func(i, j int) bool {
		return got[i].Path < got[j].Path
	})
	sort.Slice(want, func(i, j int) bool {
		return want[i].Path < want[j].Path
	})

	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d", len(got), len(want))
	}

	for i := range got {
		if !compareWalkEntries(got[i], want[i]) {
			t.Errorf("entry %d:\ngot:  %+v\nwant: %+v", i, got[i], want[i])
		}
	}
}

// Helper function to compare individual WalkEntries
func compareWalkEntries(got, want WalkEntry) bool {
	return got.Path == want.Path &&
		got.Name == want.Name &&
		got.Ext == want.Ext &&
		got.IsDir == want.IsDir &&
		got.IsRegular == want.IsRegular
}

func TestWalkFS_EarlyTermination(t *testing.T) {
	// Create a filesystem with multiple files
	fsys := fstest.MapFS{
		"dir1/file1.txt":      &fstest.MapFile{Data: []byte("content")},
		"dir1/dir2/file2.txt": &fstest.MapFile{Data: []byte("content")},
		"dir1/dir2/file3.txt": &fstest.MapFile{Data: []byte("content")},
	}

	// Count how many entries we process before stopping
	processedEntries := 0

	// Track if we found our target file
	foundTarget := false

	// Walk the filesystem but stop after finding "file2.txt"
	for entry, err := range walkFS(fsys, ".") {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		processedEntries++

		if entry.Name == "file2" && entry.Ext == ".txt" {
			foundTarget = true
			break // This will cause yield to return false on next iteration
		}
	}

	// Verify we found our target
	if !foundTarget {
		t.Error("did not find target file before termination")
	}

	// Verify we didn't process all entries
	// The full tree has 5 entries (dir1, dir1/dir2, and 3 files)
	if processedEntries >= 5 {
		t.Error("early termination failed: processed all entries")
	}
}
