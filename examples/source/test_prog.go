package main

import "fmt"

func predicate(i int) bool {
	var modCheck = 2 // сокращенный var
	return i%modCheck == 0
}

func main() {
	const (
		iterNum = 10 // блок констант
	)

	var (
		message      string = "sum is: " // длинная форма в блоке
		code1, code2        = 10, 11     // короткая форма в блоке
	)

	sum, delta := 0, 5 // множественное присваивание в короткой форме
	for i := 0; i < iterNum; i++ {
		if predicate(i) {
			sum = sum + delta // переприсваивание
		}
	}

	fmt.Printf("%d %d: %s %d\n", code1, code2, message, sum)
}
