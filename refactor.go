package easyWeb

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func err(e error) {
	if e != nil {
		panic(e)
	}
}
