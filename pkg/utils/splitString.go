package utils

type String struct {
	Line string
}

func (s String) Split(char string) [2]string {
	var out [2]string
	first := ""
	sec := ""
	isFirst := true
	for _, c := range s.Line {
		if string(c) == char {
			isFirst = false
			continue
		}
		if isFirst {
			first = first + string(c)
		} else {
			sec = sec + string(c)
		}
	}
	out[0] = first
	out[1] = sec
	return out
}
