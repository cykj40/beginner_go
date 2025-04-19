package main

import (
	"fmt"
)

func main() {
	// variables

	var name string = "Cyrus"
	fmt.Printf("This is my name %s\n", name)

	age := 43
	fmt.Printf("this is my age %d\n", age)

	var city string = "Holmdel"
	fmt.Printf("this is my city %s\n", city)

	var country, continent string = "USA", "North America"
	fmt.Printf("this is my country %s and this is my continent %s\n", country, continent)

	var (
		isEmployed bool   = true
		salary     int    = 50000
		position   string = "Software Engineer"
	)

	fmt.Printf("isEmployed: %t this is my salary %d this is my position: %s", isEmployed, salary, position)

	// zero values

	var defaultInt int
	var defaultFloat float64
	var defaultString string
	var defaultBool bool

	fmt.Printf("int: %d float: %f string: %s bool: %t\n", defaultInt, defaultFloat, defaultString, defaultBool)

	// constants
	const pi = 3.14
	const (
		Monday    = 1
		Tuesday   = 2
		Wednesday = 3
	)
	fmt.Printf("Monday %d - Tuesday %d Wednesday %d\n", Monday, Tuesday, Wednesday)

	const typedAge int = 44
	const untypedAge = 44

	fmt.Printf("typedAge: %d untypedAge: %d\n", typedAge, untypedAge)

	const (
		Jan int = iota + 1
		Feb
		Mar
		Apr
	)

	fmt.Printf("jan - %d feb - %d mar - %d apr - %d\n", Jan, Feb, Mar, Apr)

	result := add(3, 4)
	fmt.Println("this is the result", result)

	sum, product := calculateSumAndProduct(3, 4)
	fmt.Printf("this is sum: %d and this is product: %d\n", sum, product)

}

func add(a int, b int) int {
	return a + b

}

func calculateSumAndProduct(a, b int) (int, int) {
	return a + b, a * b
}
