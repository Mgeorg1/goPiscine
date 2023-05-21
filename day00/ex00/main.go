package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func calcMode(nums []int) int {
	m := make(map[int]int)
	max := -100000
	for i := range nums {
		if _, ok := m[nums[i]]; ok == true {
			m[nums[i]]++
		} else {
			m[nums[i]] = 1
		}
		if max < nums[i] {
			max = nums[i]
		}
	}
	return max
}

func calcSD(nums []int, mean float64) float64 {
	s := float64(0)
	var distance float64
	for i := range nums {
		distance = mean - float64(nums[i])
		s += distance * distance
	}
	return math.Sqrt(distance / float64(len(nums)))
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	var nums []int
	var sum, i float64 = 0, 0
	for in.Scan() {
		txt := in.Text()
		if txt == "" {
			fmt.Println("Empty input")
			os.Exit(1)
		}
		num, err := strconv.Atoi(txt)
		if err != nil {
			fmt.Println(err.Error)
			os.Exit(1)
		}
		if num < -100000 || num > 100000 {
			fmt.Println("Out of range")
			os.Exit(1)
		}
		nums = append(nums, num)
		sum += float64(num)
		i++
	}
	mean := sum / float64(len(nums))
	fmt.Printf("Mean: %f\n", mean)
	var median float64
	if len(nums)%2 != 0 {
		middle := len(nums) / 2
		median = float64(nums[middle])
	} else {
		middle := len(nums) / 2
		median = float64(nums[middle]+nums[middle-1]) / 2
	}
	fmt.Printf("Median: %f\n", median)
	mode := calcMode(nums)
	fmt.Printf("Mode: %d\n", mode)
	SD := calcSD(nums, mean)
	fmt.Printf("SD: %f\n", SD)
}
