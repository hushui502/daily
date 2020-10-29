package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofrs/flock"
)

func ExampleFlock_Locked() {
	f := flock.New(os.TempDir() + "/go-lock.lock")
	f.TryLock() // unchecked errors here

	fmt.Printf("locked: %v\n", f.Locked())

	f.Unlock()

	fmt.Printf("locked: %v\n", f.Locked())
	// Output: locked: true
	// locked: false
}

func ExampleFlock_TryLock() {
	// should probably put these in /var/lock
	fileLock := flock.New(os.TempDir() + "/go-lock.lock")

	locked, err := fileLock.TryLock()

	if err != nil {
		// handle locking error
	}

	if locked {
		fmt.Printf("path: %s; locked: %v\n", fileLock.Path(), fileLock.Locked())

		if err := fileLock.Unlock(); err != nil {
			// handle unlock error
		}
	}

	fmt.Printf("path: %s; locked: %v\n", fileLock.Path(), fileLock.Locked())
}

func ExampleFlock_TryLockContext() {
	// should probably put these in /var/lock
	fileLock := flock.New(os.TempDir() + "/go-lock.lock")

	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, 678*time.Millisecond)

	if err != nil {
		// handle locking error
	}

	if locked {
		fmt.Printf("path: %s; locked: %v\n", fileLock.Path(), fileLock.Locked())

		if err := fileLock.Unlock(); err != nil {
			// handle unlock error
		}
	}

	fmt.Printf("path: %s; locked: %v\n", fileLock.Path(), fileLock.Locked())
}
