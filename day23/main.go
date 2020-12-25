package main

import (
	"container/ring"
	"fmt"
)

const (
	elements   = 1000000
	iterations = 10000000
)

func main() {
	cups := []int{6, 1, 4, 7, 5, 2, 8, 3, 9}
	// cups := []int{3, 8, 9, 1, 2, 5, 4, 6, 7}

	r := ring.New(elements)
	c := r
	for _, cup := range cups {
		c.Value = cup
		c = c.Next()
	}

	for i := 10; i <= elements; i++ {
		c.Value = i
		c = c.Next()
	}

	index := make(map[int]*ring.Ring, elements)
	index[r.Value.(int)] = r
	for c := r.Next(); c != r; c = c.Next() {
		index[c.Value.(int)] = c
	}

	c = r
	for i := 0; i < iterations; i++ {
		cut := c.Unlink(3)
		dest := selectDestinationElement(index, cut, elements, c.Value.(int))
		dest.Link(cut)
		c = c.Next()
	}

	one := index[1]
	fmt.Printf("Lösung: %d\n", one.Next().Value.(int)*one.Move(2).Value.(int))
}

func findElem(r *ring.Ring, value int) *ring.Ring {
	if r.Value.(int) == value {
		return r
	}

	for c := r.Next(); c != r; c = c.Next() {
		if c.Value.(int) == value {
			return c
		}
	}

	return nil
}

func selectDestinationElement(index map[int]*ring.Ring, cut *ring.Ring, overallLen, current int) *ring.Ring {
	for {
		current--
		if current < 1 {
			current = overallLen
		}

		if in(cut, current) {
			continue
		}

		if r, ok := index[current]; ok {
			return r
		}
	}
}

func in(r *ring.Ring, sample interface{}) bool {
	if r.Value == sample {
		return true
	}
	for c := r.Next(); c != r; c = c.Next() {
		if c.Value == sample {
			return true
		}
	}
	return false
}
