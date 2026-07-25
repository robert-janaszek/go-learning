package jsonfixer

func pop(slice *[]byte) {
	length := len(*slice)
	*slice = (*slice)[:length-1]
}
