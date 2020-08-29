package server

import (
	"strconv"
)

func itoa(i int) string {
	return strconv.Itoa(i)
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		i = 0
	}
	return i
}
