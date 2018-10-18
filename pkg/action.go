package utopia

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ExecuteCommandline parses and runs a command in directory cwd
func ExecuteCommandline(cwd string, args []string) error {

	if len(args) < 3 {
		return fmt.Errorf("to few arguments")
	}

	command := args[1]

	svcs, err := services(cwd, args[2])
	if err != nil {
		return fmt.Errorf("failed to select services: %v", err)
	}

	if contains([]string{"configure", "reconfigure", "config", "conf", "cfg", "c"}, command) {
		err := Configure(cwd, svcs)
		if err != nil {
			return fmt.Errorf("failed to auto configure: %v", err)
		}
		return nil
	}

	if contains([]string{"deploy"}, command) {
		err := Deploy(cwd, svcs)
		if err != nil {
			return fmt.Errorf("failed to deploy: %v", err)
		}
		return nil
	}

	if len(args) < 4 {
		return fmt.Errorf("to few arguments")
	}

	if contains([]string{"execute", "exec", "exe", "e"}, command) {
		err := Exec(cwd, svcs, args[3:])
		if err != nil {
			return fmt.Errorf("failed to execute: %v", err)
		}
		return nil
	}

	return fmt.Errorf("command unknown")
}

func services(directory string, ls string) ([]string, error) {
	if ls == "-" {
		return []string{}, nil
	}
	if ls != "all" {
		return strings.Split(ls, ","), nil
	}
	services, err := subDirectories(filepath.Join(directory, "services"))
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %v", err)
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("could not find services")
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