package driver

import (
	"strconv"
	"time"
)

func parseDuration(d string) time.Duration {
	duration, err := time.ParseDuration(d)
	if err != nil {
		panic(err)
	}
	return duration
}

func parseInt(i string) int {
	res, err := strconv.ParseInt(i, 10, 0)
	if err != nil {
		panic(err)
	}
	return int(res)
}

func readWithDefault(m map[string]string, key string, def string) string {
	value, ok := m[key]
	if ok {
		return value
	}

	return def
}
