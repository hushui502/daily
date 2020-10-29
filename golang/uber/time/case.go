package main

import "time"

func main() {
	poll(10 * time.Second)
}

func isActive(now, start, stop time.Time) bool {
	return (start.Before(now) || start.Equal(now)) && now.Before(stop)
}

func poll(delay time.Duration) {
	for {
		time.Sleep(delay)
	}
}
