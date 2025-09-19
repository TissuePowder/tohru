package env

import (
	"os"
	"strconv"
)

func GetStr(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return v
}

func GetInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return i
}

func GetBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		panic(err)
	}
	return b
}
