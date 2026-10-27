package main

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestBuildVersion(t *testing.T) {
	tests := []struct {
		name string
		info *debug.BuildInfo
		ok   bool
		want string
	}{
		{name: "unavailable", want: "unknown"},
		{name: "nil info", ok: true, want: "unknown"},
		{name: "failed read", info: &debug.BuildInfo{}, want: "unknown"},
		{name: "empty version", info: &debug.BuildInfo{}, ok: true, want: "(devel)"},
		{name: "development", info: &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, ok: true, want: "(devel)"},
		{name: "release", info: &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}}, ok: true, want: "v1.2.3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildVersion(tt.info, tt.ok); got != tt.want {
				t.Fatalf("buildVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestImportExport(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	config, err := filepath.Abs(filepath.Join("..", "..", "bibfuse.toml"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	t.Chdir(dir)
	input := "input.bib"
	const url = "https://example.com/reference"
	if err := os.WriteFile(input, []byte("@misc{offline, title={Offline entry}, url={"+url+"}}"), 0o600); err != nil {
		t.Fatal(err)
	}
	opts := options{
		config:  config,
		dbFile:  "bib.db",
		outFile: "out.bib",
	}
	// Exercise relative and absolute paths, duplicate imports, and database reopen.
	for _, files := range [][]string{{input}, {filepath.Join(dir, input)}, nil} {
		if err := run(opts, files); err != nil {
			t.Fatal(err)
		}
		output, err := os.ReadFile(opts.outFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(output), "@misc{offline,") || !strings.Contains(string(output), url) {
			t.Fatalf("missing entry or literal URL in output: %s", output)
		}
		if strings.Count(string(output), "@misc{offline,") != 1 {
			t.Fatalf("duplicate entry in output: %s", output)
		}
		opts.dbFile = filepath.Join(dir, "bib.db")
		opts.outFile = filepath.Join(dir, "out.bib")
	}
}
