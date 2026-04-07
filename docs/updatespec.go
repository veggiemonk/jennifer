//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	version, err := goVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get go version: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Detected Go version: %s\n", version)

	url := "https://raw.githubusercontent.com/golang/go/refs/tags/go" + version + "/doc/go_spec.html"
	fmt.Printf("Fetching %s\n", url)

	output := filepath.Join(filepath.Dir(os.Args[0]), "go_spec.html")
	// If run via "go run", os.Args[0] is a temp path. Use working directory instead.
	if wd, err := os.Getwd(); err == nil {
		output = filepath.Join(wd, "go_spec.html")
	}

	if err := fetch(url, output); err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch spec: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved to %s\n", output)
}

// goVersion returns the Go version string (e.g., "1.26.1") from "go version".
func goVersion() (string, error) {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", err
	}
	// Output: "go version go1.26.1 darwin/arm64"
	fields := strings.Fields(string(out))
	if len(fields) < 3 || !strings.HasPrefix(fields[2], "go") {
		return "", fmt.Errorf("unexpected go version output: %s", out)
	}
	return strings.TrimPrefix(fields[2], "go"), nil
}

func fetch(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
