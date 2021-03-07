Bug-risks are issue in code that can cause errors and breakages in production.
A bug is flaw in the code that produces undesired or incorrect results.
Code often has bug-risks due to poor coding practices, lack of version control,
miscommunication of requirements, unrealistic time schedules for development,
and buggy third-party tools,in this post, let's take a look at few commonly
seen bug-risks in Go.

- Infinite recursive call 

- Assignment to nil map

```bigquery
make(map[string]int)
make(map[string]int, 100)
```

Bad Pattern:
```bigquery
var countedData map[string][]ChartElement
```

Good Pattern:
```bigquery
countedData := make(map[string][]ChartElement)
```

- Method modifies receiver
For example:
```bigquery
type data struct {
    num   int
    key   *string
    items map[string]bool
}

func (d data) vmethod() {
    d.num = 8
}

func (d data) run() {
    d.vmethod()
    fmt.Printf("%+v", d) // Output: {num:1 key:0xc0000961e0 items:map[1:true]}
}
```
if num must be modified:
```bigquery
type data struct {
	num   int
	key   *string
	items map[string]bool
}

func (d *data) vmethod() {
	d.num = 8
}

func (d *data) run() {
	d.vmethod()
	fmt.Printf("%+v", d) // Output: &{num:8 key:0xc00010a040 items:map[1:true]}
}
```

- Possibly undesired value being used in goroutine
  In the example below, the value of index and value used in the goroutine are from the outer scope. 
  Because the goroutines run asynchronously, the value of index and value could be (and usually are) different from the intended value
```bigquery
mySlice := []string{"A", "B", "C"}
for index, value := range mySlice {
	go func() {
		fmt.Printf("Index: %d\n", index)
		fmt.Printf("Value: %s\n", value)
	}()
}
```
To overcome this problem, a local scope must be created, like in the example below.
```bigquery
mySlice := []string{"A", "B", "C"}
for index, value := range mySlice {
	index := index
	value := value
	go func() {
		fmt.Printf("Index: %d\n", index)
		fmt.Printf("Value: %s\n", value)
	}()
}
```
Another way to handle this could be by passing the values as args to the goroutines.
```bigquery
mySlice := []string{"A", "B", "C"}
for index, value := range mySlice {
	go func(index int, value string) {
		fmt.Printf("Index: %d\n", index)
		fmt.Printf("Value: %s\n", value)
	}(index, value)
}
```

- Deferring Close before checking for a possible error
```bigquery
f, err := os.Open("/tmp/file.md")
if err != nil {
    return err
}
defer f.Close()
```
But this pattern is harmful for writable files because deferring a function call ignores its return value, 
and the Close() method can return errors. 
For instance, if you wrote data to the file, it might have been cached in memory and not flushed to disk by the time you called Close. 
This error should be explicitly handled.

```bigquery
f, err := os.Open("/tmp/file.md")
if err != nil {
	return err
}

defer func() {
	closeErr := f.Close()
	if closeErr != nil {
		if err == nil {
			err = closeErr
		} else {
			log.Println("Error occured while closing the file :", closeErr)
		}
	}
}()
return err
```