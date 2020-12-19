package sort

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	arr := []int{1,5,9,6,3,7,5,10}
	fmt.Println("排序前：",arr)
	bubbleSort(arr,len(arr))
	fmt.Println("排序后：",arr)
}

func TestInsertionSort(t *testing.T) {
	arr := []int{1,5,9,6,3,7,5,10}
	fmt.Println("排序前：",arr)
	insertSort(arr,len(arr))
	fmt.Println("排序后：",arr)
}

func TestSelectionSort(t *testing.T) {
	arr := []int{1,5,9,6,3,7,5,10}
	fmt.Println("排序前：",arr)
	selectSort(arr,len(arr))
	fmt.Println("排序后：",arr)
}

func TestMergeSort(t *testing.T) {
	arr := []int{5, 4}
	MergeSort(arr)
	t.Log(arr)

	arr = []int{5, 4, 3, 2, 1}
	MergeSort(arr)
	t.Log(arr)
}

func createRandomArr(length int) []int {
	arr := make([]int, length, length)
	for i := 0; i < length; i++ {
		arr[i] = rand.Intn(100)
	}
	return arr
}

func TestQuickSort(t *testing.T) {
	arr := []int{5, 4}
	QuickSort(arr)
	t.Log(arr)

	arr = createRandomArr(100)
	QuickSort(arr)
	t.Log(arr)
}
