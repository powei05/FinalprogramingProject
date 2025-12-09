package main

import "math"

func DistanceOrdered(a, b OrderedPair) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	return math.Sqrt(dx*dx + dy*dy)
}

func NormOrdered(p OrderedPair) float64 {
	return math.Sqrt(p.x*p.x + p.y*p.y)
}

func NormalizeOrdered(p OrderedPair) OrderedPair {
	n := NormOrdered(p)
	if n == 0 {
		return OrderedPair{}
	}
	return OrderedPair{x: p.x / n, y: p.y / n}
}

func AddOrdered(a, b OrderedPair) OrderedPair {
	return OrderedPair{x: a.x + b.x, y: a.y + b.y}
}

func SubOrdered(a, b OrderedPair) OrderedPair {
	return OrderedPair{x: a.x - b.x, y: a.y - b.y}
}

func ScaleOrdered(p OrderedPair, k float64) OrderedPair {
	return OrderedPair{x: p.x * k, y: p.y * k}
}

func WrapPosition(p OrderedPair, width float64) OrderedPair {
	x := p.x
	y := p.y
	for x < 0 {
		x += width
	}
	for x >= width {
		x -= width
	}
	for y < 0 {
		y += width
	}
	for y >= width {
		y -= width
	}
	return OrderedPair{x: x, y: y}
}

func ReflectFromCircle(position OrderedPair, center OrderedPair, radius float64) OrderedPair {
	dx := position.x - center.x
	dy := position.y - center.y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		return position
	}
	if dist <= radius {
		scale := (radius + 1.0) / dist
		return OrderedPair{
			x: center.x + dx*scale,
			y: center.y + dy*scale,
		}
	}
	return position
}

func LerpOrdered(a, b OrderedPair, t float64) OrderedPair {
	return OrderedPair{
		x: a.x + (b.x-a.x)*t,
		y: a.y + (b.y-a.y)*t,
	}
}
