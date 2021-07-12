package main

func main() {
	limit := NewConcurrencyLimiter(10)
	for i := 0; i < 1000; i++ {
		limit.Execute(func() {
			// do some work
		})
	}
	limit.Wait()
}
