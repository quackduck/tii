// This software is distributed under the MIT License.
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
	helpMsg = `Tii - Instantly install command when not found

On most GNU/Linux systems, when a command is not found, a message showing what
can be run to install the command is printed. However, macOS does not have
this. This program adds a similar function with support for macOS (only for
macOS, as of now). Instead of simply printing the command, Tii also offers to
run it for you.

Usage: tii [--help/-h | --version/-v | <command>]

Examples:
   tii fish
   tii cowsay
   tii --help

Environment:
	TII_DISABLE_INTERACTIVE: If this variable is set to "true", Tii will
		disable interactive output (prompting for confirmation) and not install
		any packages.
	TII_AUTO_INSTALL_EXACT_MATCHES: If this variable is set to "true", Tii will
		automatically install exact matches without prompting for confirmation

If Tii was installed correctly, using commands which are not found will
automatically trigger it. The name Tii is an acronym for "Then Install It".`
	formulaeLocation        = "/usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula"
	underline               = color.New(color.Underline).SprintFunc()
	disablePrompts          = os.Getenv("TII_DISABLE_INTERACTIVE") == "true" //nolint // complains about using the literal string "true" 3 times
	autoInstallExactMatches = os.Getenv("TII_AUTO_INSTALL_EXACT_MATCHES") == "true"
	version                 = "development" // This will be set at build time using ldflags: go build -ldflags="-s -w -X main.version=$(git describe --tags --abbrev=0)"
)

func main() {
	if _, err := exec.LookPath("brew"); err != nil {
		handleErrStr("Homebrew is not installed. Install it to use Tii")
		runWithPrompt("Install Homebrew", `/bin/bash -c "\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
		return
	}
	if len(os.Args) > 2 {
		handleErrStr("Too many arguments")
		fmt.Println(helpMsg)
		return
	}
	if hasOption, _ := argsHaveOption("help", "h"); hasOption || len(os.Args) == 1 {
		fmt.Println(helpMsg)
		return
	}
	if hasOption, _ := argsHaveOption("version", "v"); hasOption {
		fmt.Println("Tii " + version)
		return
	}
	if disablePrompts {
		fmt.Println("Running Tii in non-interactive mode. ($TII_DISABLE_INTERACTIVE is true)")
	}
	findPkg(os.Args[1])
}

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
			fmt.Println("Found exact match: " + color.YellowString(formulaName))
			gotExactMatch = true
			if autoInstallExactMatches {
				fmt.Println("Installing it because auto-install is enabled. ($TII_AUTO_INSTALL_EXACT_MATCHES is true)")
				run("brew install " + formulaName)
			}
			if runWithPrompt("Install it", "brew install "+formulaName) {
				return
			}
			break
		} else if strings.Contains(formulaName, search) || strings.Contains(search, formulaName) {
			possibleMatches = append(possibleMatches, formulaName)
		}
	}
	if len(possibleMatches) > 0 {
		fmt.Println("Presenting possible matches [" + color.CyanString(strconv.Itoa(len(possibleMatches))) + "]")
		for i, name := range possibleMatches {
			fmt.Println(color.CyanString(strconv.Itoa(i+1)) + ": " + name)
		}
		if ok, i := promptInt("Enter number to install or press enter to quit", 1, len(possibleMatches)); ok {
			runWithPrompt("Install it", "brew install "+possibleMatches[i-1])
		}
	}
	if !gotExactMatch {
		fmt.Println("No exact matches found for " + color.YellowString(search) + ".")
		runWithPrompt("Update Homebrew formulae database", "brew update")
	}
}

func promptBool(promptStr string) (yes bool) {
	if disablePrompts {
		return false
	}
	for {
		fmt.Print(underline(promptStr) + " (y/N) > ")
		color.Set(color.FgCyan)
		if !scanner.Scan() {
			break
		}
		color.Unset()
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
	if disablePrompts {
		return false, 0
	}
	for {
		fmt.Print(underline(promptStr) + ": ")
		color.Set(color.FgCyan)
		if !scanner.Scan() {
			break
		}
		color.Unset()
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
	yes := promptBool(promptStr + " with " + color.YellowString(command) + "?")
	if yes {
		run(command)
	}
	return yes
}

func run(command string) {
	// run it with the users shell
	cmd := exec.Command(os.Getenv("SHELL"), "-c", command) //nolint //"Subprocess launched with function call as argument or cmd arguments"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		handleErrStr("An error occurred while trying to run " + command)
		handleErr(err)
	}
	cmd.Stderr = nil
	cmd.Stdout = nil
	cmd.Stdin = nil
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
	handleErrStr(err.Error())
}

func handleErrStr(str string) {
	_, _ = fmt.Fprintln(os.Stderr, color.RedString("Error: ")+str)
}
