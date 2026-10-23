package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Regenerate the expected output after an intentional change with:
//
//	go test ./cmd/logseq-claude-indexer -run TestGolden -update
//
// then review the diff of testdata/golden/ like any other code change.
var update = flag.Bool("update", false, "rewrite testdata/golden from the current output")

const (
	goldenRepo = "../../testdata/synthetic"
	goldenDir  = "../../testdata/golden"
)

// goldenNow is a fixed clock for the golden run, the morning after the last
// synthetic journal (2025_11_03), so the recent-timeline and upcoming windows
// are exercised.
var goldenNow = time.Date(2025, 11, 4, 9, 0, 0, 0, time.UTC)

func runGolden(t *testing.T, outDir string) {
	t.Helper()
	repo, err := filepath.Abs(goldenRepo)
	if err != nil {
		t.Fatal(err)
	}
	if err := generate(repo, outDir, goldenNow, log.New(io.Discard, "", 0)); err != nil {
		t.Fatalf("generate: %v", err)
	}
}

func readDir(t *testing.T, dir string) map[string]string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	files := make(map[string]string)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		files[e.Name()] = string(b)
	}
	return files
}

// TestGolden compares every generated index against a checked-in expected
// copy. Unit tests check that the code does what was intended; this checks
// what it actually produces, so a change in output shows up as a reviewable
// diff rather than being noticed later in a real knowledge base.
func TestGolden(t *testing.T) {
	out := t.TempDir()
	runGolden(t, out)
	got := readDir(t, out)

	if *update {
		if err := os.RemoveAll(goldenDir); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(goldenDir, 0755); err != nil {
			t.Fatal(err)
		}
		for name, content := range got {
			if err := os.WriteFile(filepath.Join(goldenDir, name), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
		return
	}

	want := readDir(t, goldenDir)
	for name, w := range want {
		g, ok := got[name]
		if !ok {
			t.Errorf("%s: expected file was not generated", name)
			continue
		}
		if g != w {
			t.Errorf("%s: output differs from golden copy; run with -update and review the diff if the change is intended", name)
		}
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			t.Errorf("%s: generated but has no golden copy; run with -update", name)
		}
	}
}

// TestGenerateDeterministic runs the full pipeline twice with the same clock
// and requires byte-identical output. It guards against map-iteration order,
// unstable sorts and hidden clock reads, all of which have caused churn in
// real repos before.
func TestGenerateDeterministic(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	runGolden(t, a)
	runGolden(t, b)
	first, second := readDir(t, a), readDir(t, b)
	if len(first) != len(second) {
		t.Fatalf("runs produced %d and %d files", len(first), len(second))
	}
	for name, content := range first {
		if second[name] != content {
			t.Errorf("%s: two runs with the same input and clock produced different output", name)
		}
	}
}
