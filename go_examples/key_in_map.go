package main

import "fmt"

func main() {

	fmt.Println("\n##############################\n")
	strDict := map[string]int{"japan": 1, "china": 2, "canada": 3}
	value, ok := strDict["china"]
	if ok {
		fmt.Println("Key found value is: ", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("\n##############################\n")
	if value, exist := strDict["china"]; exist {
		fmt.Println("Key found value is: ", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("\n##############################\n")
	intMap := map[int]string{
		0: "zero",
		1: "one",
	}
	fmt.Printf("Key 0 exists: %t\nKey 1 exists: %t\nKey 2 exists: %t",
		intMap[0] != "", intMap[1] != "", intMap[2] != "")

	fmt.Println("\n##############################\n")
	t := map[int]string{}
	if len(t) == 0 {
		fmt.Println("\nEmpty Map")
	}
	if len(intMap) == 0 {
		fmt.Println("\nEmpty Map")
	} else {
		fmt.Println("\nNot Empty Map")
	}

}
