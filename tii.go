// This software is distributed under the MIT License.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"

	"github.com/fatih/color"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

var (
	scanner = bufio.NewScanner(os.Stdin)
	helpMsg = `Tii - Instantly install command when not found

On most GNU/Linux systems, when a command is not found, a message showing what
can be run to install the command is printed. However, macOS does not have
this. This program adds a similar function with support for macOS using
Homebrew. Instead of simply printing the best matches, Tii shows package 
descriptions and also offers to run an install command for you.

Usage: tii [--help/-h | --version/-v | --refresh-cache/-r | <command>]

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

Files:
   $XDG_DATA_HOME/tii: used to cache package list info. If $XDG_DATA_HOME is
      not set, ~/.local/share is used instead. Refresh the cache using the 
      --refresh-cache option.

If Tii was installed correctly, using commands which are not found will
automatically trigger it. The name Tii is an acronym for "Then Install It".

See source and report bugs at github.com/quackduck/tii`
	//prefix                  = "" // set in init()
	underline               = color.New(color.Underline).SprintFunc()
	disablePrompts          = os.Getenv("TII_DISABLE_INTERACTIVE") == "true" //nolint // complains about using the literal string "true" 3 times
	autoInstallExactMatches = os.Getenv("TII_AUTO_INSTALL_EXACT_MATCHES") == "true"
	version                 = "development" // This will be set at build time using ldflags: go build -ldflags="-s -w -X main.version=$(git describe --tags --abbrev=0)"
	formulaURL              = "https://formulae.brew.sh/api/formula.json"
	caskURL                 = "https://formulae.brew.sh/api/cask.json"
	dataDir                 = os.Getenv("XDG_DATA_HOME")
	dataFile                = "pkginfo.json"
)

type Formula struct {
	Name string `json:"name"`
	//FullName string `json:"full_name"`
	Desc string `json:"desc"`
}

type Cask struct {
	Name string `json:"token"`
	//FullNames []string `json:"name"`
	Desc string `json:"desc"`
}

func main() {
	if dataDir == "" {
		dataDir = os.Getenv("HOME") + "/.local/share"
	}
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
	if hasOption, _ := argsHaveOption("refresh-cache", "r"); hasOption {
		_, _, err := getCachedPackageInfo(true)
		if err != nil {
			handleErr(err)
			return
		}
		fmt.Println("Cache refreshed!")
		return
	}
	if disablePrompts {
		fmt.Println("Running Tii in non-interactive mode. ($TII_DISABLE_INTERACTIVE is true)")
	}
	findPkg(os.Args[1])
}

func findPkg(search string) {
	list, descriptions, err := getCachedPackageInfo(false)
	if err != nil {
		handleErr(err)
		return
	}
	matches := fuzzy.RankFindFold(search, list)
	sort.Sort(matches)

	if len(matches) > 0 {
		if matches[0].Target == search {
			fmt.Println("Found exact match: " + color.YellowString(matches[0].Target) + color.HiBlackString(" ("+descriptions[list[matches[0].OriginalIndex]]+")"))
			if autoInstallExactMatches {
				fmt.Println("Installing it because auto-install is enabled. ($TII_AUTO_INSTALL_EXACT_MATCHES is true)")
				run("brew install " + matches[0].Target)
				return
			}
			runWithPrompt("Install it", "brew install "+matches[0].Target)
			return
		}

		fmt.Println("Presenting fuzzy matches")
		for i, match := range matches {
			if match.Distance > 10 {
				break
			}
			fmt.Println(color.CyanString(strconv.Itoa(i+1)) + ": " + match.Target + color.HiBlackString(" ("+descriptions[list[match.OriginalIndex]]+")"))
			if i == 9 {
				fmt.Println("... and more")
				break
			}
		}
		if ok, i := promptInt("Enter number to install or press enter to quit", 1, len(matches)); ok {
			if runWithPrompt("Install it", "brew install "+matches[i-1].Target) {
				return
			}
		}
	}
	fmt.Println("No exact matches found for " + color.YellowString(search) + ".")
	runWithPrompt("Update Homebrew formulae database", "brew update")
}

func getCachedPackageInfo(forceRefresh bool) ([]string, map[string]string, error) {
	if _, err := os.Stat(dataDir + "/tii/" + dataFile); os.IsNotExist(err) || forceRefresh {
		err = os.MkdirAll(dataDir+"/tii", 0755)
		if err != nil {
			return nil, nil, err
		}
		f, err := os.Create(dataDir + "/tii/" + dataFile)
		if err != nil {
			return nil, nil, err
		}
		defer f.Close()
		list, descriptions, err := fetchPackageInfo()
		if err != nil {
			return nil, nil, err
		}
		err = json.NewEncoder(f).Encode(descriptions) // descriptions also has the full list, so no need to save list
		if err != nil {
			return nil, nil, err
		}
		return list, descriptions, nil
	}
	f, err := os.Open(dataDir + "/tii/" + dataFile)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	var descriptions map[string]string
	err = json.NewDecoder(f).Decode(&descriptions)
	if err != nil {
		return nil, nil, err
	}
	var list []string
	for name := range descriptions {
		list = append(list, name)
	}
	return list, descriptions, nil
}

// returns the list, the map to descriptions, and an error
func fetchPackageInfo() ([]string, map[string]string, error) {
	var list []string
	var formulae []Formula
	var casks []Cask
	var descriptions = make(map[string]string, 1000)

	resp, err := http.Get(formulaURL)
	if err != nil {
		handleErr(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&formulae)
	if err != nil {
		handleErr(err)
	}
	resp, err = http.Get(caskURL)
	if err != nil {
		handleErr(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&casks)
	if err != nil {
		handleErr(err)
	}
	for _, formula := range formulae {
		list = append(list, formula.Name)
		descriptions[formula.Name] = formula.Desc
	}
	for _, cask := range casks {
		list = append(list, cask.Name)
		descriptions[cask.Name] = cask.Desc
	}
	return list, descriptions, nil
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
