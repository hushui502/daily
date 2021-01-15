package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v33/github"
	"io"
	"io/fs"
	"sync"
	"time"
)

// Ensure io/fs interfaces are implemented
var ErrNotLoaded = fmt.Errorf("gist not loaded: %w", fs.ErrInvalid)

// FS represents a filesystem based on a Github Gist.
type FS struct {
	id string
	client *github.Client
	gist *github.Gist
	mu sync.RWMutex
}

// New returns a FS based on a given Gist ID, without the username portion.
// Example "https://gist.github.com/hushui502/e7d5d645bfe2b3d696e16da98a5c8871"
//		id=e7d5d645bfe2b3d696e16da98a5c8871
func New(id string) *FS {
	return &FS{
		client: github.NewClient(nil),
		id:     id,
	}
}

func NewWithClient(client *github.Client, id string) *FS {
	return &FS{
		client: client,
		id:     id,
	}
}

func (fsys *FS) GetID() string {
	return fsys.id
}

func (fsys *FS) Load(ctx context.Context) error {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()

	gist, _, err := fsys.client.Gists.Get(ctx, fsys.id)
	if err != nil {
		return err
	}

	fsys.gist = gist

	return nil
}

// file represents a file stored in a gist and implements fs.File methods.
// It is built out of a github.GistFile
type file struct {
	gistFile *github.GistFile
	modtime time.Time
	reader io.Reader
	mu sync.Mutex
}

// Open opens the name file for reading and return implements fs.File methods.
func (fsys *FS) Open(name string) (fs.File, error) {
	fsys.mu.RLock()
	defer fsys.mu.RUnlock()

	if fsys.gist == nil {
		return nil, ErrNotLoaded
	}

	if name == "./" || name == "." {
		return fsys.openRoot(), nil
	}

	f, ok := fsys.gist.Files[github.GistFilename(name)]
	if !ok {
		return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
	}

	return fsys.wrapFile(&f), nil
}

// wrapFile wraps a github.GistFile into a file which implements the fs.File interfaces
func (fsys *FS) wrapFile(f *github.GistFile) *file {
	return &file{
		gistFile: f,
		modtime:  fsys.gist.GetUpdatedAt(),
		reader:   bytes.NewReader([]byte(f.GetContent())),
	}
}

// ReadFile reads and returns the content of the named file
func (fsys *FS) ReadFile(name string) ([]byte, error) {
	fsys.mu.RLock()
	defer fsys.mu.RUnlock()

	if fsys.gist == nil {
		return nil, ErrNotLoaded
	}

	gistFile, ok := fsys.gist.Files[github.GistFilename(name)]
	if !ok {
		return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
	}

	return []byte(gistFile.GetContent()), nil
}

// ReadDir reads and returns the entire named directory which contains all files
// that are stored in the Gist supporting the filesystem
func (fsys *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	fsys.mu.RLock()
	defer fsys.mu.RUnlock()

	if fsys.gist == nil {
		return nil, ErrNotLoaded
	}

	if name == "." || name == "./" {
		return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
	}

	return fsys.openRoot().(*rootDir).ReadDir(-1)
}

func (f *file) isClosed() bool {
	return f.reader == nil
}

func (f *file) Read(b []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.isClosed() {
		return 0, fs.ErrClosed
	}

	return f.reader.Read(b)
}

func (f *file) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.gistFile = nil
	f.reader = nil

	return nil
}

// Stat provides stat about the file. The modtime notably, is set to
// when the underlying Gist was last updated.
func (f *file) Stat() (fs.FileInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.isClosed() {
		return nil, fs.ErrClosed
	}

	return f, nil
}

func (f *file) Name() string { return f.gistFile.GetFilename() }
func (f *file) Size() int64  { return int64(f.gistFile.GetSize()) }

// Mode always return 0444.
func (f *file) Mode() fs.FileMode { return 0444 }

// ModTime always return the time of the underlying gist last update.
func (f *file) ModTime() time.Time { return f.modtime }

func (f *file) IsDir() bool                { return false }
func (f *file) Sys() interface{}           { return f.gistFile }
func (f *file) Type() fs.FileMode          { return f.Mode().Type() }
func (f *file) Info() (fs.FileInfo, error) { return f, nil }

func (f *file) ReadDir(count int) ([]fs.DirEntry, error) {
	return nil, &fs.PathError{
		Op:   "read",
		Path: f.Name(),
		Err:  errors.New("is not a directory"),
	}
}

type rootDir struct {
	files []*file
	offset int
	modtime time.Time
	mu sync.Mutex
}

func (fsys *FS) openRoot() fs.File {
	files := make([]*file, 0, len(fsys.gist.Files))
	for _, f := range fsys.gist.Files {
		files = append(files, fsys.wrapFile(&f))
	}

	return &rootDir{
		files:   files,
		modtime: fsys.gist.GetUpdatedAt(),
	}
}

func (d *rootDir) Close() error               { return nil }
func (d *rootDir) Stat() (fs.FileInfo, error) { return d, nil }
func (d *rootDir) Name() string               { return "./" }
func (d *rootDir) Size() int64                { return 0 }
func (d *rootDir) Mode() fs.FileMode          { return fs.ModeDir | 0444 }

// ModTime always return the time of the underlying gist last update.
func (d *rootDir) ModTime() time.Time { return d.modtime }

func (d *rootDir) IsDir() bool       { return true }
func (d *rootDir) Type() fs.FileMode { return d.Mode().Type() }
func (d *rootDir) Sys() interface{}  { return nil }

func (d *rootDir) Read(b []byte) (int, error) {
	return 0, &fs.PathError{
		Op:   "read",
		Path: d.Name(),
		Err:  errors.New("is a directory"),
	}
}

func (d *rootDir) ReadDir(count int) ([]fs.DirEntry, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	n := len(d.files) - d.offset

	if count > 0 && n > count {
		n = count
	}

	if n == 0 {
		if count <= 0 {
			return nil, nil
		}
	}

	files := make([]fs.DirEntry, n)
	for i := range files {
		files[i] = d.files[d.offset+i]
	}
	d.offset += n

	return files, nil
}
