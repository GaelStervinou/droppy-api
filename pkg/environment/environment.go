package environment

import (
	"slices"
)

var DEV_ENV = "dev"
var PROD_ENV = "prod"
var TEST_ENV = "test"

var currentEnv = DEV_ENV

func SetEnv(env string) {
	if slices.Contains([]string{DEV_ENV, PROD_ENV, TEST_ENV}, env) {
		currentEnv = env
		return
	}
	currentEnv = DEV_ENV
}

func IsDev() bool {
	return currentEnv == DEV_ENV
}

func IsProd() bool {
	return currentEnv == PROD_ENV
}

func IsTest() bool {
	return currentEnv == TEST_ENV
}

func GetEnv() string {
	return currentEnv
}
