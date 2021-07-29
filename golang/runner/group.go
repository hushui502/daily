package runner

type actor struct {
	execute   func() error
	interrupt func(error)
}

type Group struct {
	actors []actor
}

func (g *Group) Add(executes func() error, interrupt func(error)) {
	g.actors = append(g.actors, actor{execute: executes, interrupt: interrupt})
}

// Run all actors (functions) concurrently.
// When the first actor returns error, all others will be interrupted.
func (g *Group) Run() error {
	if len(g.actors) == 0 {
		return nil
	}

	errors := make(chan error, len(g.actors))
	for _, a := range g.actors {
		go func(a actor) {
			errors <- a.execute()
		}(a)
	}

	// Wait for the first actor to stop.
	err := <-errors

	// Signal all actors to stop.
	for _, a := range g.actors {
		a.interrupt(err)
	}

	// Wait for all actors to stop.
	for i := 1; i < cap(errors); i++ {
		<-errors
	}

	// Return the original error.
	return err
}
