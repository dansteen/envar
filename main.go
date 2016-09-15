package main

import "syscall"
import "log"
import "os"

import "os/exec"
import "fmt"
import "strings"

// shellSplit will split a string in the way that a shell would expect, respecting quoted characters
// culled from here: https://gist.github.com/jmervine/d88c75329f98e09f5c87
func shellSplit(s string) []string {
	split := strings.Split(s, " ")

	var result []string
	var inquote string
	var block string
	for _, i := range split {
		if inquote == "" {
			if strings.HasPrefix(i, "'") || strings.HasPrefix(i, "\"") {
				inquote = string(i[0])
				block = strings.TrimPrefix(i, inquote) + " "
			} else {
				result = append(result, i)
			}
		} else {
			if !strings.HasSuffix(i, inquote) {
				block += i + " "
			} else {
				block += strings.TrimSuffix(i, inquote)
				inquote = ""
				result = append(result, block)
				block = ""
			}
		}
	}
	return result
}

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
	newArgs := make([]string, 0)
	for _, arg := range os.Args {
		// the replaced variables may contain multiple switches/arguments.  We need to break them out into their own elements
		argParts := shellSplit(os.ExpandEnv(arg))
		// add to our new set of args
		newArgs = append(newArgs, argParts...)
	}
	os.Args = newArgs

	// look for our command.
	command, err := exec.LookPath(os.Args[1])
	if err != nil {
		log.Fatalf("Could not find %s in PATH", os.Args[1])
	}
	os.Args[1] = command

	err = syscall.Exec(os.Args[1], os.Args[1:], os.Environ())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal("You find yourself in an incorrect location. Something is desperately wrong...")
}
