package main

import "fmt"

func main() {
	sum := 0
	alt := false
	number := "4000056655665556"
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if alt {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		sum += digit
		alt = !alt
	}
	fmt.Println(sum%10 == 0)
}
