package notify

import (
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (n Notify) Send(subject, message string) error {
	if n.Disabled {
		return nil
	}

	var eg errgroup.Group

	for _, service := range n.notifiers {
		if service != nil {
			s := service
			eg.Go(func() error {
				return s.Send(subject, message)
			})
		}
	}

	err := eg.Wait()
	if err != nil {
		err = errors.Wrap(ErrSendNotification, err.Error())
	}

	return err
}
