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

type Metric func(a, b interface{}) float64

type VPTree struct {
	root           *node
	tau            float64
	distanceMetric Metric
}

func New(metric Metric, items []interface{}) (t *VPTree) {
	t = &VPTree{
		distanceMetric: metric,
	}
	t.root = t.buildFromPoints(items)
	return
}

func (vp *VPTree) Search(target interface{}, k int) (results []interface{}, distances []float64) {
	h := make(priorityQueue, 0, k)

	vp.tau = math.MaxFloat64
	vp.search(vp.root, target, k, &h)

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

		n.Threshold = vp.distanceMetric(items[median], n.Item)
		n.Left = vp.buildFromPoints(items[:median])
		n.Right = vp.buildFromPoints(items[median:])
	}
	return
}

func (vp *VPTree) search(n *node, target interface{}, k int, h *priorityQueue) {
	if n == nil {
		return
	}

	dist := vp.distanceMetric(n.Item, target)

	if dist < vp.tau {
		if h.Len() == k {
			heap.Pop(h)
		}
		heap.Push(h, &heapItem{n.Item, dist})
		if h.Len() == k {
			vp.tau = h.Top().(*heapItem).Dist
		}
	}

	if n.Left == nil && n.Right == nil {
		return
	}

	if dist < n.Threshold {
		if dist-vp.tau <= n.Threshold {
			vp.search(n.Left, target, k, h)
		}

		if dist+vp.tau >= n.Threshold {
			vp.search(n.Right, target, k, h)
		}
	} else {
		if dist+vp.tau >= n.Threshold {
			vp.search(n.Right, target, k, h)
		}

		if dist-vp.tau <= n.Threshold {
			vp.search(n.Left, target, k, h)
		}
	}
}
