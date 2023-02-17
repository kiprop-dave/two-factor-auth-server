package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenerateOtp() string {
	rand.Seed(time.Now().UnixNano())
	otpText := ""
	for i := 0; i < 6; i++ {
		num := rand.Intn(10)
		otpText = otpText + strconv.Itoa(num)
	}
	fmt.Println(otpText)
	return otpText
}
