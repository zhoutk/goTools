package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {
	_, err := IntFromInt64(math.MaxInt32 + 1)
	if err != nil {
		fmt.Println(err)
	}
}

func ConvertInt64ToInt(i64 int64) int {
	if math.MinInt32 <= i64 && i64 <= math.MaxInt32 {
		return int(i64)
	}
	panic("can't convert int64 to int")
}

func IntFromInt64(i64 int64) (i int, err error) {//这里
	defer func() {
		if err2 := recover(); err2 != nil {
			i = 0//这里
			err = errors.New("ttt")//这里
		}
	}()
	i = ConvertInt64ToInt(i64)
	return i, nil
}
