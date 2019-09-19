package main

import (
	"fmt"
	"strconv"
	"strings"
)

var COUPON_SUBSTITUTE_MAP = map[int]string{
	5: "sku_v2_14_3",
}

func main() {
	split()
}

func split() {
	s := "abc"
	parts := strings.SplitN(s, "=", 2)
	fmt.Printf("parts = %+v\n", parts)

	c := strings.Split(s, "")[2]
	fmt.Println(c)

	fmt.Println(getSkuTypeByCouponSubstitute(5))
}

func getSkuTypeByCouponSubstitute(substitute int) int {
	val, _ := strconv.Atoi(strings.Split(COUPON_SUBSTITUTE_MAP[substitute], "_")[3])
	return val
}
