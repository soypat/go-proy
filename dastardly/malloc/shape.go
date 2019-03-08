package malloc


type Shape []int


func (s Shape) Size() int {
	size:=1
	for x:=range s {
		size = x*size
	}
	return size
}

func New(nums ...int) (s Shape) {
	for x := range nums {
		s = append(s,x)
}
	return s
}