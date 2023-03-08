package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func GenMsgIncr(userID string) string {
	return userID + "_" + Int64ToString(GetCurrentTimestampByNano())
}

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}

// Get the current timestamp by Nano
func GetCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}
