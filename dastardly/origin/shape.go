package origin


type Shape []int

var scalarShape = Shape{}

func ScalarShape() Shape {
	return scalarShape
}

func (s Shape) Size() int {
	size:=1
	for x:=range s {
		size = x*size
	}
	return size
}

func NewShape(nums ...int) (s Shape) {
	for x := range nums {
		s = append(s,x)
}
	return s
}