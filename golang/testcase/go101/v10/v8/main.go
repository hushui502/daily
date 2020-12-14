package main

type Point struct {
	x int
	y int
}

func (p Point) Add(q Point) Point {
	return Point{x: p.x + q.x, y: p.y + q.y}
}

func (p Point) Sub(q Point) Point {
	return Point{x: p.x - q.x, y: p.y - q.y}
}

type Path []Point

func (path Path) TranslateBy(offset Point, add bool) {
	var op func(p, q Point) Point
	if add {
		op = Point.Add
	} else {
		op = Point.Sub
	}
	for i := range path {
		path[i] = op(path[i], offset)
	}
}
