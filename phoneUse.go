package main

import (
	"fmt"

	"github.com/ttacon/libphonenumber"
)

func main() {
	num, err := libphonenumber.Parse("+8613671827770", "")
	fmt.Println(num, err)
	fmt.Println(*num.CountryCode, *num.NationalNumber)

	num, err = libphonenumber.Parse("13671827770", "CN")
	fmt.Println(num, err)

	num, err = libphonenumber.Parse("+13648889954", "")
	fmt.Println(num, err)

	num, err = libphonenumber.Parse("+27872406111", "")
	fmt.Println(num, err)

	formattedNum := libphonenumber.Format(num, libphonenumber.NATIONAL)
	fmt.Println(formattedNum)

	formattedNum = libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
	fmt.Println(formattedNum)

	formattedNum = libphonenumber.Format(num, libphonenumber.E164)
	fmt.Println(formattedNum)
}
