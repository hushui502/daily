package hsfs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// ensure file system implements interface
var _ fs.FS = (*FS)(nil)


// 文件系统可以hash文件名
// 可以使用缓存
// 如果文件数据改变，文件名相对也会改变，之前的缓存也会失效
type FS struct {
	fsys fs.FS

	mu sync.RWMutex
	m map[string]string		// lookup (path to hash path)
	r map[string]string		// reverse lookup (hash path to path)
}

func NewFS(fsys fs.FS) *FS {
	return &FS{
		fsys: fsys,
		m:    make(map[string]string),
		r:    make(map[string]string),
	}
}

// 返回一个对命名文件的引用
// 如果name是一个hash名，则调用之前使用过的文件
func (fsys *FS) Open(name string) (fs.File, error) {
	f, _, err := fsys.open(name)
	return f, err
}

func (fsys *FS) open(name string) (_ fs.File, isHashed bool, err error) {
	// lookup original file by hashed name
	fsys.mu.RLock()
	if hashname, ok := fsys.r[name]; ok {
		name = hashname
		isHashed = true
	}
	fsys.mu.RUnlock()

	// parse filename to see if it contains a hash
	// if so, check if hash name matches
	if base, hash := ParseName(name); hash != "" {
		if fsys.HashName(base) == name {
			name = base
			isHashed = true
		}
	}

	f, err := fsys.fsys.Open(name)

	return f, isHashed, err
}

// 如果文件存在，返回一个hash name
// 否则返回原文件名
func (fsys *FS) HashName(name string) string {
	// lookup cached formatted name, if exists
	fsys.mu.RLock()
	if s := fsys.m[name]; s != "" {
		fsys.mu.RUnlock()
		return s
	}
	fsys.mu.RUnlock()

	// read file contents, return original filename if we receive an error
	buf, err := fs.ReadFile(fsys.fsys, name)
	if err != nil {
		return name
	}

	// compute hash and build filename
	hash := sha256.Sum224(buf)
	hashname := FormatName(name, hex.EncodeToString(hash[:]))

	// store in lookups
	fsys.mu.Lock()
	fsys.m[name] = hashname
	fsys.r[hashname] = name
	fsys.mu.Unlock()

	return hashname
}


// 格式化一个hash name
// 根据是否有扩展名分别有不同的逻辑处理
func FormatName(filename, hash string) string {
	if filename == "" {
		return ""
	} else if hash == "" {
		return filename
	}

	dir, base := path.Split(filename)
	if i := strings.Index(base, "."); i != -1 {
		return path.Join(dir, fmt.Sprintf("%s-%s%s", base[:i], hash, base[i:]))
	}

	return path.Join(dir, fmt.Sprintf("%s-%s", base, hash))
}


// 解析名字，返回base和hash两部分
func ParseName(filename string) (base, hash string) {
	if filename == "" {
		return "", ""
	}

	dir, base := path.Split(filename)

	pre, ext := base, ""
	if i := strings.Index(base, "."); i != -1 {
		pre = base[:i]
		ext = base[i:]
	}

	if !hashSuffixRegex.MatchString(pre) {
		return filename, ""
	}

	return path.Join(dir, pre[:len(pre)-65]+ext), pre[len(pre)-64:]
}

var hashSuffixRegex = regexp.MustCompile(`-[0-9a-f]{64}`)

func FileServer(fsys fs.FS) http.Handler {
	hfsys, ok := fsys.(*FS)
	if !ok {
		hfsys = NewFS(fsys)
	}

	return &fsHandler{fsys: hfsys}
}

type fsHandler struct {
	fsys *FS
}

func (h *fsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path
	if filename == "/" {
		filename = "."
	} else {
		filename = strings.TrimPrefix(filename, "/")
	}
	filename = path.Clean(filename)

	// read file from attached file system
	f, isHashed, err := h.fsys.open(filename)
	if os.IsNotExist(err) {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "505 Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// fetch file info, disallow directories from being displayed.
	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "500 Interval Server Error", http.StatusInternalServerError)
		return
	} else if fi.IsDir() {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	// Determined content type based on file extension.
	if ext := path.Ext(filename); ext != "" {
		w.Header().Set("Content-Type", mime.TypeByExtension(ext))
	}

	// Cache the file aggressively if the file contains a hash
	if isHashed {
		w.Header().Set("Cache-Control", `public, max-age=31536000`)
	}

	// Set content length
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

	// Flush header and write content
	w.WriteHeader(http.StatusOK)
	if r.Method != "HEAD" {
		io.Copy(w, f)
	}
}