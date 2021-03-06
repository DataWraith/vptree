package vptree

import (
	"container/heap"
	"math"
	"math/rand"
	"sync"
	"testing"
)

type Coordinate struct {
	X float64
	Y float64
}

func CoordinateMetric(a, b interface{}) float64 {
	c1 := a.(Coordinate)
	c2 := b.(Coordinate)

	return math.Sqrt(math.Pow(c1.X-c2.X, 2) + math.Pow(c1.Y-c2.Y, 2))
}

// This helper function compares two sets of coordinates/distances to make sure
// they are the same.
func compareCoordDistSets(t *testing.T, actualCoords []interface{}, expectedCoords []Coordinate, actualDists, expectedDists []float64) {
	if len(actualCoords) != len(expectedCoords) {
		t.Fatalf("Expected %v coordinates, got %v", len(expectedCoords), len(actualCoords))
	}

	if len(actualDists) != len(expectedDists) {
		t.Fatalf("Expected %v distances, got %v", len(expectedDists), len(actualDists))
	}

	for i := 0; i < len(actualCoords); i++ {
		if actualCoords[i] != expectedCoords[i] {
			t.Errorf("Expected actualCoords[%v] to be %v, got %v", i, expectedCoords[i], actualCoords[i])
		}
		if actualDists[i] != expectedDists[i] {
			t.Errorf("Expected actualDists[%v] to be %v, got %v", i, expectedDists[i], actualDists[i])
		}
	}
}

// This helper function finds the k nearest neighbours of target in items. It's
// slower than the VPTree, but its correctness is easy to verify, so we can
// test the VPTree against it.
func nearestNeighbours(target Coordinate, items []Coordinate, k int) (coords []Coordinate, distances []float64) {
	pq := &priorityQueue{}

	// Push all items onto a heap
	for _, v := range items {
		heap.Push(pq, &heapItem{v, CoordinateMetric(v, target)})
	}

	// Pop all but the k smallest items
	for pq.Len() > k {
		heap.Pop(pq)
	}

	// Extract the k smallest items and distances
	for pq.Len() > 0 {
		hi := heap.Pop(pq)
		coords = append(coords, hi.(*heapItem).Item.(Coordinate))
		distances = append(distances, hi.(*heapItem).Dist)
	}

	// Reverse coords and distances, because we popped them from the heap
	// in large-to-small order
	for i, j := 0, len(coords)-1; i < j; i, j = i+1, j-1 {
		coords[i], coords[j] = coords[j], coords[i]
		distances[i], distances[j] = distances[j], distances[i]
	}

	return
}

// This test makes sure vptree's behavior is sane with no input items
func TestEmpty(t *testing.T) {
	vp := New(CoordinateMetric, nil)
	qp := Coordinate{0, 0}

	coords, distances := vp.Search(qp, 3)

	if len(coords) != 0 {
		t.Error("coords should have been of length 0")
	}

	if len(distances) != 0 {
		t.Error("distances should have been of length 0")
	}
}

// This test creates a small VPTree and makes sure its search function returns
// the right results
func TestSmall(t *testing.T) {
	items := []Coordinate{
		Coordinate{24, 57},
		Coordinate{35, 28},
		Coordinate{55, 48},
		Coordinate{68, 42},
	}

	target := Coordinate{12, 34}

	vpitems := make([]interface{}, len(items))
	for i, v := range items {
		vpitems[i] = interface{}(v)
	}

	vp := New(CoordinateMetric, vpitems)
	coords1, distances1 := vp.Search(target, 3)
	coords2, distances2 := nearestNeighbours(target, items, 3)

	compareCoordDistSets(t, coords1, coords2, distances1, distances2)
}

// This test creates a bunch of random input items and tests against the
// simpler, but slower nearestNeighbours function
func TestRandom(t *testing.T) {
	items := make([]Coordinate, 0, 10)

	// Generate 1000 random coordinates
	for i := 0; i < 1000; i++ {
		items = append(items, Coordinate{X: rand.Float64(), Y: rand.Float64()})
	}

	// Build a VPTree
	vpitems := make([]interface{}, len(items))
	for i, v := range items {
		vpitems[i] = interface{}(v)
	}
	vp := New(CoordinateMetric, vpitems)

	// Random query point
	q := Coordinate{X: rand.Float64(), Y: rand.Float64()}

	// Select number of nearest neighbours
	k := rand.Intn(100) + 1

	// Get the k nearest neighbours and their distances
	coords1, distances1 := vp.Search(q, k)
	coords2, distances2 := nearestNeighbours(q, items, k)

	compareCoordDistSets(t, coords1, coords2, distances1, distances2)
}

// This test creates a random tree and tests concurrent queries
func TestConcurrent(t *testing.T) {
	var items []Coordinate

	// Generate 1000 random coordinates
	for i := 0; i < 1000; i++ {
		items = append(items, Coordinate{X: rand.Float64(), Y: rand.Float64()})
	}

	// Build a VPTree
	vpitems := make([]interface{}, len(items))
	for i, v := range items {
		vpitems[i] = interface{}(v)
	}
	vp := New(CoordinateMetric, vpitems)

	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {

		wg.Add(1)

		go func() {
			for j := 0; j < 100; j++ {
				// Random query point
				q := Coordinate{X: rand.Float64(), Y: rand.Float64()}

				// Get the k nearest neighbours and their distances
				coords1, distances1 := vp.Search(q, 10)
				coords2, distances2 := nearestNeighbours(q, items, 10)

				compareCoordDistSets(t, coords1, coords2, distances1, distances2)
			}
			wg.Done()
		}()

	}

	wg.Wait()
}
