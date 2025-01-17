package main

import "fmt"
import "day05/ex02/presents"

func main() {
	pres := []presents.Present{
		{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
	}

	coolestPresents, err := presents.GetNCoolestPresents(pres, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, present := range coolestPresents {
		fmt.Printf("Value: %d, Size: %d\n", present.Value, present.Size)
	}
}
