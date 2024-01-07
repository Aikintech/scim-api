package facades

import (
	"os"
	"strconv"
)

type env struct {
}

func Env() *env {
	return &env{}
}

// getEnv gets env from os.
func (e *env) getEnv(envName string, defaultValue ...interface{}) interface{} {
	val := os.Getenv(envName)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

// GetString gets string type env from os.
func (e *env) GetString(path string, defaultValue ...interface{}) string {
	return e.getEnv(path, defaultValue...).(string)
}

// GetInt gets int type env from os.
func (e *env) GetInt(path string, defaultValue ...interface{}) int {
	val := e.getEnv(path, defaultValue...)
	switch v := val.(type) {
	case int:
		return v
	case string:
		intVal, err := strconv.Atoi(v)
		if err == nil {
			return intVal
		}
	}
	return 0
}

// GetBool gets bool type config from application.
func (e *env) GetBool(path string, defaultValue ...interface{}) bool {
	val := e.getEnv(path, defaultValue...)
	switch v := val.(type) {
	case bool:
		return v
	case string:
		boolVal, err := strconv.ParseBool(v)
		if err == nil {
			return boolVal
		}
	}
	return false
}
