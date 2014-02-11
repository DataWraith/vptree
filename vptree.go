package vptree

import (
	"container/heap"
	"math"
	"math/rand"
)

type node struct {
	Index     int
	Threshold float64
	Left      *node
	Right     *node
}

type heapItem struct {
	Index int
	Dist  float64
}

type Metric func(a, b interface{}) float64

type VPTree struct {
	root           *node
	tau            float64
	items          []interface{}
	distanceMetric Metric
}

func New(metric Metric, items []interface{}) (t *VPTree) {
	t = &VPTree{
		items:          items,
		distanceMetric: metric,
	}
	if len(items) > 0 {
		t.root = t.buildFromPoints(0, len(items)-1)
	}
	return
}

func (vp *VPTree) Search(target interface{}, k int) (results []interface{}, distances []float64) {
	h := make(priorityQueue, 0, k)

	vp.tau = math.MaxFloat64
	vp.search(vp.root, target, k, &h)

	for h.Len() > 0 {
		hi := h.Pop()
		results = append(results, vp.items[hi.(*heapItem).Index])
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

func (vp *VPTree) buildFromPoints(lower int, upper int) (n *node) {
	if upper == lower {
		return nil
	}

	n = &node{Index: lower}

	if upper-lower > 1 {
		// Choose an arbitrary point and move it to the start
		i := rand.Intn(upper-lower-1) + lower
		vp.items[lower], vp.items[i] = vp.items[i], vp.items[lower]

		median := (upper + lower) / 2

		// Partition around the median distance
		pivotDistance := vp.distanceMetric(vp.items[median], vp.items[lower])
		left := lower + 1
		right := upper
		for left < right {
			for vp.distanceMetric(vp.items[left], vp.items[lower]) < pivotDistance {
				left += 1
			}
			for vp.distanceMetric(vp.items[right], vp.items[lower]) > pivotDistance {
				right -= 1
			}
			if left <= right {
				vp.items[left], vp.items[right] = vp.items[right], vp.items[left]
				left += 1
				right -= 1
			}
		}

		// What was the median?
		n.Threshold = vp.distanceMetric(vp.items[lower], vp.items[median])

		n.Index = lower
		n.Left = vp.buildFromPoints(lower+1, median)
		n.Right = vp.buildFromPoints(median, upper)
	}

	return
}

func (vp *VPTree) search(n *node, target interface{}, k int, h *priorityQueue) {
	if n == nil {
		return
	}

	dist := vp.distanceMetric(vp.items[n.Index], target)

	if dist < vp.tau {
		if h.Len() == k {
			heap.Pop(h)
		}
		heap.Push(h, &heapItem{n.Index, dist})
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
