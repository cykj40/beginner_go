package main

import (
	"fmt"
)

func main() {
	// variables

	age := 30

	if age >= 18 {
		fmt.Println("You are a adult")
	} else if age >= 13 {
		fmt.Println("You are a teenager")
	} else {
		fmt.Println("You are a child")
	}

	day := "Tuesday"

	switch day {
	case "Monday":
		fmt.Println("start of the week")
	case "Tuesday", "Wednesday", "Thursday":
		fmt.Println("Midweek")
	case "Friday":
		fmt.Println("TGIF")
	default:
		fmt.Println("it's the weekend")
	}

	for i := 0; i < 5; i++ {
		fmt.Println("this is i", i)
	}

	counter := 0
	for counter < 3 {
		fmt.Println("counter is", counter)
		counter++
	}

	iterations := 0
	for {
		if iterations > 3 {
			break
		}
		fmt.Println("iterations is", iterations)
		iterations++
	}

	matrix := [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}

	fmt.Printf("this is our matrix %v\n", matrix)

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Println("matrix[", i, "][", j, "] =", matrix[i][j])
		}
	}

}

func add(a int, b int) int {
	return a + b

}

func calculateSumAndProduct(a, b int) (int, int) {
	return a + b, a * b
}
