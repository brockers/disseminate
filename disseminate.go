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
	Disseminate PackageDisseminate `json:"ai_disseminate"`
}


func check(e error, s string) {
	if e != nil {
		fmt.Println(s)
		panic(e)
	}
}

func getPackage() PackageJSON {
	raw, err := ioutil.ReadFile("./package.json")
	check(err,"Unable to open package.json")

	var c PackageJSON
	json.Unmarshal(raw, &c)

	return c
}

func (p PackageDisseminate) toString() string {
	return toJson(p)
}

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

	// var number string
	var message string
	var hash []string
	var commitTime []string
	now := time.Now()

	// Regular Expression filters
	var authTag = regexp.MustCompile(`Author: .*\n`)
	var dateTag = regexp.MustCompile(`Date:\s+(.*)\n`)
	var commTag = regexp.MustCompile(`commit ([a-z0-9]*)\n`)
	var mer1Tag = regexp.MustCompile(`Merge pull request .*\n`)
	var mer2Tag = regexp.MustCompile(`Merge: .*\n`)

	// Build our timestamp
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
	flag.Parse()

	maxLenth = *maxLenthPtr
	minLenth = *minLenthPtr
	is_noMerge = *is_noMergePtr

	// Get out package information
	pkgs := getPackage()

	// exec.Command requires a single string for third arg.  Combine strings
	s := []string{ "git log -n", strconv.Itoa(*numbPtr) }

	gitlogCmd := exec.Command("bash", "-c", strings.Join(s, "  "))
	gitlogOut, err := gitlogCmd.CombinedOutput()

	if err != nil {
		// log := spew.Sdump(gitlogCmd)
		// p("ERROR:", log)
		p("ERROR", gitlogOut)
		panic(err)
	}

	message = string(gitlogOut)
	// First test is if the most recent update actually has a merge message
	if message == "" {
		p("ERROR: No message was obtained from git log")
		os.Exit(1)
	}

	hash = commTag.FindStringSubmatch(message)
	commitTime = dateTag.FindStringSubmatch(message)

	// Store the message as a hash document
	if hash[1] == "" {
		p(message)
		panic("Was not able to get the hash value for the git commit message. ")
	}

	// Store the datestamp
	if commitTime[1] == "" {
		p(message)
		panic("Was not able to get the date/time value for the git commit message. ")
	}

	has_merge := mer1Tag.MatchString(message)

	if ! has_merge && ! is_noMerge {
		p("ERROR: Last commit was not a merge")
		os.Exit(1)
	} else {

		// Now we remove the author, date, and commit from the message
		message = authTag.ReplaceAllString(message, "")
		message = dateTag.ReplaceAllString(message, "")
		message = commTag.ReplaceAllString(message, "")
		message = mer1Tag.ReplaceAllString(message, "")
		message = mer2Tag.ReplaceAllString(message, "")
		message = strings.TrimSpace(message)

		// Message tests, to make sure they have a certain "quality"
		if message == "" {
			p("ERROR: Commit message is empty.")
			os.Exit(1)
		} else if len(message) <= minLenth {
			p("ERROR: Commit message does not meet minimum allowable length.")
			os.Exit(1)
		} else if len(message) >= maxLenth {
			p("ERROR: Commit message exceeds the maximum allowable length.")
			os.Exit(1)
		} else  {
			// p("SUCCESS:", hash[1])
			// p(message)
			// timestamp, _ := now.MarshalJSON()
			res := &response{
				Commit: hash[1],
				Date: commitTime[1],
				Time: strings.Join(timestamp, ""),
				Message: message,
			 	Disseminate: pkgs.Disseminate}
			resJson, _ := json.Marshal(res)
			p(string(resJson))
		}
	}

	// Lets try reading from a file
	// dat, err := ioutil.ReadFile("./package.json")
	// check(err)
	//
	// x := new(PackageJSON)
	// dec := json.NewDecoder(dat).Decode(x)
	// p("============ PRINT FROM JSON DOCK ======================")
	// for _, p := range pkgs {
		// fmt.Println(pkgs.Disseminate.toString())
	// }

	// fmt.Print(string(dat))
	// fmt.Print(len(dat))
}
