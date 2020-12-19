package heap

type Heap struct {
	a []int
	n int
	count int
}

func NewHeap(capacity int) *Heap {
	return &Heap{
		a:     make([]int, capacity),
		n:     capacity,
		count: 0,
	}
}

func (heap *Heap) insert(data int) {
	if heap.count == heap.n {
		return
	}

	heap.count++
	heap.a[heap.count] = data

	// compare with parent node
	i := heap.count
	parent := i / 2
	for parent > 0 && heap.a[parent] < heap.a[i] {
		swap(heap.a, parent, i)
		i = parent
		parent = parent / 2
	}
}

func (heap *Heap) removeMax() {
	if heap.count == 0 {
		return
	}

	swap(heap.a, 1, heap.count)
	heap.count--

	heapifyUpToDown(heap.a, heap.count)
}

func heapifyUpToDown(a []int, count int) {
	for i := 1; i <= count/2; {
		maxIndex := i
		if a[i] < a[i*2] {
			maxIndex = i*2
		}
		if i*2+1 <= count && a[maxIndex] < a[i*2+1] {
			maxIndex = i*2 + 1
		}

		if maxIndex == 1 {
			break
		}

		swap(a, i, maxIndex)
		i = maxIndex
	}
}

func swap(a []int, i, j int) {
	a[i], a[j] = a[j], a[i]
}