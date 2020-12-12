package observer

import "fmt"

type ISubject interface {
	Register(observer IObserver)
	Remove(observer IObserver)
	Notify(msg string)
}

type IObserver interface {
	Update(msg string)
}

type Subject struct {
	observers []IObserver
}

func (s Subject) Register(observer IObserver) {
	s.observers = append(s.observers, observer)
}

func (s Subject) Remove(observer IObserver) {
	for i, ob := range s.observers {
		if ob == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
		}
	}
}

func (s Subject) Notify(msg string) {
	for _, o := range s.observers {
		o.Update(msg)
	}
}

type Observer1 struct {

}

func (Observer1) Update(msg string) {
	fmt.Printf("Observer1: %s\n", msg)
}

// Obsever2 Obsever2
type Observer2 struct{}

// Update 实现观察者接口
func (Observer2) Update(msg string) {
	fmt.Printf("Obsever2: %s\n", msg)
}
