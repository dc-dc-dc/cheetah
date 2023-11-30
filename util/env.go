package util

import "os"

var (
	truthyVals = NewSetFromItr([]string{"1", "true", "True", "TRUE", "t", "T", "y", "Y", "yes", "Yes", "YES", "on", "On", "ON"}...)
)

func GetEnv(env, def string) string {
	if val, ok := os.LookupEnv(env); ok {
		return val
	}
	return def
}

func IsTruthy(env string) bool {
	return truthyVals.Contains(os.Getenv(env))
}
