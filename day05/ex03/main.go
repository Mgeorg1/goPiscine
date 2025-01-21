package main

import "day05/ex02/presents"

func getItems(arr [][]int, pres []presents.Present) []presents.Present {
	i := len(arr) - 1
	j := len(arr[0]) - 1
	ret := make([]presents.Present, 0)
	for i > 0 {
		if arr[i][j] != arr[i-1][j] {
			ret = append(ret, pres[i-1])
			j -= pres[i-1].Size
		}
		i--
	}
	return ret
}

func grabPresents(pres []presents.Present, capacity int) []presents.Present {
	arr := make([][]int, len(pres)+1)
	for i := range arr {
		arr[i] = make([]int, capacity+1)
	}
	arr[0][0] = 0
	for i := 1; i <= len(pres); i++ {
		for j := 0; j <= capacity; j++ {
			if pres[i-1].Size > j {
				arr[i][j] = arr[i-1][j]
			} else {
				arr[i][j] = max(arr[i-1][j], arr[i-1][j-pres[i-1].Size]+pres[i-1].Value)
			}
		}
	}
	return getItems(arr, pres)
}

func main() {
	pres := []presents.Present{
		{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
	}

	selectedPresents := grabPresents(pres, 3)
	for _, present := range selectedPresents {
		println(present.Value, present.Size)
	}
}
