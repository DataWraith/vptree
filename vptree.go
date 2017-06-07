package vptree

import (
	"container/heap"
	"math"
	"math/rand"
)

type node struct {
	Item      interface{}
	Threshold float64
	Left      *node
	Right     *node
}

type heapItem struct {
	Item interface{}
	Dist float64
}

// A Metric is a function that measures the distance between two provided
// interface{}-values. The function *must* be a metric in the mathematical
// sense, that is, the metric d must fullfill the following requirements:
//
//	* d(x, y) >= 0
//	* d(x, y) = 0 if and only if x = y
//	* d(x, y) = d(y, x)
//	* d(x, z) <= d(x, y) + d(y, z) (triangle inequality)
type Metric func(a, b interface{}) float64

// A VPTree struct represents a Vantage-point tree. Vantage-point trees are
// useful for nearest-neighbour searches in high-dimensional metric spaces.
type VPTree struct {
	root           *node
	distanceMetric Metric
}

// New creates a new VP-tree using the metric and items provided. The metric
// measures the distance between two items, so that the VP-tree can find the
// nearest neighbour(s) of a target item.
func New(metric Metric, items []interface{}) (t *VPTree) {
	t = &VPTree{
		distanceMetric: metric,
	}
	t.root = t.buildFromPoints(items)
	return
}

// Search searches the VP-tree for the k nearest neighbours of target. It
// returns the up to k narest neighbours and the corresponding distances in
// order of least distance to largest distance.
func (vp *VPTree) Search(target interface{}, k int) (results []interface{}, distances []float64) {
	if k < 1 {
		return
	}

	h := make(priorityQueue, 0, k)

	tau := math.MaxFloat64
	vp.search(vp.root, &tau, target, k, &h)

	for h.Len() > 0 {
		hi := heap.Pop(&h)
		results = append(results, hi.(*heapItem).Item)
		distances = append(distances, hi.(*heapItem).Dist)
	}

	// Reverse results and distances, because we popped them from the heap
	// in large-to-small order
	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
		distances[i], distances[j] = distances[j], distances[i]
	}

	return
}

func (vp *VPTree) buildFromPoints(items []interface{}) (n *node) {
	if len(items) == 0 {
		return nil
	}

	n = &node{}

	// Take a random item out of the items slice and make it this node's item
	idx := rand.Intn(len(items))
	n.Item = items[idx]
	items[idx], items = items[len(items)-1], items[:len(items)-1]

	if len(items) > 0 {
		// Now partition the items into two equal-sized sets, one
		// closer to the node's item than the median, and one farther
		// away.
		median := len(items) / 2
		pivotDist := vp.distanceMetric(items[median], n.Item)
		items[median], items[len(items)-1] = items[len(items)-1], items[median]

		storeIndex := 0
		for i := 0; i < len(items)-1; i++ {
			if vp.distanceMetric(items[i], n.Item) <= pivotDist {
				items[storeIndex], items[i] = items[i], items[storeIndex]
				storeIndex++
			}
		}
		items[len(items)-1], items[storeIndex] = items[storeIndex], items[len(items)-1]
		median = storeIndex

		n.Threshold = pivotDist
		n.Left = vp.buildFromPoints(items[:median])
		n.Right = vp.buildFromPoints(items[median:])
	}
	return
}

func (vp *VPTree) search(n *node, tau *float64, target interface{}, k int, h *priorityQueue) {
	if n == nil {
		return
	}

	dist := vp.distanceMetric(n.Item, target)

	if dist < *tau {
		if h.Len() == k {
			heap.Pop(h)
		}
		heap.Push(h, &heapItem{n.Item, dist})
		if h.Len() == k {
			*tau = h.Top().(*heapItem).Dist
		}
	}

	if n.Left == nil && n.Right == nil {
		return
	}

	if dist < n.Threshold {
		if dist-*tau <= n.Threshold {
			vp.search(n.Left, tau, target, k, h)
		}

		if dist+*tau >= n.Threshold {
			vp.search(n.Right, tau, target, k, h)
		}
	} else {
		if dist+*tau >= n.Threshold {
			vp.search(n.Right, tau, target, k, h)
		}

		if dist-*tau <= n.Threshold {
			vp.search(n.Left, tau, target, k, h)
		}
	}
}
