package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/module"
)

var cachedir = filepath.Join(os.Getenv("HOME"), "gomodproxy-cache")

func main() {
	if err := os.MkdirAll(cachedir, 0755); err != nil {
		log.Fatalf("creating cache: %v", err)
	}
	http.HandleFunc("/mod/", handleMod)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMod(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/mod/")

	// MODULE/@v/list
	if mod, ok := suffixed(path, "/@v/list"); ok {
		mod, err := module.UnescapePath(mod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("list", mod)

		versions, err := listVersions(mod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		for _, v := range versions {
			fmt.Fprintln(w, v)
		}
		return
	}

	// MODULE/@latest
	if mod, ok := suffixed(path, "/@latest"); ok {
		mod, err := module.UnescapePath(mod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("latest", mod)

		latest, err := resolve(mod, "latest")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		info := InfoJSON{Version: latest.Version, Time: latest.Time}
		json.NewEncoder(w).Encode(info)
		return
	}

	// MODULE/@v/VERSION.{info,mod,zip}
	if rest, ext, ok := lastCut(path, "."); ok && isOneOf(ext, "mod", "info", "zip") {
		if mod, version, ok := cut(rest, "/@v/"); ok {
			mod, err := module.UnescapePath(mod)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			version, err := module.UnescapeVersion(version)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			log.Printf("%s %s@%s", ext, mod, version)

			m, err := download(mod, version)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			// The version may be a query such as a branch name.
			// Branches move, so we suppress HTTP caching in that case.
			// (To avoid repeated calls to download, the proxy could use
			// the module name and resolved m.Version as a key in a cache.)
			if version != m.Version {
				w.Header().Set("Cache-Control", "no-store")
				log.Printf("%s %s@%s => %s", ext, mod, version, m.Version)
			}

			// Return the relevant cached file.
			var filename, mimetype string
			switch ext {
			case "info":
				filename = m.Info
				mimetype = "application/json"
			case "mod":
				filename = m.GoMod
				mimetype = "text/plain; charset=UTF-8"
			case "zip":
				filename = m.Zip
				mimetype = "application/zip"
			}
			w.Header().Set("Content-Type", mimetype)
			if err := copyFile(w, filename); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "bad request", http.StatusBadRequest)
}

// download runs 'go mod download' and returns information about a
// specific module version. It also downloads the module's dependencies.
func download(name, version string) (*ModuleDownloadJSON, error) {
	var mod ModuleDownloadJSON
	if err := runGo(&mod, "mod", "download", "-json", name+"@"+version); err != nil {
		return nil, err
	}
	if mod.Error != "" {
		return nil, fmt.Errorf("failed to download module %s: %v", name, mod.Error)
	}

	return &mod, nil
}

// listVersions runs 'go list -m -versions' and returns an unordered list
// of versions of the specified module.
func listVersions(name string) ([]string, error) {
	var mod ModuleListJSON
	if err := runGo(&mod, "list", "-m", "-json", "-versions", name); err != nil {
		return nil, err
	}
	if mod.Error != nil {
		return nil, fmt.Errorf("failed to list module %s: %v", name, mod.Error.Err)
	}

	return mod.Versions, nil
}

func resolve(name, query string) (*ModuleListJSON, error) {
	var mod ModuleListJSON
	if err := runGo(&mod, "list", "-m", "-json", name+"@"+query); err != nil {
		return nil, err
	}
	if mod.Error != nil {
		return nil, fmt.Errorf("failed to list module %s: %v", name, mod.Error.Err)
	}

	return &mod, nil
}

// runGo runs the Go command and decodes its JSON output into result.
func runGo(result interface{}, args ...string) error {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	cmd := exec.Command("go", args...)
	cmd.Dir = tmpdir
	cmd.Env = []string{
		"USER=" + os.Getenv("USER"),
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"NETRC=", // don't allow go command to read user's secrets
		"GOPROXY=direct",
		"GOCACHE=" + cachedir,
		"GOMODCACHE=" + cachedir,
		"GOSUMDB=",
	}
	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %v (stderr=<<%s>>)", cmd, err, cmd.Stderr)
	}
	
	if err := json.Unmarshal(cmd.Stdout.(*bytes.Buffer).Bytes(), result); err != nil {
		return fmt.Errorf("internal error decoding %s JSON output: %v", cmd, err)
	}
	
	return nil
}

// ModuleDownloadJSON is the JSON schema of the output of 'go help mod download'.
type ModuleDownloadJSON struct {
	Path     string // module path
	Version  string // module version
	Error    string // error loading module
	Info     string // absolute path to cached .info file
	GoMod    string // absolute path to cached .mod file
	Zip      string // absolute path to cached .zip file
	Dir      string // absolute path to cached source root directory
	Sum      string // checksum for path, version (as in go.sum)
	GoModSum string // checksum for go.mod (as in go.sum)
}

// ModuleListJSON is the JSON schema of the output of 'go help list'.
type ModuleListJSON struct {
	Path      string          // module path
	Version   string          // module version
	Versions  []string        // available module versions (with -versions)
	Replace   *ModuleListJSON // replaced by this module
	Time      *time.Time      // time version was created
	Update    *ModuleListJSON // available update, if any (with -u)
	Main      bool            // is this the main module?
	Indirect  bool            // is this module only an indirect dependency of main module?
	Dir       string          // directory holding files for this module, if any
	GoMod     string          // path to go.mod file used when loading this module, if any
	GoVersion string          // go version used in module
	Retracted string          // retraction information, if any (with -retracted or -u)
	Error     *ModuleError    // error loading module
}

type ModuleError struct {
	Err string
}

// InfoJSON is the JSON schema of the .info and @latest endpoints.
type InfoJSON struct {
	Version string
	Time *time.Time
}

// suffixed reports whether x has the specified suffix,
// and returns the prefix.
func suffixed(x, suffix string) (rest string, ok bool) {
	if y := strings.TrimSuffix(x, suffix); y != x {
		return y, true
	}
	return
}

// See https://github.com/golang/go/issues/46336
func cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func lastCut(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

// copyFile writes the content of the named file to dest.
func copyFile(dest io.Writer, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(dest, f)
	return err
}

func isOneOf(s string, items ...string) bool {
	for _, item := range items {
		if s == item {
			return true
		}
	}
	return false
}
