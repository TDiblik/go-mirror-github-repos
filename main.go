package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	GH_TOKEN := loadEnvRequired("GH_TOKEN")
	MIRROR_PATH := loadEnvRequired("MIRROR_PATH")
	EXCLUDED_ORGS, EXCLUDED_ORGS_RAW := loadEnvList("EXCLUDED_ORGS")
	EXCLUDED_REPOSITORIES, EXCLUDED_REPOSITORIES_RAW := loadEnvList("EXCLUDED_REPOSITORIES")
	log.Println("Running with the following environment variables:")
	log.Println("	MIRROR_PATH: \"" + MIRROR_PATH + "\"")
	log.Println("	EXCLUDED_ORGS: \"" + EXCLUDED_ORGS_RAW + "\"")
	log.Println("	EXCLUDED_REPOSITORIES_ENV: \"" + EXCLUDED_REPOSITORIES_RAW + "\"")

	// Get list of repositories
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	if err != nil {
		log.Fatalf("error creating /user/repos http request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+GH_TOKEN)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error calling /user/repos http request: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading /user/repos http body: %s\n", err)
		os.Exit(1)
	}

	// Parse the list
	var repos []Repository
	err = json.Unmarshal(data, &repos)
	if err != nil {
		log.Fatalf("error unmarshalling /user/repos JSON data: %s\n", err)
		os.Exit(1)
	}
	log.Println("Successfuly called and parsed /user/repos")

	// Clone/Update each repo
	for _, repo := range repos {
		if slices.Contains(EXCLUDED_ORGS, repo.Owner.Login) {
			log.Println("skipping " + repo.FullName + ", because it's inside the EXCLUDED_ORGS env variable")
			continue
		}
		if slices.Contains(EXCLUDED_REPOSITORIES, repo.FullName) {
			log.Println("skipping " + repo.FullName + ", because it's inside the EXCLUDED_REPOSITORIES env variable")
			continue
		}

		log.Println("running for " + repo.FullName)
		command := "cd \"" + MIRROR_PATH + "\"; "

		owner_directory := filepath.Join(MIRROR_PATH, repo.Owner.Login)
		owner_directory_exists, err := doesPathExists(owner_directory)
		if err != nil {
			log.Fatalf("unable to find out whether path (%s) exists: %s", owner_directory, err)
			os.Exit(1)
		}
		if !owner_directory_exists {
			log.Println("will freshly create: ", repo.Owner.Login)
			command += "mkdir \"" + repo.Owner.Login + "\"; "
		}
		command += "cd \"" + repo.Owner.Login + "\"; "

		repo_directory := filepath.Join(MIRROR_PATH, repo.Name)
		repo_directory_exists, err := doesPathExists(repo_directory)
		if err != nil {
			log.Fatalf("unable to find out whether path (%s) exists: %s", repo_directory, err)
			os.Exit(1)
		}
		if !repo_directory_exists {
			log.Println("will freshly clone: ", repo.FullName)
			command += "git clone https://" + GH_TOKEN + "@github.com/" + repo.FullName + ".git; "
		} else {
			log.Println("will only mirror changes to: ", repo.FullName)
		}
		command += "cd \"" + repo.Name + "\"; "

		command += `
			base_branch=$(git branch --show-current);

			git fetch --all;
			for branch in $(git branch -r | grep -v '\->'); do
				git branch --track "${branch#origin/}" "$branch" || true
			done

			for branch in $(git for-each-ref --format='%(refname:short)' refs/heads/); do
				git checkout "$branch"
				git reset --hard "origin/$branch"
			done;
			
			git checkout "$base_branch"; 
		`

		log.Println("started work...")
		stdout, err := exec.Command("bash", "-c", command).Output()
		if err != nil {
			log.Fatalf("Error running the command (%s): %s\n", command, err)
			os.Exit(1)
		}
		log.Println("Result of the commands executed: \n ----------------- \n", string(stdout), "-----------------")
	}
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

func doesPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
