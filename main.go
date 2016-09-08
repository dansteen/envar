package main

import "syscall"
import "log"
import "os"
import "os/exec"
import "fmt"
import "regexp"
import "strings"

func main() {
	// check if we have not passed in any commands, or if the only command is -h
	if len(os.Args) < 2 || len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Printf("Usage: %s <command> <command_options>\n\n", os.Args[0])
		os.Exit(0)
	}

	// buld our regex to identify environment variables
	regex, err := regexp.Compile(`\${?[A-Z_][A-Z0-9_]*}?`)
	if err != nil {
		log.Fatal("Could not compile regex")
	}

	// we run through our provided command and replace any environment variables with their values
	for index, variable := range os.Args {
		// get our matches
		matches := regex.FindAllString(variable, -1)
		// replace each match
		for _, match := range matches {
			variable = strings.Replace(variable, match, os.Getenv(strings.Trim(match, "${}")), 1)
		}
		// replace our command line parameter
		os.Args[index] = variable
	}

	// try to find the command
	os.Args[1], err = exec.LookPath(os.Args[1])
	if err != nil {
		log.Fatal("Could not find command in PATH")
	}

	err = syscall.Exec(os.Args[1], os.Args[1:], os.Environ())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal("You find yourself in an incorrect location. Something is desperately wrong...")
}
