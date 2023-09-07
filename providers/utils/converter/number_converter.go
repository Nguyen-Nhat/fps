package converter

func IntArrToInt32Arr(arr []int) []int32 {
	if len(arr) == 0 {
		return []int32{}
	}

	var output []int32
	for _, item := range arr {
		output = append(output, int32(item))
	}

	return output
}
