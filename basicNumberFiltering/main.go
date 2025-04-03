package main

import (
	"math"
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

func conditionTransformer(expression string) func(num int) bool {
	if strings.HasPrefix(expression, "multiple of") {
		strMutiple := strings.Split(expression, " ")
		intMutiple, err := strconv.Atoi(strMutiple[len(strMutiple)-1])
		multipleFn := func(num int) bool {
			return num%intMutiple == 0
		}
		if err == nil {
			return multipleFn
		}
	} else if strings.HasPrefix(expression, "greater than") {
		strGreater := strings.Split(expression, " ")
		intGreater, err := strconv.Atoi(strGreater[len(strGreater)-1])
		greaterFN := func(num int) bool {
			return num > intGreater
		}
		if err == nil {
			return greaterFN
		}
	} else if strings.HasPrefix(expression, "less than") {
		strLesser := strings.Split(expression, " ")
		intLesser, err := strconv.Atoi(strLesser[len(strLesser)-1])
		lessFn := func(num int) bool {
			return num < intLesser
		}
		if err == nil {
			return lessFn
		}
	} else if expression == "even" {
		evenFn := func(num int) bool {
			return num%2 == 0
		}
		return evenFn
	} else if expression == "odd" {
		oddFn := func(num int) bool {
			return num%2 != 0
		}
		return oddFn
	} else if expression == "prime" {
		primeFn := func(num int) bool {
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
		return primeFn
	}
	return func(num int) bool {
		return false
	}
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
	ans := make([]int, 0)
	for _, num := range nums {
		flag := true
		for _, condition := range conditions {
			check := conditionTransformer(condition)
			if !check(num) {
				flag = false
				break
			}
		}
		if flag {
			ans = append(ans, num)
		}
	}

	return ans
}

func story8(nums []int, conditions []string) []int {
	ans := make([]int, 0)
	for _, num := range nums {
		for _, condition := range conditions {
			check := conditionTransformer(condition)
			if check(num) {
				ans = append(ans, num)
				break
			}
		}
	}
	return ans
}
