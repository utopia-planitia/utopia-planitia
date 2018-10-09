package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	utopia "github.com/utopia-planitia/utopiactl/pkg/utopia"
)

const help = `
usage:
	clusterctl configure [service-selector]
	clusterctl exec [service-selector] [command]
example:
	clusterctl configure all
	clusterctl exec service1,service2 git fetch --all
`

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to determine current working directory: %v", err)
	}

	command := os.Args[1]

	svcs, err := services(cwd, os.Args[2])
	if err != nil {
		log.Fatalf("failed to select services: %v", err)
	}

	if contains([]string{"configure", "reconfigure", "config", "cfg", "c"}, command) {
		err := utopia.Customize(cwd, svcs)
		if err != nil {
			log.Fatalf("failed to auto configure: %v", err)
		}
		return
	}

	if contains([]string{"execute", "exec", "exe", "e"}, command) {
		err := utopia.Exec(cwd, svcs, os.Args[3:])
		if err != nil {
			log.Fatalf("failed to execute: %v", err)
		}
		return
	}

	printHelp()
}

func printHelp() {
	log.Printf(help)
}

func services(directory string, ls string) ([]string, error) {
	if ls != "" && ls != "all" {
		return strings.Split(ls, ","), nil
	}
	services, err := subDirectories(filepath.Join(directory, "services"))
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %v", err)
	}
	return services, nil
}

func subDirectories(path string) ([]string, error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %v", err)
	}

	ls := []string{}
	for _, content := range contents {
		if !content.IsDir() {
			continue
		}
		ls = append(ls, content.Name())
	}
	return ls, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
