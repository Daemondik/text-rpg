package helpers

import (
	"os"
	"strconv"
	"strings"
)

func Env(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func EnvAsInt(name string, defaultVal int) int {
	valueStr := Env(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func EnvAsBool(name string, defaultVal bool) bool {
	valStr := Env(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func EnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := Env(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
