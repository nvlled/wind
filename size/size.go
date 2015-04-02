package size

import (
	"fmt"
)

type T interface {
	Size()
	String() string
	Equals(s T) bool
	LessThan(s T) bool
	Value(alloc int) int
	Add(s T) T
}

type ConstT int
type RangeT struct{ min, max int }
type FreeT struct{}

type Folder func(sizes []T) T
type Allocator func(x int, sizes []T) []int

func Const(x int) ConstT    { return ConstT(x) }
func Range(x, y int) RangeT { return RangeT{x, y} }

var Free = FreeT{}

func (_ ConstT) Size() {}
func (_ RangeT) Size() {}
func (_ FreeT) Size()  {}

func (s ConstT) String() string { return fmt.Sprintf("ConstT(%d)", int(s)) }
func (s RangeT) String() string {
	return fmt.Sprintf("RangeT(%d, %d)", s.min, s.max)
}
func (s FreeT) String() string { return "FreeT" }

func (c ConstT) Equals(s T) bool {
	switch t := s.(type) {
	case ConstT:
		return int(c) == int(t)
	}
	return false
}

func (r RangeT) Equals(s T) bool {
	switch t := s.(type) {
	case RangeT:
		return r.min == t.min && r.max == t.max
	}
	return false
}

func (r FreeT) Equals(s T) bool {
	switch s.(type) {
	case FreeT:
		return true
	}
	return false
}

func (c ConstT) LessThan(s T) bool {
	x := int(c)
	switch t := s.(type) {
	case ConstT:
		return x < int(t)
	case RangeT:
		return x < t.Length()
	}
	return true
}

func (r RangeT) LessThan(s T) bool {
	length := r.Length()
	switch t := s.(type) {
	case ConstT:
		return length < int(t)
	case RangeT:
		return length < t.Length()
	}
	return true
}

func (f FreeT) LessThan(s T) bool {
	return false
}

func (s ConstT) Value(alloc int) int { return lower(alloc, int(s)) }
func (s RangeT) Value(alloc int) int { return lower(alloc, s.Length()) }
func (s FreeT) Value(alloc int) int  { return alloc }

func (e FreeT) Add(s T) T {
	return e
}

func (c ConstT) Add(s T) T {
	x := int(c)
	switch v := s.(type) {
	case FreeT:
		return v.Add(c)
	case ConstT:
		return Const(x + int(v))
	case RangeT:
		return reduct(Range(x+v.min, x+v.max))
	}
	panic("non-exhaustive case analysis")
}

func (r RangeT) Add(s T) T {
	switch v := s.(type) {
	case FreeT:
		return v.Add(r)
	case ConstT:
		return v.Add(r)
	case RangeT:
		return reduct(Range(r.min+v.min, r.max+v.max))
	}
	return s.Add(r)
}

func Sum(sizes []T) T {
	var total T = ConstT(0)
	for _, s := range sizes {
		total = total.Add(s)
	}
	return total
}

func (r RangeT) Length() int {
	// return r.max
	return zero(r.max - r.min + 1)
}

func Max(sizes []T) T {
	var x T = Const(0)
	for _, s := range sizes {
		if x.LessThan(s) {
			x = s
		}
	}
	return x
}

func Int(n int) T {
	if n < 0 {
		return Free
	}
	return Const(n)
}

// reduction rules:
//  range(n, m) = 0 where n > m
//  range(n, n) = const(n)
func reduct(s T) T {
	switch v := s.(type) {
	case RangeT:
		if v.min == v.max {
			return Const(v.min)
		}
		if v.min > v.max {
			return Const(0)
		}
	}
	return s
}
