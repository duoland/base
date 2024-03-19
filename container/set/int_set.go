package set

// IntSet is not routine-safe
type IntSet struct {
	setElems map[int]struct{}
}

// NewIntSet returns a new int set
func NewIntSet(elems ...int) *IntSet {
	setElems := make(map[int]struct{})
	for _, e := range elems {
		setElems[e] = struct{}{}
	}
	return &IntSet{
		setElems: setElems,
	}
}

// NewIntSetFunc returns a new int set with a decoration
func NewIntSetFunc(decoration func(int) int, elems ...int) *IntSet {
	setElems := make(map[int]struct{})
	for _, e := range elems {
		setElems[decoration(e)] = struct{}{}
	}
	return &IntSet{
		setElems: setElems,
	}
}

// Contains returns true when elem in the set
func (s *IntSet) Contains(elem int) bool {
	if _, exists := s.setElems[elem]; exists {
		return true
	}
	return false
}

// Elems returns the elems of the int set
func (s *IntSet) Elems() []int {
	elems := make([]int, 0, len(s.setElems))
	for k := range s.setElems {
		elems = append(elems, k)
	}
	return elems
}

// ElemsFunc returns the elems that meet the specified func
func (s *IntSet) ElemsFunc(fn func(int) bool) []int {
	elems := make([]int, 0, len(s.setElems))
	for k := range s.setElems {
		if fn(k) {
			elems = append(elems, k)
		}
	}
	return elems
}

// Append appends the elems into the int set and return itself
func (s *IntSet) Append(elems ...int) *IntSet {
	for _, e := range elems {
		s.setElems[e] = struct{}{}
	}
	return s
}

// Union returns a new int set which holds the union elems of both set
func (s *IntSet) Union(p *IntSet) *IntSet {
	unionElems := make(map[int]struct{})
	for k, v := range s.setElems {
		unionElems[k] = v
	}

	for k, v := range p.setElems {
		unionElems[k] = v
	}

	return &IntSet{
		setElems: unionElems,
	}
}

// Intersect returns a new int set which holds the elems intersect of both set
func (s *IntSet) Intersect(p *IntSet) *IntSet {
	intersectElems := make(map[int]struct{})
	for k, v := range s.setElems {
		if _, ok := p.setElems[k]; ok {
			intersectElems[k] = v
		}
	}

	return &IntSet{
		setElems: intersectElems,
	}
}

// Difference returns a new int set which holds the elems in s but not in p
func (s *IntSet) Difference(p *IntSet) *IntSet {
	diffElems := make(map[int]struct{})
	for k, v := range s.setElems {
		if _, ok := p.setElems[k]; !ok {
			diffElems[k] = v
		}
	}
	return &IntSet{
		setElems: diffElems,
	}
}

// Exists checks whether the specified elem exists in set.
func (s *IntSet) Exists(val int) bool {
	if _, ok := s.setElems[val]; ok {
		return true
	}
	return false
}
