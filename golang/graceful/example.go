package graceful


/*
func main() {
	m := newManager()

	m.AddRunningJob(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Println("working job-1")
				time.Sleep(time.Second)
			}
		}
	})

	m.AddRunningJob(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Println("working job-2")
				time.Sleep(time.Second)
			}
		}
	})

	// Add shutdown 01
	m.AddShutdownJob(func() error {
		log.Println("shutdown job 01 and wait 1 second")
		time.Sleep(1 * time.Second)
		return nil
	})

	// Add shutdown 02
	m.AddShutdownJob(func() error {
		log.Println("shutdown job 02 and wait 2 second")
		time.Sleep(2 * time.Second)
		return nil
	})

	<-m.Done()
}
*/