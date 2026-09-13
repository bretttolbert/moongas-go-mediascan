package shared

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFilesFiltersMediaFiles(t *testing.T) {
	root := t.TempDir()
	albumDir := filepath.Join(root, "Artist", "Album")
	if err := os.MkdirAll(albumDir, 0o755); err != nil {
		t.Fatalf("create album directory: %v", err)
	}

	trackPath := filepath.Join(albumDir, "Track.mp3")
	if err := os.WriteFile(trackPath, []byte("not a complete mp3"), 0o644); err != nil {
		t.Fatalf("create track: %v", err)
	}
	if err := os.WriteFile(filepath.Join(albumDir, "cover.jpg"), []byte("cover"), 0o644); err != nil {
		t.Fatalf("create cover: %v", err)
	}
	excludedPath := filepath.Join(albumDir, "Excluded.mp3")
	if err := os.WriteFile(excludedPath, []byte("not a complete mp3"), 0o644); err != nil {
		t.Fatalf("create excluded track: %v", err)
	}

	files := ScanFiles(MediascanConf{
		MediaDirs:    []string{root},
		MediaExts:    []string{".mp3"},
		ExcludePaths: []string{"Excluded.mp3"},
	})

	if len(files) != 1 {
		t.Fatalf("expected 1 media file, got %d", len(files))
	}
	if files[0].Path != trackPath {
		t.Errorf("got path %q, want %q", files[0].Path, trackPath)
	}
	if files[0].Title != "Track" {
		t.Errorf("got title %q, want %q", files[0].Title, "Track")
	}
}
