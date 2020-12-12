package command

type Command func() error

func StartCommandFunc() Command {
	return func() error {
		return nil
	}
}

func ArchiveCommandFunc() Command {
	return func() error {
		return nil
	}
}
