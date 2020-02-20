package main

import "fmt"

func predicate(i int) bool {
	var modCheck = 2
	fmt.Printf("%s = %+v\n", "modCheck", modCheck)
	return i%modCheck == 0
}

func main() {
	const (
		iterNum = 10
	)
	fmt.Printf("%s = %+v\n", "iterNum", iterNum)

	var (
		message      string = "sum is: "
		code1, code2        = 10, 11
	)
	fmt.Printf("%s = %+v\n", "message", message)
	fmt.Printf("%s = %+v\n", "code1", code1)
	fmt.Printf("%s = %+v\n", "code2", code2)

	sum, delta := 0, 5
	fmt.Printf("%s = %+v\n", "sum", sum)
	fmt.Printf("%s = %+v\n", "delta", delta)
	for i := 0; i < iterNum; i++ {
		if predicate(i) {
			sum = sum + delta
			fmt.Printf("%s = %+v\n", "sum", sum)
		}
	}

	fmt.Printf("%d %d: %s %d\n", code1, code2, message, sum)
}
