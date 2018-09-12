package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/utopia-planitia/utopia-planitia/pkg"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to determine current working directory: %v", err)
	}

	repos, err := utopia.Repositories(cwd, os.Args[1:])
	if err != nil {
		log.Fatalf("failed to setup config: %v", err)
	}

	customizePath := filepath.Join(cwd, "customize")

	for _, repo := range repos {

		log.Println(repo)

		if repo == "customize" {
			continue
		}

		repoPath := filepath.Join(cwd, repo)

		err := utopia.Vars(repoPath, customizePath)
		if err != nil {
			log.Fatalf("failed to copy vars from %v: %v", customizePath, err)
		}

		filepath.Walk(repoPath, utopia.Walk(customizePath, repo, cwd))

		// prerender for example certificates
	}
}
