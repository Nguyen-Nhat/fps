package converter

func ArrToMapIdentifyKeyInt[T any](arrInput []T, funcGetKey func(T) int) map[int]T {
	return ArrToMapKeyInt(arrInput, funcGetKey, func(t T) T { return t })
}

func ArrToMapKeyInt[T any, M any](arrInput []T, funcGetKey func(T) int, funcGetValue func(T) M) map[int]M {
	mapRes := make(map[int]M, len(arrInput))

	for _, e := range arrInput {
		mapRes[funcGetKey(e)] = funcGetValue(e)
	}

	return mapRes
}
