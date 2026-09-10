package shared

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveMediaDirs_WithRootDir(t *testing.T) {
	root := t.TempDir()
	musicDir := filepath.Join(root, "data", "Music")
	otherDir := filepath.Join(root, "data", "MusicOther")
	if err := os.MkdirAll(musicDir, 0o755); err != nil {
		t.Fatalf("mkdir music: %v", err)
	}
	if err := os.MkdirAll(otherDir, 0o755); err != nil {
		t.Fatalf("mkdir other: %v", err)
	}

	resolved, err := ResolveMediaDirs(root, []string{"data/Music", "data/MusicOther"})
	if err != nil {
		t.Fatalf("ResolveMediaDirs returned error: %v", err)
	}

	if len(resolved) != 2 {
		t.Fatalf("expected 2 resolved dirs, got %d", len(resolved))
	}
	if resolved[0] != musicDir || resolved[1] != otherDir {
		t.Fatalf("unexpected resolved dirs: %#v", resolved)
	}
}

func TestResolveMediaDirs_RejectsMissingRootDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	_, err := ResolveMediaDirs(missing, []string{"data/Music"})
	if err == nil {
		t.Fatal("expected error for missing root dir")
	}
}

func TestResolveMediaDirs_RejectsMissingJoinedDir(t *testing.T) {
	root := t.TempDir()

	_, err := ResolveMediaDirs(root, []string{"data/Music"})
	if err == nil {
		t.Fatal("expected error for missing assembled media dir")
	}
}
