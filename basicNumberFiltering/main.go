package main

import (
	"math"
	"slices"
	"strconv"
	"strings"
)

func returnEven(nums []int) []int {
	ans := make([]int, 0)
	for _, num := range nums {
		if num%2 == 0 {
			ans = append(ans, num)
		}
	}
	return ans
}

func returnOdd(nums []int) []int {
	ans := make([]int, 0)
	for _, num := range nums {
		if num%2 != 0 {
			ans = append(ans, num)
		}
	}
	return ans
}
func checkPrime(num int) bool {
	if num < 2 {
		return false
	}
	for i := 2; float64(i) <= math.Sqrt(float64(num)); i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}
func returnPrime(nums []int) []int {
	ans := make([]int, 0)
	for _, num := range nums {
		if checkPrime(num) {
			ans = append(ans, num)
		}
	}
	return ans
}

func returnMutiples(nums []int, multiple int) []int {
	ans := make([]int, 0)
	for _, num := range nums {
		if num%multiple == 0 {
			ans = append(ans, num)
		}
	}
	return ans
}

func returnGreaterThan(nums []int, greaterInt int) []int {
	ans := make([]int, 0)
	for _, num := range nums {
		if num > greaterInt {
			ans = append(ans, num)
		}
	}
	return ans
}

func returnLessThan(nums []int, lesserInt int) []int {
	ans := make([]int, 0)
	for _, num := range nums {
		if num < lesserInt {
			ans = append(ans, num)
		}
	}
	return ans
}
func conditionEvaluator(expression string, nums []int) []int {
	ans := make([]int, 0)
	if strings.HasPrefix(expression, "multiple of") {
		strMutiple := strings.Split(expression, " ")
		intMutiple, err := strconv.Atoi(strMutiple[len(strMutiple)-1])
		if err == nil {
			ans = returnMutiples(nums, intMutiple)
		}
	} else if strings.HasPrefix(expression, "greater than") {
		strGreater := strings.Split(expression, " ")
		intGreater, err := strconv.Atoi(strGreater[len(strGreater)-1])
		if err == nil {
			ans = returnGreaterThan(nums, intGreater)
		}
	} else if strings.HasPrefix(expression, "less than") {
		strLesser := strings.Split(expression, " ")
		intLesser, err := strconv.Atoi(strLesser[len(strLesser)-1])
		if err == nil {
			ans = returnLessThan(nums, intLesser)
		}
	} else if expression == "even" {
		ans = returnEven(nums)
	} else if expression == "odd" {
		ans = returnOdd(nums)
	} else if expression == "prime" {
		ans = returnPrime(nums)
	}
	return ans
}
func story4(nums []int) []int {
	ans := make([]int, 0)
	nums = returnOdd(nums)
	for _, num := range nums {
		if checkPrime(num) {
			ans = append(ans, num)
		}
	}
	return ans
}
func story5(nums []int) []int {
	ans := returnEven(nums)
	multiple := 5
	ans = returnMutiples(ans, multiple)
	return ans
}

func story6(nums []int) []int {
	ans := returnOdd(nums)
	mutiple := 3
	ans = returnMutiples(ans, mutiple)
	return ans
}
func story7(nums []int, conditions []string) []int {
	for _, condition := range conditions {
		nums = conditionEvaluator(condition, nums)
	}
	return nums
}

func story8(nums []int, conditions []string) []int {
	ans := make([]int, 0)
	for _, condition := range conditions {
		temp := conditionEvaluator(condition, nums)
		for _, t := range temp {
			if !slices.Contains(ans, t) {
				ans = append(ans, t)
			}
		}
	}
	slices.Sort(ans)
	return ans
}
