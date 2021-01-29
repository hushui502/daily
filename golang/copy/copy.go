package copy

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	tmpPermissionForDirectory = os.FileMode(0755)
)

func Copy(src, dest string, opt ...Options) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	return switchboard(src, dest, info, assure(opt...))
}

func switchboard(src, dest string, info os.FileInfo, opt Options) error {
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return onsymlink(src, dest, info, opt)
	case info.IsDir():
		return dcopy(src, dest, info, opt)
	default:
		return fcopy(src, dest, info, opt)
	}
}

func copy(src, dest string, info os.FileInfo, opt Options) error {
	skip, err := opt.Skip(src)
	if err != nil {
		return err
	}

	if skip {
		return nil
	}

	return switchboard(src, dest, info, opt)
}

func fcopy(src, dest string, info os.FileInfo, opt Options) (err error) {
	if err = os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return
	}

	f, err := os.Create(dest)
	if err != nil {
		return
	}
	defer fclose(f, &err)

	if err = os.Chmod(f.Name(), info.Mode()|opt.AddPermission); err != nil {
		return
	}

	s, err := os.Open(src)
	if err != nil {
		return
	}
	defer fclose(s, &err)

	if _, err = io.Copy(f, s); err != nil {
		return
	}

	if opt.Sync {
		err = f.Sync()
	}

	return
}

func dcopy(srcdir, destdir string, info os.FileInfo, opt Options) (err error) {
	originalMode := info.Mode()

	if err = os.MkdirAll(destdir, tmpPermissionForDirectory); err != nil {
		return
	}
	defer chmod(destdir, originalMode|opt.AddPermission, &err)

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())

		if err = copy(cs, cd, content, opt); err != nil {
			return
		}
	}

	return
}

func onsymlink(src, dest string, info os.FileInfo, opt Options) error {
	switch opt.OnSymlink(src) {
	case Shallow:
		return lcopy(src, dest)
	case Deep:
		orig, err := os.Readlink(src)
		if err != nil {
			return err
		}
		info, err := os.Lstat(orig)
		if err != nil {
			return err
		}
		return copy(orig, dest, info, opt)
	case Skip:
		fallthrough
	default:
		return nil
	}
}

func lcopy(src, dest string) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}

func fclose(f *os.File, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}

func chmod(dir string, mode os.FileMode, reported *error) {
	if err := os.Chmod(dir, mode); *reported == nil {
		*reported = err
	}
}

func assure(opts ...Options) Options {
	if len(opts) == 0 {
		return getDefaultOptions()
	}
	defopt := getDefaultOptions()
	if opts[0].OnSymlink == nil {
		opts[0].OnSymlink = defopt.OnSymlink
	}
	if opts[0].Skip == nil {
		opts[0].Skip = defopt.Skip
	}

	return opts[0]
}