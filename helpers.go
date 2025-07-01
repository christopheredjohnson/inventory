package main

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sign(n int) int {
	if n < 0 {
		return -1
	} else if n > 0 {
		return 1
	}
	return 0
}
