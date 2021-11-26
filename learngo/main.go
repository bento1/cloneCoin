package main

import (
	"fmt"

	"github.com/cloneCoin/person"
)

func plus(a, b int, name string) (int, string) {
	return a + b, name
}
func plus2(a ...int) int {
	result := 0
	for _, item := range a {
		result += item
	}
	return result
}

func main() {
	result, name := plus(2, 2, "dong")
	fmt.Println(result, name)
	result2 := plus2(2, 3, 4, 5, 6, 7, 8, 9, 10)
	fmt.Println(result2)
	x := 34834945
	fmt.Printf("%d\n", x) //format
	fmt.Printf("%b\n", x) //formating binary
	fmt.Printf("%o\n", x) //formating oxital 8진법
	stringx := fmt.Sprintf("%b", x)
	fmt.Println(stringx)
	foods := [3]string{"potato", "pizza", "pasta"}
	for _, dish := range foods {
		fmt.Println(dish)
	}
	for i := 0; i < len(foods); i++ {
		fmt.Println(foods[i])
	}
	foods2 := []string{"potato", "pizza", "pasta"}
	foods2 = append(foods2, "chicken")
	for _, dish := range foods2 {
		fmt.Println(dish)
	}
	foods3 := append(foods2, "chicken")
	for _, dish := range foods3 {
		fmt.Println(dish)
	}
	a := 2
	b := &a
	a = 12
	fmt.Println(&b, &a, *b)
	dongun := person.Person{}
	fmt.Println(dongun)
	dongun.SetDetails("dongun", 32)
	fmt.Println(dongun)

}
