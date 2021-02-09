package notify

type Notifier interface {
	Send(string, string) error
	AddReceivers(...string)
}