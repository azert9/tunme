package utils

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must1[T1 any](arg1 T1, err error) T1 {
	if err != nil {
		panic(err)
	}
	return arg1
}
