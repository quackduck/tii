package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	scanner = bufio.NewScanner(os.Stdin)
	helpMsg = `tii - Directly install command when not found

On most GNU/Linux systems, when a command is not found, a message showing what 
can be run to install the command is printed. However, macOS does not 
have this. This program supports a similar function with support for macOS 
(only for macOS, as of now). Instead of simply printing the command, tii also 
offers to run it for you. 

Usage: tii [--help/-h | <command>]
Examples: 
   tii fish
   tii cowsay
   tii --help

If tii was installed correctly, using commands which are not found will 
automatically trigger it. The name tii is an acronym for "Then Install It".`
	formulaeLocation = "/usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula"
)

func main() {
	if _, err := exec.LookPath("brew"); err != nil {
		fmt.Println("Homebrew is not installed. Install it to use tii")
		runWithPrompt("Install Homebrew", `/bin/bash -c "\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
		return
	}
	if len(os.Args) > 2 {
		handleErrStr("error: too many arguments")
		fmt.Println(helpMsg)
		return
	}
	if hasOption, _ := argsHaveOption("help", "h"); hasOption || len(os.Args) == 1 {
		fmt.Println(helpMsg)
		return
	}
	findPkg(os.Args[1])
}

// takes in a string to search for a package
func findPkg(search string) {
	file, err := os.Open(formulaeLocation)
	if err != nil {
		handleErrStr("Could not open " + formulaeLocation)
		handleErr(err)
		return
	}
	defer file.Close()
	list, err := file.Readdirnames(0) // >=0 to read all files and folders
	_ = file.Close()                  // close file right after reading so it's run even if ^C or os.Exit() used later in the function
	if err != nil {
		handleErrStr("An error occurred while trying to list files in " + formulaeLocation)
		handleErr(err)
	}
	possibleMatches := make([]string, 0, 5)
	gotExactMatch := false
	for _, name := range list {
		formulaName := strings.TrimSuffix(name, filepath.Ext(name))
		if formulaName == search {
			fmt.Println("Found exact match")
			gotExactMatch = true
			runWithPrompt("Run", "brew install "+formulaName)
			break
		} else if strings.Contains(formulaName, search) {
			possibleMatches = append(possibleMatches, formulaName)
		}
	}
	if len(possibleMatches) > 0 {
		fmt.Println("Presenting possible matches [" + strconv.Itoa(len(possibleMatches)) + "]")
		for i, name := range possibleMatches {
			fmt.Println(strconv.Itoa(i+1) + ": " + color.GreenString(name))
		}
		if ok, i := promptInt("Enter number to install or press enter to quit: ", 1, len(possibleMatches)); ok {
			runWithPrompt("Run", "brew install "+possibleMatches[i+1])
		}
	}
	if !gotExactMatch {
		fmt.Println("No exact matches found for " + color.YellowString(search) + ".")
		runWithPrompt("Update Homebrew formulae database, "brew update")
	}
}

func promptBool(promptStr string) (yes bool) {
	for {
		fmt.Print(promptStr + " (y/N) > ")
		if !scanner.Scan() {
			break
		}
		switch scanner.Text() {
		case "y", "Y", "yes", "Yes", "YES", "true", "True", "TRUE":
			return true
		case "", "n", "N", "no", "No", "NO", "false", "False", "FALSE":
			return false
		default:
			continue
		}
	}
	return true
}

// quits if user enters enter
func promptInt(promptStr string, lowerLimit int, upperLimit int) (bool, int) {
	for {
		fmt.Print(promptStr)
		if !scanner.Scan() {
			break
		}
		if scanner.Text() == "" {
			break
		}
		if i, err := strconv.Atoi(scanner.Text()); err == nil && lowerLimit <= i && i <= upperLimit {
			return true, i
		}
	}
	return false, 0
}

func runWithPrompt(promptStr string, command string) (ran bool) {
	yes := promptBool(promptStr + " with " + "`" + command + "`" + "?")
	if yes {
		// run it with the users shell
		cmd := exec.Command(os.Getenv("SHELL"), "-c", command) //nolint //"Subprocess launched with function call as argument or cmd arguments"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			handleErrStr("An error occurred while trying to run " + command)
			handleErr(err)
			return false
		}
	}
	return yes
}

func argsHaveOption(long string, short string) (hasOption bool, foundAt int) {
	for i, arg := range os.Args {
		if arg == "--"+long || arg == "-"+short {
			return true, i
		}
	}
	return false, 0
}

func handleErr(err error) {
  fmt.Fprintln(os.Stderr, color.Red("error: "+err.Error()))
}

func handleErrStr(str string) {
  fmt.Fprintln(os.Stderr, color.Red(str))
}
