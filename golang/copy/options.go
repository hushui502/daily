package copy

import "os"

type Options struct {
	OnSymlink func(src string) SymlinkAction
	Skip func(src string) (bool, error)
	AddPermission os.FileMode
	Sync bool
}

type SymlinkAction int

const (
	Deep SymlinkAction = iota
	Shallow
	Skip
)

func getDefaultOptions() Options {
	return Options{
		OnSymlink: func(src string) SymlinkAction {
			return Shallow
		},
		Skip: func(src string) (b bool, err error) {
			return false, nil
		},
		AddPermission: 0,
		Sync:          false,
	}
}

