package vptree

import (
	"container/heap"
	"math"
	"math/rand"
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

	// Reverse results and distances, because we popped them from the heap
	// in large-to-small order
	for i, j := 0, len(coords)-1; i < j; i, j = i+1, j-1 {
		coords[i], coords[j] = coords[j], coords[i]
		distances[i], distances[j] = distances[j], distances[i]
	}

	return
}

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

	if len(coords1) != len(coords2) {
		t.Fatalf("Expected %v coordinates, got %v", len(coords2), len(coords1))
	}

	if len(distances1) != len(distances2) {
		t.Fatalf("Expected %v distances, got %v", len(distances2), len(distances1))
	}

	for i := 0; i < len(coords1); i++ {
		if coords1[i] != coords2[i] {
			t.Errorf("Expected coords1[%v] to be %v, got %v", i, coords2[i], coords1[i])
		}
		if distances1[i] != distances2[i] {
			t.Errorf("Expected distances1[%v] to be %v, got %v", i, distances2[i], distances1[i])
		}
	}
}

// This test creates a bunch of random input items and tests against a linear search
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

	if len(coords1) != len(coords2) {
		t.Fatalf("Expected %v coordinates, got %v", len(coords2), len(coords1))
	}

	if len(distances1) != len(distances2) {
		t.Fatalf("Expected %v distances, got %v", len(distances2), len(distances1))
	}

	for i := 0; i < len(coords1); i++ {
		if coords1[i] != coords2[i] {
			t.Errorf("Expected coords1[%v] to be %v, got %v", i, coords2[i], coords1[i])
		}
		if distances1[i] != distances2[i] {
			t.Errorf("Expected distances1[%v] to be %v, got %v", i, distances2[i], distances1[i])
		}
	}
}
