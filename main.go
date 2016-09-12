package main

import "syscall"
import "log"
import "os"

import "os/exec"
import "fmt"
import "strings"

func main() {
	// check if we have not passed in any commands, or if the only command is -h
	if len(os.Args) < 2 || len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Printf("Usage: %s <command> <command_options>\n\n", os.Args[0])
		os.Exit(0)
	}

	// we process environment variables as well for environment variable references
	for _, value := range os.Environ() {
		// split the value into parts
		varParts := strings.SplitN(value, "=", 2)
		// run through the value and resolve the variables present down as far as possible
		valuePart := os.ExpandEnv(varParts[1])
		for oldValuePart := ""; oldValuePart != valuePart; valuePart = os.ExpandEnv(valuePart) {
			oldValuePart = valuePart
		}
		// once we are done processing we set the environment variable
		err := os.Setenv(varParts[0], valuePart)
		if err != nil {
			log.Fatal("Could not set env variable %s=%s", varParts[0], valuePart)
		}
	}

	// we run through our provided command and replace any environment variables with their values
	for index, arg := range os.Args {
		os.Args[index] = os.ExpandEnv(arg)
	}

	// try to find the command (we break this into two lines because of http://stackoverflow.com/questions/21345274/go-fails-to-infer-type-in-assignment-non-name-on-left-side-of )
	command, err := exec.LookPath(os.Args[1])
	os.Args[1] = command
	if err != nil {
		log.Fatal("Could not find command in PATH")
	}

	err = syscall.Exec(os.Args[1], os.Args[1:], os.Environ())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal("You find yourself in an incorrect location. Something is desperately wrong...")
}
