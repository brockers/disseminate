package main

import "os"
import "flag"
import "time"
import "os/exec"
import "fmt"
import "regexp"
import "strings"
import "strconv"
import "encoding/json"
import "io/ioutil"

// Regular Expression filters
var authTag = regexp.MustCompile(`Author: .*\n`)
var dateTag = regexp.MustCompile(`Date:\s+(.*)\n`)
var mer1Tag = regexp.MustCompile(`Merge pull request .*\n`)
var mer2Tag = regexp.MustCompile(`Merge: .*\n`)
var commTag = regexp.MustCompile(`commit ([a-z0-9]*)\n`)

type response struct {
	Commit  string `json:"commit"`
	Message string `json:"message"`
	Date    string `json:"timestamp"`
	Time    string `json:"buildtime"`
	Disseminate PackageDisseminate `json:"package"`
}

type PackageDisseminate struct {
	Product string `json:"product"`
	Website string `json:"website"`
}

type PackageJSON struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	License     string `json:"license"`
	Private     bool   `json:"private"`
	Disseminate PackageDisseminate `json:"disseminate"`
}

// Standarize our error message
func check(e error, s string) {
	if e != nil {
		fmt.Println(s)
		panic(e)
	}
}

// No panic but we want a exit with a message
func warn(s string){
	fmt.Println("ERROR:", s)
	os.Exit(1)
}

// Get the package file and convert it from json
func getPackage(f string) PackageJSON {
	raw, err := ioutil.ReadFile(f)
	check(err,"Unable to open configuration file")

	var c PackageJSON
	json.Unmarshal(raw, &c)

	if c.Disseminate == (PackageDisseminate{}) {
		warn("No disseminate tag in configuration file. { disseminate: {} } is required.")
	}

	return c
}

// Log messages from git
func getGitlogMessage(n string) string {

	// exec.Command requires a single string for third arg.  Combine strings
	s := []string{ "git log -n", n }

	gitlogCmd := exec.Command("bash", "-c", strings.Join(s, "  "))
	gitlogOut, err := gitlogCmd.CombinedOutput()
	check(err, string(gitlogOut))

	m := string(gitlogOut)
	// First test is if the most recent update actually has a merge message
	if m == "" {
		warn("No message was obtained from git log")
	}

	return m
}

// Grab has from our commit message,
func getHashString(m string) string {

	var hsh []string

	hsh = commTag.FindStringSubmatch(m)
	// Store the message as a hsh document
	if hsh[1] == "" {
		warn("Was not able to get the hash value for the git commit message. ")
	}

	return hsh[1]
}

// Get time from commit message
func getCommitTime(m string) string {

	var cTime []string

	cTime = dateTag.FindStringSubmatch(m)
	// If we do not have one get very upset
	if cTime[1] == "" {
		warn("Was not able to get the date/time value for the git commit message. ")
	}

	return cTime[1]
}

// See if it is a merge, and it a merge is required.
func checkMergeRequirement(m string, is_nMrg bool) {

	var has_merge bool

	has_merge = mer1Tag.MatchString(m)

	if ! has_merge && ! is_nMrg {
		warn("Last commit was not a merge")
	}
}

// Convert our json into an actual string
func (p PackageDisseminate) toString() string {
	return toJson(p)
}

// Decode json document into native go structure
func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	check(err, "Failed to decide json document.")

	return string(bytes)
}

func main(){

	// Message )onditions
	var maxLenth int = 280
	var minLenth int = 30
	var is_noMerge bool = false
	var is_print bool = false
	var configFile string = "./package.json"

	// var number string
	var message string
	var hash string
	var commitTime string

	// Build our timestamp
	now := time.Now()
	timestamp := []string{
		strconv.Itoa(now.Year()),
		fmt.Sprintf("%02d", int(now.Month())),
		strconv.Itoa(now.Day()),
		strconv.Itoa(now.Hour()),
		strconv.Itoa(now.Minute()) }

	// Less simpler printing
	p := fmt.Println

	// wordPtr := flag.String("filter", "merge", "a string")
	numbPtr := flag.Int("count", 1, "Number of previous commit messages to return.")
	maxLenthPtr := flag.Int("max", maxLenth, "Maximum allowable length of the commit message.")
	minLenthPtr := flag.Int("min", minLenth, "Minimum allowable length of the commit message.")
	is_noMergePtr := flag.Bool("nomerge", is_noMerge, "Include non-merge commits in results.")
	is_printPtr := flag.Bool("print", is_print, "Print output to terminal instead of to a file.")
	configFilePtr := flag.String("config", configFile, "Disseminate configuration file in JSON format.")
	flag.Parse()

	maxLenth = *maxLenthPtr
	minLenth = *minLenthPtr
	// is_noMerge = *is_noMergePtr
	is_print = *is_printPtr

	// Get out package information
	pkgs := getPackage(*configFilePtr)
	// Get our package messages
	message = getGitlogMessage(strconv.Itoa(*numbPtr))
	// Hash values for uniq storage
	hash = getHashString(message)
	// Need different build and commit times
	commitTime = getCommitTime(message)
	// If we have  a merge requirement let us know
	checkMergeRequirement(message, *is_noMergePtr)

	// Now we remove the author, date, and commit from the message
	message = authTag.ReplaceAllString(message, "")
	message = dateTag.ReplaceAllString(message, "")
	message = commTag.ReplaceAllString(message, "")
	message = mer1Tag.ReplaceAllString(message, "")
	message = mer2Tag.ReplaceAllString(message, "")
	message = strings.TrimSpace(message)

	// Message tests, to make sure they have a certain "quality"
	if message == "" {
		warn("Commit message is empty.")
	} else if len(message) <= minLenth {
		warn("Commit message does not meet minimum allowable length.")
	} else if len(message) >= maxLenth {
		warn("Commit message exceeds the maximum allowable length.")
	} else  {
		res := &response{
			Commit: hash,
			Date: commitTime,
			Time: strings.Join(timestamp, ""),
			Message: message,
			Disseminate: pkgs.Disseminate}
		resJson, _ := json.Marshal(res)

		// Print or write.
		if is_print {
			p(string(resJson))
		} else {
			err := ioutil.WriteFile("./disseminate.json", resJson, 0644)
			check(err, "Unable to write to file disseminate.json")
			p("disseminate.json update")
		}
	}
}
