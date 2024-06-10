package main

import (
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Package struct {
	FULLPKGNAME string `json:"FULLPKGNAME"`
	COMMENT     string `json:"COMMENT"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: opkg [search, install, delete, info] <args>")
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "search":
		if len(args) == 0 {
			fmt.Println("Usage: opkg search <search_query>")
			return
		}
		query := strings.Join(args, " ")
		packages, err := fSearch(query)
		if err != nil {
			fmt.Printf("Error fetching packages: %v\n", err)
			return
		}

		for _, pkg := range packages {
			fmt.Printf("%s - %s\n", pkg.FULLPKGNAME, pkg.COMMENT)
		}

	case "info":
		_ = runCommand("pkg_info", args)

	case "delete":
		_ = runCommand("pkg_delete", args)

	case "install":
		_ = runCommand("pkg_add", args)

	default:
		fmt.Println("Unknown command. Available commands: search, info, delete, install")
	}
}

func fSearch(query string) ([]Package, error) {
	url := fmt.Sprintf("https://openbsd.app/?search=%s&format=json", query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var packages []Package
	err = json.Unmarshal(body, &packages)
	if err != nil {
		return nil, err
	}

	return packages, nil
}

func runCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
