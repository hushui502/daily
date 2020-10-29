package v1

import (
	"testing"
	"time"
)

func TestExample_TimeToGo(t *testing.T) {
	nowMock := func() time.Time {
		t, _ := time.Parse(time.Kitchen, time.Kitchen)
		return t
	}

	e := Example{nowMock}

	if e.TimeToGo() != "its time to go! 0000-01-01 15:04:00 +0000 UTC" {
		t.Errorf("test failed")
	}

}
