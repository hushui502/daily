package notify

func (n *Notify) useService(service Notifier) {
	if service == nil {
		return
	}
	n.notifiers = append(n.notifiers, service)
}

func (n *Notify) UseServices(services...Notifier) {
	for _, s := range services {
		n.useService(s)
	}
}