package malloc

type Dense struct {
	Shape
	body

	who
	old *Dense // Guardar version transpuesta
}





type body struct{
	header *[]byte // Aqui van los elementos de la matriz
}

type who struct {
	symmetric bool
	zeros bool
	triangularU bool
	triangularL bool
	diagonal bool
	isTransposed bool
}


func malloc(shape Shape) *[]byte {
	header := make([]byte, shape.Size()*8 )
	return &header
}

