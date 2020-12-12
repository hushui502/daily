package command

type ICommand interface {
	Execute() error
}

type StartCommand struct {}

func NewStartCommand() *StartCommand {
	return &StartCommand{}
}

func (c *StartCommand) Execute() error {
	return nil
}

type ArchiveCommand struct {}

func NewArchiveCommand() *StartCommand {
	return &StartCommand{}
}
func (c *ArchiveCommand) Execute() error {
	return nil
}
