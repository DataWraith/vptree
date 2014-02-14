# vptree

vptree is a port of Steve Hanov's C++
[implementation](http://stevehanov.ca/blog/index.php/?id=130) of [Vantage-point
trees](https://en.wikipedia.org/wiki/Vantage-point_tree) to the Go programming
language. Vantage-point trees are useful for nearest-neighbour searches in
high-dimensional metric spaces.


## Installation

	go get github.com/DataWraith/vptree


## Usage

First, you need to define the metric space you want to search in:

```go
import "fmt"
import "github.com/DataWraith/vptree"

type Coordinate struct {
	X float64
	Y float64
}

func CoordinateMetric(a, b interface{}) float64 {
	c1 := a.(Coordinate)
	c2 := b.(Coordinate)

	return math.Sqrt(math.Pow(c1.X-c2.X, 2) + math.Pow(c1.Y-c2.Y, 2))
}
```

Coordinate is the user-defined type that you want to search nearest neighbours
for. CoordinateMetric is a function that defines the distance between two
Coordinates, in this case the Euclidean Distance. Note that leaving out the
square-root operation will sabotage this, since Squared Euclidean Distance
is not a metric. A metric in the mathematical sense is required for VPTree to
operate correctly; a metric `d` must have the following properties:

* d(x, y) >= 0
* d(x, y) = 0 if and only if x = y
* d(x, y) = d(y, x)
* d(x, z) <= d(x, y) + d(y, z) (triangle inequality)

The next step is to build the tree:

```go
func NearestNeighbor() {
	// Define some coordinates
	coordinates := []Coordinate{
		Coordinate{24, 57},
		Coordinate{35, 28},
		Coordinate{55, 48},
		Coordinate{68, 42}
	}

	// Convert the slice of coordinates into a slice of interface{}
	vpitems := make([]interface{}, len(coordinates))
	for i, v := range coordinates {
		vpitems[i] = interface{}(v)
	}

	// Build the tree
	tree := vptree.New(CoordinateMetric, vpitems)
```

Now you can search the tree for the k nearest neighbours of a query point.

```go
	// Define the query point
	q := Coordinate{12, 34}

	// Number of neighbours to return
	k := 3

	// Get Coordinates and distances of the k nearest neighbours
	neighbours, distances := tree.Search(q, k)

	fmt.Println(neighbours)
	// Output: [{35 28} {24 57} {55 48}]

	fmt.Println(distances)
	// Output: [23.769728648009426 25.942243542145693 45.221676218380054]
}
```
