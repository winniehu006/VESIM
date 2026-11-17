package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureSQLiteDirectoryCreatesParent(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "nested", "deeper", "vesim.db")
	if err := ensureSQLiteDirectory(dsn); err != nil {
		t.Fatalf("ensureSQLiteDirectory() error = %v", err)
	}
	info, err := os.Stat(filepath.Dir(dsn))
	if err != nil {
		t.Fatalf("parent directory missing: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("parent path is not a directory")
	}
}

func TestEnsureSQLiteDirectoryHandlesDSNForms(t *testing.T) {
	directory := t.TempDir()
	cases := []struct {
		name    string
		dsn     string
		wantDir string
	}{
		{
			name:    "query parameters",
			dsn:     filepath.Join(directory, "a", "vesim.db") + "?cache=shared&mode=rwc",
			wantDir: filepath.Join(directory, "a"),
		},
		{
			name:    "file prefix",
			dsn:     "file:" + filepath.Join(directory, "b", "vesim.db"),
			wantDir: filepath.Join(directory, "b"),
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := ensureSQLiteDirectory(testCase.dsn); err != nil {
				t.Fatalf("ensureSQLiteDirectory(%q) error = %v", testCase.dsn, err)
			}
			if info, err := os.Stat(testCase.wantDir); err != nil || !info.IsDir() {
				t.Fatalf("directory %q was not created: %v", testCase.wantDir, err)
			}
		})
	}
}

func TestEnsureSQLiteDirectorySkipsNonFileDSNs(t *testing.T) {
	for _, dsn := range []string{"", "  ", ":memory:", "file::memory:?cache=shared", "vesim.db"} {
		if err := ensureSQLiteDirectory(dsn); err != nil {
			t.Fatalf("ensureSQLiteDirectory(%q) error = %v", dsn, err)
		}
	}
	// A bare file name must not be turned into a directory.
	if info, err := os.Stat("vesim.db"); err == nil && info.IsDir() {
		t.Fatalf("bare DSN created a directory named vesim.db")
	}
}
