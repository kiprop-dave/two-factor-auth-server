package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateUserId() string {
	rand.Seed(time.Now().UnixNano())
	userId := "S"
	for i := 0; i < 4; i++ {
		num := rand.Intn(10)
		userId = userId + strconv.Itoa(num)
	}
	return userId
}
