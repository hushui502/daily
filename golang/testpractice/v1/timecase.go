package v1

import (
	"fmt"
	"time"
)

type Example struct {
	now func() time.Time
}

func NewExample() *Example {
	return &Example{now: time.Now}
}

func (e *Example) TimeToGo() string {
	return fmt.Sprintf("its time to go! %s", e.now().String())
}
