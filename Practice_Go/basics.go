package main

import "fmt"

func main() {
	// var name string = "Shwejit"
	// age := 22 --> int used as data type for it
	// var marks float64 = 45.3
	// fmt.Println(age)
	// fmt.Println(name)
	// fmt.Println(marks)

	// var nums []int --> declaring a slice using var , but cannot use scnln to input the value in it beacuse the array size is not defined here , so can use appent to dynamically input the value one by one
	// nums = append(nums, 10)
	// nums = append(nums, 20)
	// fmt.Println(nums)

	// arr := [] int {16,17,18} --> to normally create a slice

	//--> to create a slice using make function, defining the size as well 
	var num int
	fmt.Println("Enter the size of array:")
	fmt.Scanln(&num)

	arr := make([]int, num)

	var n string

	//--> Normal method of using for loop in a slice to iterate through values in the slice
	// for i := 0; i < len(arr); i++ {
	// 	fmt.Print("Enter the age: ", arr[i], " name: ")
	// 	fmt.Scan(&n)
	// 	fmt.Println(driving(arr[i], n))
	// }

	//--> GO method of using for loop for a slice to iterate through values in the slice
	for _, i := range arr {
		fmt.Println("Enter the age: ")
		fmt.Scan(&i)
		fmt.Println("Enter the name: ")
		fmt.Scan(&n)
		fmt.Println(driving(i, n))
	}
}
func driving(age int, n string) string {

	if age < 18 {
		return n + " is not eligible for driving license"
	} else {
		return n + " is eligible for driving license"
	}
}
