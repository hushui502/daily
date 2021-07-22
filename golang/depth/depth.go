package depth

import (
	"errors"
	"go/build"
	"os"
)

// ErrRootPkgNotResolved is returned when the root Pkg of the Tree cannot be resolved,
// typically because it does not exist.
var ErrRootPkgNotResolved = errors.New("unable to resolve root package")

// Importer defines a type that can import a package and return its details.
type Importer interface {
	Import(name, srcDir string, im build.ImportMode) (*build.Package, error)
}

// Tree represents the top level of a Pkg and the configuration used to
// initialize and represent its contents.
type Tree struct {
	Root *Pkg

	ResolveInternal bool
	ResolveTest     bool
	MaxDepth        int

	Importer Importer

	importCache map[string]struct{}
}

func (t *Tree) Resolve(name string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	t.Root = &Pkg{
		Name:   name,
		Tree:   t,
		SrcDir: pwd,
		Test:   false,
	}

	// Reset the import cache each time to ensure a reused Tree doesn't
	// reuse the same cache.
	t.importCache = nil

	// Allow custom importers, but use build.Default if none is provided.
	if t.Importer == nil {
		t.Importer = &build.Default
	}

	t.Root.Resolve(t.Importer)
	if !t.Root.Resolved {
		return ErrRootPkgNotResolved
	}

	return nil
}

func (t *Tree) shouldResolveInternal(parent *Pkg) bool {
	if t.ResolveInternal {
		return true
	}

	return parent == t.Root
}

func (t *Tree) isAtMaxDepth(p *Pkg) bool {
	if t.MaxDepth == 0 {
		return false
	}

	return p.depth() >= t.MaxDepth
}

func (t *Tree) hasSeenImport(name string) bool {
	if t.importCache == nil {
		t.importCache = make(map[string]struct{})
	}

	if _, ok := t.importCache[name]; ok {
		return true
	}
	t.importCache[name] = struct{}{}
	return false
}
