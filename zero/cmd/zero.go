package cmd

import (
	"fmt"
	"strconv"
)

func Add(first string, second string) (result string) {
	num1, err := strconv.ParseFloat(first, 64)
	if err != nil {
		fmt.Print("num1 is not provided")
		return ""
	}
	num2, err := strconv.ParseFloat(second, 64)
	if err != nil {
		fmt.Print("num2 is not provided")
		return ""
	}

	res := num1 + num2
	return fmt.Sprintf("%f", res)
}

func Subtract(from string, subTract string) (result string) {
	num1, err := strconv.ParseFloat(from, 64)
	if err != nil {
		fmt.Println("Error: first value is invalid")
		return
	}
	num2, err := strconv.ParseFloat(from, 64)
	if err != nil {
		fmt.Println("Error: second value is invalid")
		return
	}
	res := num1 - num2
	return fmt.Sprintf("%f", res)
}

func Multiply(first string, second string, shouldRoundUp bool) (result string) {
	num1, err := strconv.ParseFloat(first, 64)
	if err != nil {
		fmt.Println("Error: First value is not a decimal")
		return
	}
	num2, err := strconv.ParseFloat(second, 64)
	if err != nil {
		fmt.Println("Error: Second value is not a decimal")
		return
	}
	if shouldRoundUp {
		return fmt.Sprintf("%.2f", num1*num2)
	}
	return fmt.Sprintf("%f", num1*num2)
}

func Divide(divide string, by string, shouldRoundUp bool) (e error, result string) {
	num1, err := strconv.ParseFloat(divide, 64)
	if err != nil {
		return fmt.Errorf("first value is not a number"), ""
	}
	num2, err := strconv.ParseFloat(by, 64)
	if err != nil {
		return fmt.Errorf("second value is not a number"), ""
	}
	if shouldRoundUp {
		return nil, fmt.Sprintf("%.2f", num1/num2)
	}
	return nil, fmt.Sprintf("%f", num1/num2)
}
