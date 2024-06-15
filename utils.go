package main

import (
	"log"
	"os"
	"strings"
)

// Includes only fields used, the API returns MUCH more
type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    Owner  `json:"owner"`
}
type Owner struct {
	Login string `json:"login"`
}

func loadEnvRequired(env_name string) string {
	env_value := os.Getenv(env_name)
	if len(env_value) == 0 {
		log.Fatal(env_name + " env variable cannot be empty")
		os.Exit(1)
	}
	return env_value
}

func loadEnvList(env_name string) ([]string, string) {
	env_value_parsed := make([]string, 0)
	env_value_raw := os.Getenv(env_name)
	if len(env_value_raw) != 0 {
		env_value_parsed = strings.Split(env_value_raw, ",")

		for i, value := range env_value_parsed {
			env_value_parsed[i] = strings.TrimSpace(value)
		}
	}

	return env_value_parsed, env_value_raw
}

func pathExistsDef(path string) bool {
	exists, err := pathExists(path)
	if err != nil {
		log.Fatalf("unable to find out whether path (%s) exists: %s", path, err)
		os.Exit(1)
	}
	return exists
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
