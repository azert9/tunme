package utils

type Integer interface {
	int | int8 | int16 | int32 | int64
}

type Number interface {
	Integer | float32 | float64
}

func Min[T Number](num ...T) T {

	result := num[0]

	for i := 1; i < len(num); i++ {
		if num[i] < result {
			result = num[i]
		}
	}

	return result
}

func Max[T Number](num ...T) T {

	result := num[0]

	for i := 1; i < len(num); i++ {
		if num[i] > result {
			result = num[i]
		}
	}

	return result
}

func Pow[A Integer, B Integer](n A, e B) A {

	if e < 0 {
		panic("not supported")
	}

	acc := n

	for {
		if e == 0 {
			return 1
		} else if e == 1 {
			return acc
		} else if e%2 == 0 {
			acc *= acc
			e /= 2
		} else {
			acc *= n
			e--
		}
	}
}
