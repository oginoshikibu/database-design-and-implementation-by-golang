package file

import (
	"os"
	"path/filepath"
	"testing"
)

// Fig. 3.12 Testing the SimpleDB file manager
func TestFileMgr(t *testing.T) {
	// Setup: create a temporary directory for testing
	testDir := filepath.Join(os.TempDir(), "filetest")
	defer func() { _ = os.RemoveAll(testDir) }()

	// Create FileMgr with block size 400
	fm, err := NewFileMgr(testDir, 400)
	if err != nil {
		t.Fatalf("failed to create FileMgr: %v", err)
	}

	// Create a block reference
	blk := NewBlockID("testfile", 2)

	// Create a page and write data to it
	p1 := NewPage(fm.BlockSize())

	pos1 := 88
	testString := "abcdefghijklm"
	p1.SetString(pos1, testString)

	size := MaxLength(len(testString))
	pos2 := pos1 + size
	p1.SetInt(pos2, 345)

	// Write the page to the block
	err = fm.Write(blk, p1)
	if err != nil {
		t.Fatalf("failed to write block: %v", err)
	}

	// Read the block into a new page
	p2 := NewPage(fm.BlockSize())
	err = fm.Read(blk, p2)
	if err != nil {
		t.Fatalf("failed to read block: %v", err)
	}

	// Verify the data
	gotInt := p2.GetInt(pos2)
	if gotInt != 345 {
		t.Errorf("offset %d: expected 345, got %d", pos2, gotInt)
	}

	gotString := p2.GetString(pos1)
	if gotString != testString {
		t.Errorf("offset %d: expected '%s', got '%s'", pos1, testString, gotString)
	}

	t.Logf("offset %d contains %d", pos2, gotInt)
	t.Logf("offset %d contains %s", pos1, gotString)
}
