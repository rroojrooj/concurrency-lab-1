package main

import (
	"log"
	"os"
	"runtime/trace"
	"sync"
)

// merge takes two sorted sub-arrays from slice and sorts them.
// The resulting array is put back in slice.
// merge takes two sorted sub-arrays from slice and merges them into a single sorted array.
func merge(slice []int32, middle int) {

	// Create a copy of the original slice to work with while merging
	sliceClone := make([]int32, len(slice))
	copy(sliceClone, slice)

	// Split the slice into two sub-arrays a and b using the middle index
	a := sliceClone[middle:] // Sub-array from middle to the end
	b := sliceClone[:middle] // Sub-array from the start to middle

	// Initialize two pointers, i and j, to traverse through a and b respectively
	i := 0
	j := 0

	// Loop through the entire slice to populate it with merged values
	for k := 0; k < len(slice); k++ {
		// If all elements from a have been merged, populate the slice with the remaining elements from b
		if i >= len(a) {
			slice[k] = b[j]
			j++
			// If all elements from b have been merged, populate the slice with the remaining elements from a
		} else if j >= len(b) {
			slice[k] = a[i]
			i++
			// Compare the current elements of a and b; If the element from a is greater, take the element from b
		} else if a[i] > b[j] {
			slice[k] = b[j]
			j++
			// If the element from b is greater or equal, take the element from a
		} else {
			slice[k] = a[i]
			i++
		}
	}
}

// Sequential merge sort.
func mergeSort(slice []int32) {
	if len(slice) > 1 {
		middle := len(slice) / 2
		mergeSort(slice[:middle])
		mergeSort(slice[middle:])
		merge(slice, middle)
	}
}

// TODO: Parallel merge sort.
func parallelMergeSort(slice []int32) {
	mergeSort(slice)
	if len(slice) <= 1 {
		return // Base case: A slice of length 1 or less is always sorted.
	}

	middle := len(slice) / 2

	// Use sync.WaitGroup to wait for both halves to be sorted.
	var wg sync.WaitGroup

	// Sort the left half in a new goroutine.
	wg.Add(1) // Increment the counter.
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine completes.
		parallelMergeSort(slice[:middle])
	}()

	// Sort the right half in another new goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		parallelMergeSort(slice[middle:])
	}()

	// Wait for both goroutines to complete.
	wg.Wait()

	// Merge the two halves.
	merge(slice, middle)
}

// main starts tracing and in parallel sorts a small slice.
func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	slice := make([]int32, 0, 100)
	for i := int32(100); i > 0; i-- {
		slice = append(slice, i)
	}

	parallelMergeSort(slice)
}
