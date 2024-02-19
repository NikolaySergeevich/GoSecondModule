package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello")
	fmt.Println(mergerTwoArr([4]int{1, 7, 8, 9}, [5]int{3, 4, 5, 6, 6}))
}
//задача первая
func mergerTwoArr(ar1 [4]int, ar2 [5]int) (res [9]int) {
	ind1 := 0
	ind2 := 0
	for i := 0; i < 9; i++ {
		if ind1 == len(ar1) {
			res[i] = ar2[ind2]
			ind2 = ind2 + 1
			continue
		} else if ind2 == len(ar2) {
			res[i] = ar1[ind1]
			ind1 = ind1 + 1
			continue
		}
		if ar1[ind1] == ar2[ind2] {
			res[i] = ar1[ind1]
			res[i+1] = ar1[ind1]
			i = i + 1
			ind1 = ind1 + 1
			ind2 = ind2 + 1
		} else if ar1[ind1] > ar2[ind2] {
			res[i] = ar2[ind2]
			ind2 = ind2 + 1
		} else {
			res[i] = ar1[ind1]
			ind1 = ind1 + 1
		}
	}
	return
}
